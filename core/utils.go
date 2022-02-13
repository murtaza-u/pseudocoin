package core

import (
	"bytes"
	"encoding/binary"
	"os"
)

func IntToBytes(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(), err
}

func DBExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
