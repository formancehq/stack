package core

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type LogType int16

const (
	SetMetadataLogType         LogType = iota // "SET_METADATA"
	NewTransactionLogType                     // "NEW_TRANSACTION"
	RevertedTransactionLogType                // "REVERTED_TRANSACTION"
)

func (l LogType) String() string {
	switch l {
	case SetMetadataLogType:
		return "SET_METADATA"
	case NewTransactionLogType:
		return "NEW_TRANSACTION"
	case RevertedTransactionLogType:
		return "REVERTED_TRANSACTION"
	}

	return ""
}

// TODO(polo): create Log struct and extended Log struct
type Log struct {
	ID        uint64     `json:"id"`
	Type      LogType    `json:"type"`
	Data      marshaller `json:"data"`
	Hash      []byte     `json:"hash"`
	Date      Time       `json:"date"`
	Reference string     `json:"reference"`
}

func (l *Log) UnmarshalJSON(data []byte) error {
	type auxLog Log
	type log struct {
		auxLog
		Data json.RawMessage `json:"data"`
	}
	rawLog := log{}
	if err := json.Unmarshal(data, &rawLog); err != nil {
		return err
	}

	var err error
	rawLog.auxLog.Data, err = HydrateLogFromJSON(rawLog.Type, rawLog.Data)
	if err != nil {
		return err
	}
	*l = Log(rawLog.auxLog)
	return err
}

func (l Log) WithDate(date Time) Log {
	l.Date = date
	return l
}

func (l Log) WithReference(reference string) Log {
	l.Reference = reference
	return l
}

type AccountMetadata map[string]metadata.Metadata

func (m AccountMetadata) Marshal(buf *Buffer) {
	buf.writeUInt64(uint64(len(m)))
	if len(m) == 0 {
		return
	}
	accounts := collectionutils.Keys(m)
	if len(accounts) > 1 {
		sort.Strings(accounts)
	}

	for _, account := range accounts {
		buf.writeString(account)
		marshalMetadata(buf, m[account])
	}
}

func (m AccountMetadata) Unmarshal(buf *Buffer) {
	numberOfEntries := buf.readUInt64()
	for i := uint64(0); i < numberOfEntries; i++ {
		account := buf.readString()
		metadata := metadata.Metadata{}
		unmarshalMetadata(buf, metadata)
		m[account] = metadata
	}
}

type NewTransactionLogPayload struct {
	Transaction     Transaction     `json:"transaction"`
	AccountMetadata AccountMetadata `json:"accountMetadata"`
}

func (n *NewTransactionLogPayload) Unmarshal(buf *Buffer) {
	n.AccountMetadata.Unmarshal(buf)
	n.Transaction.Unmarshal(buf)
}

func (n NewTransactionLogPayload) Marshal(buf *Buffer) {
	n.AccountMetadata.Marshal(buf)
	n.Transaction.Marshal(buf)
}

func NewTransactionLogWithDate(tx Transaction, accountMetadata map[string]metadata.Metadata, time Time) Log {
	// Since the id is unique and the hash is a hash of the previous log, they
	// will be filled at insertion time during the batch process.
	return Log{
		Type: NewTransactionLogType,
		Date: time,
		Data: &NewTransactionLogPayload{
			Transaction:     tx,
			AccountMetadata: accountMetadata,
		},
	}
}

func NewTransactionLog(tx Transaction, accountMetadata map[string]metadata.Metadata) Log {
	return NewTransactionLogWithDate(tx, accountMetadata, tx.Timestamp).WithReference(tx.Reference)
}

type SetMetadataLogPayload struct {
	TargetType string            `json:"targetType"`
	TargetID   any               `json:"targetId"`
	Metadata   metadata.Metadata `json:"metadata"`
}

