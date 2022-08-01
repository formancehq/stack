package constants

const (
	ServerHttpBindAddressFlag  = "server.http.bind_address"
	StorageMongoConnStringFlag = "storage.mongo.conn_string"

	DefaultMongoConnString = "mongodb://admin:admin@localhost:27017/"
	DefaultBindAddress     = ":8080"
)
