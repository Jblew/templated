package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
	for k, v := range headers {
		joinedV := strings.Join(v, ",")
		if joinedV != "" && req.Header.Get(k) != "" {
			req.Header.Add(k, joinedV)
		}
	}

	if isVerbose {
		fmt.Printf("fetchJSON %s, headers: %+v", url, req.Header)
	}

	client := &http.Client{
		Timeout: 300 * time.Millisecond,
	}
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

func addHeadersToResponse(w http.ResponseWriter, headers map[string]string) {
	for k, v := range headers {
		if v != "" {
			w.Header().Add(k, v)
		}
	}
}

func addHeadersToRequestIfNotExist(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		if v != "" && req.Header.Get(k) != "" {
			req.Header.Add(k, v)
		}
	}
}
