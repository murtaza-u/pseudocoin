package core

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(), err
}
