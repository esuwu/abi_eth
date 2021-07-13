package fourbyte

import (
	"crypto/sha256"
	"os"
	"time"
)

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}


// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"4byte.json": _4byteJson,
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

func _4byteJson() (*asset, error) {
	bytes, err := _4byteJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4byte.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x97, 0xc1, 0x67, 0x6, 0x1e, 0x89, 0x76, 0xf7, 0x19, 0xd6, 0x8b, 0x43, 0xb4, 0x1c, 0xf6, 0xab, 0x7f, 0xc7, 0xc4, 0xca, 0x25, 0x21, 0x2, 0x13, 0x6d, 0x5b, 0xe2, 0x72, 0xb1, 0x7, 0xbc, 0x77}}
	return a, nil
}

func _4byteJsonBytes() ([]byte, error) {
	return __4byteJson, nil
}
var __4byteJson = []byte(`{
"a9059cbb": "transfer(address,uint256)"
}`)
