package storage

import (
	"fmt"
)

const (
	insertLogQuery  string = "INSERT INTO logs (id, channel, payload, created_at) VALUES (?, ?, ?, ?);"
	selectFreshLogs string = "SELECT * FROM logs WHERE channel IN (?) AND created_at > ? ORDER BY created_at ASC;"
)

func wrapWithLogQuery(query string) string {

	rawQuery := `
		BEGIN ;

		WITH wrap_query AS (
			%s
			RETURNING *
		)

		json_payload AS (
    		SELECT json_agg(wrap_query) AS payload
    		FROM wrap_query
		)

		INSERT INTO logs (id, channel, payload, created_at) VALUE (?, ?, ?, ?) ;
		SELECT * FROM wrap_query;
		
		COMMIT; 
	`

	return fmt.Sprintf(rawQuery, query)

}