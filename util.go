package main

import "encoding/json"

func marshallToString(obj interface{}) string {
	marshalledBytes, _ := json.Marshal(obj)
	return string(marshalledBytes)
}
