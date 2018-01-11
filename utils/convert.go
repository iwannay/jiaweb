package utils

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func String2Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Interface2Struct(in interface{}, out interface{}) error {
	byteData, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteData, out)
}
