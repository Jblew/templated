package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func fetchJSONFromFile(path string) (map[string]interface{}, error) {
	bytes, _ := ioutil.ReadFile(path)
	out := make(map[string]interface{})
	err := json.Unmarshal(bytes, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func fetchJSONFromURL(url string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(responseBody, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}
