package fourbyte

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// selectorRegexp is used to validate that a 4byte database selector corresponds
// to a valid ABI function declaration.
//
// Note, although uppercase letters are not part of the ABI spec, this regexp
// still accepts it as the general format is valid. It will be rejected later
// by the type checker.
var selectorRegexp = regexp.MustCompile(`^([^\)]+)\(([A-Za-z0-9,\[\]]*)\)`)

// parseSelector converts a method selector into an ABI JSON spec. The returned
// data is a valid JSON string which can be consumed by the standard abi package.
func parseSelector(unescapedSelector string) ([]byte, error) {
	// Define a tiny fake ABI struct for JSON marshalling
	type fakeArg struct {
		Type string `json:"type"`
	}
	type fakeABI struct {
		Name   string    `json:"name"`
		Type   string    `json:"type"`
		Inputs []fakeArg `json:"inputs"`
	}
	// Validate the unescapedSelector and extract it's components
	groups := selectorRegexp.FindStringSubmatch(unescapedSelector)
	if len(groups) != 3 {
		return nil, fmt.Errorf("invalid selector %q (%v matches)", unescapedSelector, len(groups))
	}
	name := groups[1]
	args := groups[2]

	// Reassemble the fake ABI and constuct the JSON
	arguments := make([]fakeArg, 0)
	if len(args) > 0 {
		for _, arg := range strings.Split(args, ",") {
			arguments = append(arguments, fakeArg{arg})
		}
	}
	return json.Marshal([]fakeABI{{name, "function", arguments}})
}

// ethDecodedArgument is an internal type to represent an argument parsed according
// to an ABI method signature.
type ethDecodedArgument struct {
	Soltype abi.Argument
	Value   interface{}
}

func (arg *ethDecodedArgument) DecodedValue() interface{} {
	return arg.Value
}

func (arg *ethDecodedArgument) InternalType() byte {
	return arg.Soltype.Type.T
}

type ArgDecoded interface {
	fmt.Stringer
	DecodedValue() interface{}
	InternalType() byte
}

