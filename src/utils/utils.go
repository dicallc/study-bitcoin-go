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

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
