package common

type ClientId Eid

const (
	CID_LENGTH = EID_LENGTH
)

func NewClientId() ClientId {
	return ClientId(NewEid())
}
func (self ClientId) String() string {
	return Eid(self).String()
}
