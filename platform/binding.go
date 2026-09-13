package platform

import (
	"encoding/json"
	"fmt"
)

func (f *Frame) IpcBinding(id uint64, method string, params json.RawMessage) {

	fmt.Println(method)

}
