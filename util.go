package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

func fetchJSONFromURL(url string, headers map[string][]string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

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

func getEligibleHeaders(allHeaders map[string][]string, passHeaders []string) map[string][]string {
	out := make(map[string][]string)
	for key, values := range allHeaders {
		if arrayHasLowerCase(passHeaders, key) {
			out[key] = values
		}
	}
	return out
}

func arrayHasLowerCase(haystack []string, needle string) bool {
	for _, v := range haystack {
		if strings.TrimSpace(strings.ToLower(v)) == strings.TrimSpace(strings.ToLower(needle)) {
			return true
		}
	}
	return false
}
