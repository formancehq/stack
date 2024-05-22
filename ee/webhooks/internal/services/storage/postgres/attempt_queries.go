package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/webhooks/internal/commons"
)	


var TableAttempts = &Table{
	Name: "attempts",
	Columns: map[string]string{
		"ID":"id",
		"HOOKID":"webhook_id",
		"HOOKNAME":"hook_name", 
		"HOOKEP": "hook_endpoint", 
		"EVENT":"event", 
		"PAYLOAD": "payload", 
		"STATUSC": "status_code",
		"DOCCURED":"date_occured", 
		"STATUS":"status", 
		"DSTATUS": "date_status", 
		"COMMENT":"comment", 
		"NEXTTRY": "next_retry_after"},
}

func attemptColumnsAsListStr() string {
	var sb strings.Builder
	sb.WriteString(TableAttempts.Columns["ID"]+",")
	sb.WriteString(TableAttempts.Columns["HOOKID"]+",")
	sb.WriteString(TableAttempts.Columns["HOOKNAME"]+",")
	sb.WriteString(TableAttempts.Columns["HOOKEP"]+",")
	sb.WriteString(TableAttempts.Columns["EVENT"]+",")
	sb.WriteString(TableAttempts.Columns["PAYLOAD"]+",")
	sb.WriteString(TableAttempts.Columns["STATUSC"]+",")
	sb.WriteString(TableAttempts.Columns["DOCCURED"]+",")
	sb.WriteString(TableAttempts.Columns["STATUS"]+",")
	sb.WriteString(TableAttempts.Columns["DSTATUS"]+",")
	sb.WriteString(TableAttempts.Columns["COMMENT"]+",")
	sb.WriteString(TableAttempts.Columns["NEXTTRY"])
	
	return sb.String()

}

func (store PostgresStore) GetAttempt(index string)(commons.Attempt, error){
	var attempt commons.Attempt

	err := store.db.NewSelect().
	ColumnExpr("*"). 
	Table(TableAttempts.Name). 
	Where("id = ?", index). 
	Scan(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, err

}

func (store PostgresStore) SaveAttempt(attempt commons.Attempt) error {
	

		query := insertQuery.
		Fill(TableAttempts.Name, 
			attemptColumnsAsListStr(), 
			insertQuery.ValuesNb(len(TableAttempts.Columns)))
	
		_,err := store.db.NewRaw(string(query), 
							attempt.ID, 
							attempt.HookID,
							attempt.HookName,
							attempt.HookEndpoint, 
							attempt.Event,
							attempt.Payload,
							attempt.LastHttpStatusCode,
							attempt.DateOccured,
							string(attempt.Status),
							attempt.DateStatus,
							string(attempt.Comment),  
							attempt.NextTry).
				Exec(context.Background())


		return err
}


func (store PostgresStore) CompleteAttempt(index string) (commons.Attempt, error) {
	return store.ChangeAttemptStatus(index, commons.SuccessStatus, "")
}

func (store PostgresStore) AbortAttempt(index string, comment string) (commons.Attempt, error){
	return store.ChangeAttemptStatus(index, commons.AbortStatus, comment)
}

func (store PostgresStore) ChangeAttemptStatus(index string, status commons.AttemptStatus, comment string) (commons.Attempt, error){
	var attempt commons.Attempt

	updateRaw := fmt.Sprintf("%s = ?, %s = ?, %s = ?", TableAttempts.Columns["STATUS"], TableAttempts.Columns["DSTATUS"], TableAttempts.Columns["COMMENT"])
	conditionRaw := fmt.Sprintf("id = ?")

	query := updateQuery.Fill(TableAttempts.Name, updateRaw, conditionRaw)

	_, err := store.db.NewRaw(string(query), string(status), "NOW()", comment, index).Exec(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, err 
}

func (store PostgresStore) UpdateAttemptNextTry(index string, nextTry time.Time, statusCode int)(commons.Attempt, error){
	var attempt commons.Attempt

	updateRaw := fmt.Sprintf("%s = ?, %s = ?", TableAttempts.Columns["NEXTTRY"], TableAttempts.Columns["STATUSC"])
	conditionRaw := fmt.Sprintf("id = ?")

	query := updateQuery.Fill(TableAttempts.Name, updateRaw, conditionRaw)

	_, err := store.db.NewRaw(string(query), nextTry,statusCode, index).Exec(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, err 
}

func (store PostgresStore) GetWaitingAttempts(page, size int) (*[]*commons.Attempt, bool, error){
	res := make([]*commons.Attempt, 0)
	hasMore := false 

	err := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableAttempts.Name). 
	Where("(status = ?) OR (status = 'to_retry')", commons.WaitingStatus).  // 'to_retry' is for V1 compatibility...
	Limit(size+1).
	Offset(size*page).
	Scan(context.Background(), &res)

	hasMore = len(res) == (size+1)

	if(hasMore){
		res = res[0:size]
	}

	return &res, hasMore, err
}

func (store PostgresStore) GetAbortedAttempts(page, size int) (*[]*commons.Attempt, bool, error){
	res := make([]*commons.Attempt, 0)
	hasMore := false 

	err := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableAttempts.Name). 
	Where("(status = ?) OR (status = 'failed')", commons.AbortStatus).  // 'failed' is for V1 compatibility...
	Limit(size+1).
	Offset(size*page).
	Scan(context.Background(), &res)

	hasMore = len(res) == (size+1)

	if(hasMore){
		res = res[0:size]
	}

	return &res, hasMore, err
}


func (store PostgresStore) LoadWaitingAttempts() (*[]*commons.Attempt, error){
	res := make([]*commons.Attempt, 0)

	err := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableAttempts.Name). 
	Where("(status = ?) OR (status = 'to_retry')", commons.WaitingStatus).  // 'to_retry' is for V1 compatibility...
	Scan(context.Background(), &res)

	return &res, err
}







