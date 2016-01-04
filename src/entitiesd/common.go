package entitiesd

const (
	CMD_INVALID    = iota // invalid cmd: 0
	CMD_RPC        = iota // RPC call from client -> server, or server -> client
	CMD_NEW_ENTITY = iota // Create entity on client
	CMD_DEL_ENTITY = iota // Destroy entity
	CMD_NEW_CLIENT = iota // New client from gate
)
