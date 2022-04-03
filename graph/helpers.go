package graph

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Fetch[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-TBA-Auth-Key", "Q7k2X4Erid8KqVXRt82GICzJAC9ZkIpBogQHJ4sIwJACU9u29bHNE6AElJWRKMk4")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	jsonResponse, _ := ioutil.ReadAll(res.Body)

	a := new(T)
	json.Unmarshal(jsonResponse, a)
	return a, nil
}