func (s *SetMetadataLogPayload) Unmarshal(buf *Buffer) {
	s.TargetType = buf.readString()
	switch s.TargetType {
	case MetaTargetTypeAccount:
		s.TargetID = buf.readString()
	case MetaTargetTypeTransaction:
		s.TargetID = buf.readUInt64()
	}
	s.Metadata = metadata.Metadata{}
	unmarshalMetadata(buf, s.Metadata)
}

func (s SetMetadataLogPayload) Marshal(buf *Buffer) {
	buf.writeString(s.TargetType)
	switch targetID := s.TargetID.(type) {
	case string:
		buf.writeString(targetID)
	case uint64:
		buf.writeUInt64(targetID)
	}
	marshalMetadata(buf, s.Metadata)
}

func (s *SetMetadataLogPayload) UnmarshalJSON(data []byte) error {
	type X struct {
		TargetType string            `json:"targetType"`
		TargetID   json.RawMessage   `json:"targetId"`
		Metadata   metadata.Metadata `json:"metadata"`
	}
	x := X{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	var id interface{}
	switch strings.ToUpper(x.TargetType) {
	case strings.ToUpper(MetaTargetTypeAccount):
		id = ""
		err = json.Unmarshal(x.TargetID, &id)
	case strings.ToUpper(MetaTargetTypeTransaction):
		id, err = strconv.ParseUint(string(x.TargetID), 10, 64)
	default:
		panic("unknown type")
	}
	if err != nil {
		return err
	}

	*s = SetMetadataLogPayload{
		TargetType: x.TargetType,
		TargetID:   id,
		Metadata:   x.Metadata,
	}
	return nil
}

func NewSetMetadataLog(at Time, metadata SetMetadataLogPayload) Log {
	// Since the id is unique and the hash is a hash of the previous log, they
	// will be filled at insertion time during the batch process.
	return Log{
		Type: SetMetadataLogType,
		Date: at,
		Data: &metadata,
	}
}

type RevertedTransactionLogPayload struct {
	RevertedTransactionID uint64
	RevertTransaction     Transaction
}

func (r *RevertedTransactionLogPayload) Unmarshal(buf *Buffer) {
	r.RevertedTransactionID = buf.readUInt64()
	r.RevertTransaction.Unmarshal(buf)
}

func (r RevertedTransactionLogPayload) Marshal(buf *Buffer) {
	buf.writeUInt64(r.RevertedTransactionID)
	r.RevertTransaction.Marshal(buf)
}

func NewRevertedTransactionLog(at Time, revertedTxID uint64, tx Transaction) Log {
	return Log{
		Type: RevertedTransactionLogType,
		Date: at,
		Data: &RevertedTransactionLogPayload{
			RevertedTransactionID: revertedTxID,
			RevertTransaction:     tx,
		},
	}
}

func HydrateLogFromRaw(_type LogType, data []byte) (marshaller, error) {
	var payload marshaller
	switch _type {
	case NewTransactionLogType:
		payload = &NewTransactionLogPayload{}
	case SetMetadataLogType:
		payload = &SetMetadataLogPayload{}
	case RevertedTransactionLogType:
		payload = &RevertedTransactionLogPayload{}
	default:
		panic("unknown type " + _type.String())
	}

	buffer := NewBuffer(data)
	payload.Unmarshal(buffer)

	return payload, nil
}

func HydrateLogFromJSON(_type LogType, data []byte) (marshaller, error) {
	var payload any
	switch _type {
	case NewTransactionLogType:
		payload = &NewTransactionLogPayload{}
	case SetMetadataLogType:
		payload = &SetMetadataLogPayload{}
	case RevertedTransactionLogType:
		payload = &RevertedTransactionLogPayload{}
	default:
		panic("unknown type " + _type.String())
	}
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(payload).Interface().(marshaller), nil
}

type Accounts map[string]Account

type LogHolder struct {
	Log      *Log
	Ingested chan struct{}
}

func NewLogHolder(log *Log) *LogHolder {
	return &LogHolder{
		Log:      log,
		Ingested: make(chan struct{}),
	}
}
