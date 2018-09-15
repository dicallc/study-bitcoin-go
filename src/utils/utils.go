package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
func CheckErr(pos string, err error) {
	if err != nil {
		fmt.Println(" err occur:", err, "pos:", pos)
		os.Exit(1)
	}
}
