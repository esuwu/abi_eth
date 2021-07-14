package fourbyte

import (
	"os"
	"time"
)

type asset struct {
	bytes  []byte
	info   os.FileInfo
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
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func _4byteJsonBytes() ([]byte, error) {
	return __4byteJson, nil
}
var __4byteJson = []byte(`{
"a9059cbb": "transfer(address,uint256)",
"23b872dd": "transferFrom(address,address,uint256)",
"ddf252ad": "Transfer(address,address,uint256)"
}`)

// "beabacc8": "transfer(address,address,uint256)",
