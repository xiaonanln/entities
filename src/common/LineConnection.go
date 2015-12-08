package common

import "net"

const (
	MAXLINE = 4096
)

type LineConnection struct {
	Connection
	recvBuf []byte
	lineBuf []byte
}

func NewLineConnection(conn net.Conn) *LineConnection {
	return &LineConnection{
		Connection: Connection{conn},
		recvBuf:    make([]byte, 4096),
	}
}

func (self *LineConnection) RecvLine(pline *[]byte) error {
	line := make([]byte, MAXLINE)
	pos := 0
	eol := false

	for _, c := range self.lineBuf {
		if c == '\n' {
			eol = true
			break
		} else {
			line[pos] = c
			pos += 1
			if pos == MAXLINE {
				// line full
				eol = true
				break
			}
		}
	}

	if eol {
		*pline = line[:pos]
		return nil
	}

	_, err := self.Read(self.recvBuf)
	if err != nil {
		return err
	}

	return nil
}

func (self *LineConnection) SendLine(line []byte) error {
	err := self.SendAll(line)
	if err != nil {
		return err
	}
	return self.SendAll([]byte{'\n'})
}
