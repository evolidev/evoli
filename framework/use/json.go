package use

import (
	"encoding/json"
)

func JsonEncode(element interface{}) string {
	data, err := json.Marshal(element)

	AbortUnless(err)

	return string(data)
}

func JsonDecode(data string) interface{} {
	var mapData interface{}
	if err := json.Unmarshal([]byte(data), &mapData); err != nil {
		// TODO log error
		return nil
	}

	return mapData
}

func JsonDecodeObject(data string) map[string]interface{} {
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &mapData); err != nil {
		// TODO log error
		return nil
	}

	return mapData
}

func JsonDecodeStruct(data string, object any) {
	if err := json.Unmarshal([]byte(data), object); err != nil {
		// TODO log error
	}
}