// String implements stringer interface, tries to use the underlying value-type
func (arg ethDecodedArgument) String() string {
	var value string
	switch val := arg.Value.(type) {
	case fmt.Stringer:
		value = val.String()
	default:
		value = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%v: %v", arg.Soltype.Type.String(), value)
}

// DecodedCallData is an internal type to represent a method call parsed according
// to an ABI method signature.
type DecodedCallData struct {
	Signature string
	Name      string
	Inputs    []ArgDecoded
}

// String implements stringer interface for decodedCallData
func (cd DecodedCallData) String() string {
	args := make([]string, len(cd.Inputs))
	for i, arg := range cd.Inputs {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", cd.Name, strings.Join(args, ","))
}

func parseCallData(calldata []byte, unescapedAbidata string) (*DecodedCallData, error) {
	// Validate the call data that it has the 4byte prefix and the rest divisible by 32 bytes
	if len(calldata) < 4 {
		return nil, fmt.Errorf("invalid call data, incomplete method signature (%d bytes < 4)", len(calldata))
	}
	sigdata := calldata[:4]

	argdata := calldata[4:]
	if len(argdata)%32 != 0 {
		return nil, fmt.Errorf("invalid call data; length should be a multiple of 32 bytes (was %d)", len(argdata))
	}
	// Validate the called method and upack the call data accordingly
	abispec, err := abi.JSON(strings.NewReader(unescapedAbidata))
	if err != nil {
		return nil, fmt.Errorf("invalid method signature (%q): %v", unescapedAbidata, err)
	}
	method, err := abispec.MethodById(sigdata)
	if err != nil {
		return nil, err
	}
	values, err := method.Inputs.UnpackValues(argdata)
	if err != nil {
		return nil, fmt.Errorf("signature %q matches, but arguments mismatch: %v", method.String(), err)
	}
	// Everything valid, assemble the call infos for the signer
	decoded := DecodedCallData{Signature: method.Sig, Name: method.RawName}
	for i := 0; i < len(method.Inputs); i++ {
		decoded.Inputs = append(decoded.Inputs, &ethDecodedArgument{
			Soltype: method.Inputs[i],
			Value:   values[i],
		})
	}
	// We're finished decoding the data. At this point, we encode the decoded data
	// to see if it matches with the original data. If we didn't do that, it would
	// be possible to stuff extra data into the arguments, which is not detected
	// by merely decoding the data.
	encoded, err := method.Inputs.PackValues(values)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(encoded, argdata) {
		was := common.Bytes2Hex(encoded)
		exp := common.Bytes2Hex(argdata)
		return nil, fmt.Errorf("WARNING: Supplied data is stuffed with extra data. \nWant %s\nHave %s\nfor method %v", exp, was, method.Sig)
	}
	return &decoded, nil
}

// Database is a 4byte database with the possibility of maintaining an immutable
// set (embedded) into the process and a mutable set (loaded and written to file).
type Database struct {
	embedded map[string]string
	custom   map[string]string
}

// New loads the standard signature database embedded in the package.
func NewDatabase() (*Database, error) {
	db := &Database{make(map[string]string), make(map[string]string)}
	db.embedded = __4byteJson

	return db, nil
}

// This method does not validate the match, it's assumed the caller will do.
func (db *Database) Selector(id []byte) (string, error) {
	if len(id) < 4 {
		return "", fmt.Errorf("expected 4-byte id, got %d", len(id))
	}
	sig := hex.EncodeToString(id[:4])
	if selector, exists := db.embedded[sig]; exists {
		return selector, nil
	}
	if selector, exists := db.custom[sig]; exists {
		return selector, nil
	}
	return "", fmt.Errorf("signature %v not found", sig)
}

func (db *Database) MethodBySelector(id Selector) (Method, error) {
	if method, ok := erc20Methods[id]; ok {
		return method, nil
	}
	// TODO(nickeskov): support ride scripts metadata
	return Method{}, fmt.Errorf("signature %v not found", id.String())
}

// ValidateCallData checks if the ABI call-data + method selector (if given) can
// be parsed and seems to match.
func (db *Database) ParseCallData(data []byte) (*DecodedCallData, error) {

	// If the data is empty, we have a plain value transfer, nothing more to do
	if len(data) == 0 {
		return nil, errors.New("transaction doesn't contain data")
	}
	// Validate the call data that it has the 4byte prefix and the rest divisible by 32 bytes
	if len(data) < 4 {
		return nil, errors.New("transaction data is not valid ABI: missing the 4 byte call prefix")
	}
	if n := len(data) - 4; n%32 != 0 {
		return nil, errors.Errorf("transaction data is not valid ABI (length should be a multiple of 32 (was %d))", n)
	}
	embedded, err := db.Selector(data[:4])
	if err != nil {
		return nil, errors.Errorf("Transaction contains data, but the ABI signature could not be found: %v", err)
	}
	info, err := verifySelector(embedded, data)
	if err != nil {
		return nil, errors.Errorf("Transaction contains data, but provided ABI signature could not be verified: %v", err)
	}
	return info, nil

}

func (db *Database) ParseCallDataNew(data []byte) (*DecodedCallData, error) {
	// If the data is empty, we have a plain value transfer, nothing more to do
	if len(data) == 0 {
		return nil, errors.New("transaction doesn't contain data")
	}
	// Validate the call data that it has the 4byte prefix and the rest divisible by 32 bytes
	if len(data) < 4 {
		return nil, errors.New("transaction data is not valid ABI: missing the 4 byte call prefix")
	}
	if n := len(data) - 4; n%32 != 0 {
		return nil, errors.Errorf("transaction data is not valid ABI (length should be a multiple of 32 (was %d))", n)
	}
	var selector Selector
	copy(selector[:], data[:len(selector)])
	method, err := db.MethodBySelector(selector)
	if err != nil {
		return nil, errors.Errorf("Transaction contains data, but the ABI signature could not be found: %v", err)
	}

	info, err := parseArgData(&method, data[len(selector):])
	if err != nil {
		return nil, errors.Errorf("Transaction contains data, but provided ABI signature could not be verified: %v", err)
	}
	return info, nil
}

// verifySelector checks whether the ABI encoded data blob matches the requested
// function signature.
func verifySelector(functionSignature string, calldata []byte) (*DecodedCallData, error) {
	// Parse the functionSignature into an ABI JSON spec
	abidata, err := parseSelector(functionSignature)
	if err != nil {
		return nil, err
	}
	// Parse the call data according to the requested functionSignature
	return parseCallData(calldata, string(abidata))
}

type decodedArg struct {
	Soltype Argument
	Value   interface{}
}

func (da *decodedArg) String() string {
	var value string
	switch val := da.Value.(type) {
	case fmt.Stringer:
		value = val.String()
	default:
		value = fmt.Sprintf("%v", val)
	}
	return fmt.Sprintf("%v: %v", da.Soltype.Type.String(), value)
}

func (da *decodedArg) DecodedValue() interface{} {
	return da.Value
}

func (da *decodedArg) InternalType() byte {
	return byte(da.Soltype.Type.T)
}

func parseArgData(method *Method, argData []byte) (*DecodedCallData, error) {
	//method, err := abi.MethodById(selector)
	//if err != nil {
	//	return nil, errors.Wrapf(err, "failed to get method by id, id=%s", selector.String())
	//}
	values, err := method.Inputs.UnpackValues(argData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack Inputs arguments ABI data")
	}

	// TODO(nickeskov): use our types
	decoded := DecodedCallData{Signature: method.Sig.String(), Name: method.RawName}
	for i := 0; i < len(method.Inputs); i++ {
		decoded.Inputs = append(decoded.Inputs, &decodedArg{
			Soltype: method.Inputs[i],
			Value:   values[i],
		})
	}
	return &decoded, nil
}
