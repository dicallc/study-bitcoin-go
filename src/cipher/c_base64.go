package main

import (
	"encoding/base64"
	"fmt"
	"utils"
)

func main() {
	msg := "hello ,世界"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	utils.CheckErr("", err)
	fmt.Println(string(decoded))

}
