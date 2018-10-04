package core

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic("failed Int to hex:", err)
	}
	return buff.Bytes()
}
