package main

import (
	"encoding/hex"
	"fmt"
	"github.com/abi_eth/fourbyte"
	"github.com/abi_eth/types"
	"strings"
)



func parse(data []byte) {
	db, err := fourbyte.New()
	if err != nil {
		fmt.Println(err)
	}
	messages := types.ValidationMessages{}
	db.ValidateCallData(nil, data, &messages)
	for _, m := range messages.Messages {
		fmt.Printf("%v: %v\n", m.Typ, m.Message)
	}
}

// Example
// ./abidump a9059cbb000000000000000000000000ea0e2dc7d65a50e77fc7e84bff3fd2a9e781ff5c0000000000000000000000000000000000000000000000015af1d78b58c40000
func main() {

	hexdata := "a9059cbb000000000000000000000000ea0e2dc7d65a50e77fc7e84bff3fd2a9e781ff5c0000000000000000000000000000000000000000000000015af1d78b58c40000"
	data, err := hex.DecodeString(strings.TrimPrefix(hexdata, "0x"))
	if err != nil {
		fmt.Println(err)
	}
	parse(data)

}
