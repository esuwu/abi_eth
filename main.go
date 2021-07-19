package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/abi_eth/fourbyte"
	"regexp"
	"strings"
)

func parse(data []byte) (*fourbyte.DecodedCallData, error) {
	db, err := fourbyte.NewDatabase()
	if err != nil {
		fmt.Println(err)
	}
	decodedData, err := db.ParseCallData(data)
	return decodedData, err
}

func parseNew(data []byte) (*fourbyte.DecodedCallData, error) {
	db, err := fourbyte.NewDatabase()
	if err != nil {
		fmt.Println(err)
	}
	decodedData, err := db.ParseCallDataNew(data)
	return decodedData, err
}

var selectorRegexp = regexp.MustCompile(`^([^\)]+)\(([A-Za-z0-9,\[\]]*)\)`)

func getJsonAbi(selector string) ([]byte, error) {
	// Define a tiny fake ABI struct for JSON marshalling
	type Arg struct {
		Type string `json:"type"`
	}
	type ABI struct {
		Name   string    `json:"name"`
		Type   string    `json:"type"`
		Inputs []Arg `json:"inputs"`
	}
	// Validate the unescapedSelector and extract it's components
	groups := selectorRegexp.FindStringSubmatch(selector)
	if len(groups) != 3 {
		return nil, fmt.Errorf("invalid selector %q (%v matches)", selector, len(groups))
	}
	name := groups[1]
	args := groups[2]

	// Reassemble the fake ABI and constuct the JSON
	arguments := make([]Arg, 0)
	if len(args) > 0 {
		for _, arg := range strings.Split(args, ",") {
			arguments = append(arguments, Arg{arg})
		}
	}
	return json.Marshal([]ABI{{name, "function", arguments}})
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
	fmt.Println(decodedData.Inputs[0].DecodedValue())
	fmt.Println(decodedData.Inputs[1].DecodedValue())
	fmt.Println(decodedData)

	fmt.Printf("\n\n ------------------------ \n\n")

	decodedDataNew, err := parseNew(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodedDataNew.Inputs[0].DecodedValue())
	fmt.Println(decodedDataNew.Inputs[1].DecodedValue())
	fmt.Println(decodedDataNew)


}
