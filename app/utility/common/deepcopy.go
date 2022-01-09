package common

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func DeepCopyByGob(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}

func DeepCopyByJson(src interface{}) (interface{}, error) {
	b, err := json.Marshal(src)
	if nil != err {
		return nil, err
	}
	var dst interface{}
	err = json.Unmarshal(b, dst)
	return dst, err
}
