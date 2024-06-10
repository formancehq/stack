package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/formancehq/webhooks/internal/commons"
)	


const (
	selectOneHookQuery string = "SELECT * FROM configs WHERE id = ?"
	selectHooksQuery  = "SELECT * FROM configs WHERE status != ?" 
	selectHooksWithPaginationQuery  = "SELECT * FROM configs WHERE  status !=  ? ORDER By name LIMIT ? OFFSET ?"
	insertHookQuery  = "INSERT INTO configs (id, name, status, event_types, endpoint, secret, created_at, date_status, retry) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	updateHookStatusQuery = "UPDATE configs SET status = ?, date_status = ? WHERE id = ?"
	updateHookSecretQuery = "UPDATE configs SET secret = ? WHERE id = ?"
	updateHookEndpointQuery = "UPDATE configs SET endpoint = ? WHERE id = ?"
	updateHookRetryQuery = "UPDATE configs SET retry = ? WHERE id = ?"
)


func (store PostgresStore) GetHook(index string) (commons.Hook, error){
	var hook commons.Hook

	err := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableHooks.Name). 
	Where("id = ?", index). 
	Scan(context.Background(), &hook)

	if err == sql.ErrNoRows {
		return hook, nil 	
	}

	return hook, err
}

func (store PostgresStore) SaveHook(hook commons.Hook) error {

	
	_,err := store.db.NewRaw(insertHookQuery, 
							hook.ID, 
							hook.Name,
							hook.Status,
							StrArray(hook.Events),
							hook.Endpoint,
							hook.Secret,
							hook.DateCreation,
							hook.DateStatus, 
							hook.Retry).
			Exec(context.Background())

	return err
}



func (store PostgresStore) ActivateHook(index string) (commons.Hook, error){
	return store.changeHookStatus(index, commons.EnableStatus)
}

func (store PostgresStore) DeactivateHook(index string) (commons.Hook, error){
	return store.changeHookStatus(index, commons.DisableStatus)
}

func (store PostgresStore) DeleteHook(index string) (commons.Hook, error) {
	return store.changeHookStatus(index, commons.DeleteStatus)
}

func (store PostgresStore) changeHookStatus(index string, status commons.HookStatus) (commons.Hook, error) {
	var hook commons.Hook
	
	_, err := store.db.NewRaw(updateHookStatusQuery, string(status), "NOW()", index).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		return hook, nil
	}

	return hook, err 
}



func (store PostgresStore) UpdateHookEndpoint(index string, endpoint string) (commons.Hook, error) {
	return store.changeHookColumns(index, "ENDPOINT", endpoint)
}

func (store PostgresStore) UpdateHookSecret(index string, secret string) (commons.Hook, error) {
	return store.changeHookColumns(index, "SECRET", secret)
}

func (store PostgresStore) UpdateHookRetry(index string, retry bool) (commons.Hook, error){
	return store.changeHookColumns(index, "RETRY", retry)
}

func (store PostgresStore) changeHookColumns(index string, columnName string, value any) (commons.Hook, error){
	var hook commons.Hook

	updateRaw := fmt.Sprintf("%s = ?",TableHooks.Columns[columnName])
	conditionRaw := fmt.Sprintf("id = ?")

	query := updateQuery.Fill(TableHooks.Name, updateRaw, conditionRaw)
	
	_, err := store.db.NewRaw(string(query), value, index).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		return hook, nil
	}

	return hook, err 

}



func (store PostgresStore) GetHooks(page, size int, filterEndpoint string) (*[]*commons.Hook, bool, error){
	res := make([]*commons.Hook, 0)
	hasMore := false 

	q := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableHooks.Name). 
	Where("status != ?", commons.DeleteStatus).
	Limit(size+1).
	Offset(size*page)

	if(filterEndpoint != ""){
		q = q.Where("endpoint = ?", filterEndpoint)
	}


	err := q.Scan(context.Background(), &res)

	if(err != nil){
		return &res, hasMore, err
	}

	hasMore = len(res) == (size+1)

	if(hasMore){
		res = res[0:size]
	}

	return &res, hasMore, err
	
}

func (store PostgresStore) LoadHooks() (*[]*commons.Hook, error){
	res := make([]*commons.Hook, 0)

	err := store.db.NewSelect().
	ColumnExpr("*").
	Table(TableHooks.Name). 
	Where("status != ?", commons.DeleteStatus).
	Scan(context.Background(), &res)

	return &res, err
}



