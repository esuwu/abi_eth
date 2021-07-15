package main

import (
	"encoding/hex"
	"fmt"
	"github.com/abi_eth/fourbyte"
	"strings"
)



func parse(data []byte) (*fourbyte.DecodedCallData, error){
	db, err := fourbyte.New()
	if err != nil {
		fmt.Println(err)
	}
	decodedData, err := db.ParseCallData(data)
	return decodedData, err
}

// Example
// ./abidump a9059cbb000000000000000000000000ea0e2dc7d65a50e77fc7e84bff3fd2a9e781ff5c0000000000000000000000000000000000000000000000015af1d78b58c40000
func main() {

	hexdata := "a9059cbb000000000000000000000000ea0e2dc7d65a50e77fc7e84bff3fd2a9e781ff5c0000000000000000000000000000000000000000000000015af1d78b58c40000"
	data, err := hex.DecodeString(strings.TrimPrefix(hexdata, "0x"))
	if err != nil {
		fmt.Println(err)
	}
	decodedData, err := parse(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodedData)

}
