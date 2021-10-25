package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_CREATE_GROUP(t *testing.T) {

	url := "http://localhost:8083/multicast/v1/groups/:MId"
	method := "GET"

	payload := strings.NewReader(`{"multicastId":"PROVA"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
