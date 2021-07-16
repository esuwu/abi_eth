package fourbyte

import (
	"github.com/pkg/errors"
)

type ABI struct {
	Methods map[Selector]Method
}

// MethodById looks up a method by the 4-byte id,
// returns nil if none found.
func (abi *ABI) MethodById(selector Selector) (Method, error) {
	if method, ok := abi.Methods[selector]; ok {
		return method, nil
	}
	return Method{}, errors.Errorf("no method with id: %#x", selector[:])
}
