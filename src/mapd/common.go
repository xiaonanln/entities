package mapd

const (
	CMD_INVALID         = iota // invalid cmd: 0
	CMD_RPC             = iota // RPC call
	CMD_PID             = iota // set client Pid
	CMD_QUERY           = iota // query Eid : Pid
	CMD_SET             = iota // just set Eid : Pid
	CMD_LOCK_EID        = iota // Lock Eid for a period of time
	CMD_TRANSFER        = iota // transfer Eid -> Pid
	CMD_REGISTER_GLOBAL = iota // register global entity
)

const (
	REPLY_INVALID       = iota
	REPLY_RPC           = iota
	REPLY_OK            = iota // Reply ok: used when there is nothing to reply
	REPLY_REGISTER_FAIL = iota
)
