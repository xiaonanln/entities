package common

import "net"

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) Connection {
	return Connection{conn}
}

func (self Connection) String() string {
	return self.conn.RemoteAddr().String()
}

func (self Connection) RecvByte() (byte, error) {
	buf := []byte{0}
	for {
		n, err := self.conn.Read(buf)
		if n >= 1 {
			return buf[0], nil
		} else if err != nil {
			return 0, err
		}
	}
}

func (self Connection) SendByte(b byte) error {
	buf := []byte{b}
	for {
		n, err := self.conn.Write(buf)
		if n >= 1 {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func (self Connection) RecvAll(buf []byte) error {
	for len(buf) > 0 {
		n, err := self.conn.Read(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func (self Connection) SendAll(data []byte) error {
	for len(data) > 0 {
		n, err := self.conn.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

func (self Connection) Read(data []byte) (int, error) {
	return self.conn.Read(data)
}

func (self Connection) Write(data []byte) (int, error) {
	return self.conn.Write(data)
}

func (self Connection) Close() error {
	return self.conn.Close()
}
