package ledger

import (
	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/go-libs/time"
)

const (
	FeaturePostCommitVolumes            = "POST_COMMIT_VOLUMES"
	FeaturePostCommitEffectiveVolumes   = "POST_COMMIT_EFFECTIVE_VOLUMES"
	FeatureHashLogs                     = "HASH_LOGS"
	FeatureAccountMetadataHistories     = "ACCOUNT_METADATA_HISTORIES"
	FeatureTransactionMetadataHistories = "TRANSACTION_METADATA_HISTORIES"
	FeatureIndexAddressSegments         = "INDEX_ADDRESS_SEGMENTS"
	FeatureIndexTransactionAccounts     = "INDEX_TRANSACTION_ACCOUNTS"

	StateInitializing = "initializing"
	StateInUse        = "in-use"

	DefaultBucket = "_default"
)

type FeatureSet map[string]string

func (f FeatureSet) With(feature, value string) FeatureSet {
	ret := FeatureSet{}
	for k, v := range f {
		ret[k] = v
	}
	ret[feature] = value

	return ret
}

var DefaultFeatures = FeatureSet{
	FeaturePostCommitVolumes:            "SYNC",
	FeaturePostCommitEffectiveVolumes:   "SYNC",
	FeatureHashLogs:                     "SYNC",
	FeatureAccountMetadataHistories:     "SYNC",
	FeatureTransactionMetadataHistories: "SYNC",
	FeatureIndexAddressSegments:         "ON",
	FeatureIndexTransactionAccounts:     "ON",
}

var MinimalFeatureSet = FeatureSet{
	FeaturePostCommitVolumes:            "DISABLED",
	FeaturePostCommitEffectiveVolumes:   "DISABLED",
	FeatureHashLogs:                     "DISABLED",
	FeatureAccountMetadataHistories:     "DISABLED",
	FeatureTransactionMetadataHistories: "DISABLED",
	FeatureIndexAddressSegments:         "OFF",
	FeatureIndexTransactionAccounts:     "OFF",
}

var FeatureConfigurations = map[string][]string{
	FeaturePostCommitVolumes:            {"SYNC", "DISABLED"},
	FeaturePostCommitEffectiveVolumes:   {"SYNC", "DISABLED"},
	FeatureHashLogs:                     {"SYNC", "DISABLED"},
	FeatureAccountMetadataHistories:     {"SYNC", "DISABLED"},
	FeatureTransactionMetadataHistories: {"SYNC", "DISABLED"},
	FeatureIndexAddressSegments:         {"ON", "OFF"},
	FeatureIndexTransactionAccounts:     {"ON", "OFF"},
}

type Configuration struct {
	Bucket   string            `json:"bucket"`
	Metadata metadata.Metadata `json:"metadata"`
	Features map[string]string `bun:"features,type:jsonb" json:"features"`
}

func (c *Configuration) SetDefaults() {
	if c.Bucket == "" {
		c.Bucket = DefaultBucket
	}
	if c.Features == nil {
		c.Features = map[string]string{}
	}

	for key, value := range DefaultFeatures {
		if _, ok := c.Features[key]; !ok {
			c.Features[key] = value
		}
	}
}

func NewDefaultConfiguration() Configuration {
	return Configuration{
		Bucket:   DefaultBucket,
		Metadata: metadata.Metadata{},
		Features: DefaultFeatures,
	}
}

type Ledger struct {
	Configuration
	Name    string    `json:"name"`
	AddedAt time.Time `json:"addedAt"`
	State   string    `json:"-"`
}

func (l Ledger) HasFeature(feature, value string) bool {
	return l.Features[feature] == value
}

func New(name string, configuration Configuration) Ledger {
	return Ledger{
		Configuration: configuration,
		Name:          name,
		AddedAt:       time.Now(),
		State:         StateInitializing,
	}
}

func NewWithDefaults(name string) Ledger {
	return New(name, NewDefaultConfiguration())
}
