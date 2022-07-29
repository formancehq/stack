package auth

// TODO: Make a first class entity for scopes
var Scopes = Array[string]{
	"transactions:read",
	"transactions:write",
	"accounts:read",
	"accounts:write",
	"stats",
	"search",
	"payments:write",
	"payments:read",
	"connectors:read",
	"connectors:write",
}
