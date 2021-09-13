package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
	fmt.Printf("Got headers of request: %+v\n", headers)
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
		if isVerbose {
			fmt.Printf("fetchJSON %s FAILED [%d], err: %+v\n", url, resp.StatusCode, err)
		}
		return out, err
	}

	err = json.Unmarshal(responseBody, &out)
	if err != nil {
		if isVerbose {
			fmt.Printf("fetchJSON %s FAILED [%d], err: %+v\n", url, resp.StatusCode, err)
		}
		return out, err
	}

	if resp.StatusCode != 200 {
		outError := ""
		_, hasError := out["error"]
		if hasError {
			outError = fmt.Sprintf("%+v", out["error"])
		}
		if isVerbose {
			fmt.Printf("fetchJSON %s FAILED [%d], headers: %+v\n", url, resp.StatusCode, req.Header)
		}
		return out, fmt.Errorf("fetchJSON %s FAILED [%d]: %s", url, resp.StatusCode, outError)
	} else {
		if isVerbose {
			fmt.Printf("fetchJSON %s [%d], headers: %+v\n", url, resp.StatusCode, req.Header)
		}
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
