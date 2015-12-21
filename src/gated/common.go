package gated

const (
	CMD_INVALID        = iota // invalid cmd: 0
	CMD_RPC            = iota // RPC call from client -> server, or server -> client
	CMD_CREATE_ENTITY  = iota // Create entity on client
	CMD_DESTROY_ENTITY = iota // Destroy entity
)
