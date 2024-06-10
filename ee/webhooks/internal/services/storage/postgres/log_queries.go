package storage 

const (

	insertLogQuery string = "INSERT INTO logs (id, channel, payload, created_at) VALUES (?, ?, ?, ?);"
	selectFreshLogs string = "SELECT * FROM logs WHERE channel IN (?) AND created_at > ? ORDER BY created_at ASC;"
	

)