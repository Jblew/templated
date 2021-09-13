package main

import (
	"encoding/json"
	"io/ioutil"
)

func marshallToString(obj interface{}) string {
	marshalledBytes, _ := json.Marshal(obj)
	return string(marshalledBytes)
}

func readJSONFile(path string) (map[string]interface{}, error) {
	bytes, _ := ioutil.ReadFile(path)
	out := make(map[string]interface{})
	err := json.Unmarshal(bytes, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}
