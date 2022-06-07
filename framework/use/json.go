package use

import "encoding/json"

func JsonEncode(element interface{}) string {
	data, err := json.Marshal(element)

	AbortUnless(err)

	return string(data)
}

func JsonDecode(element interface{}, data string) {
	err := json.Unmarshal([]byte(data), element)

	AbortUnless(err)
}
