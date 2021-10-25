package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

var host string

func Test_main(t *testing.T) {
	host = "8081"
	Test_CREATE_GROUP(t)
	host = "8082"
	Test_CREATE_GROUP(t)
	host = "8083"
	Test_CREATE_GROUP(t)
	host = "8080"
	Test_CREATE_GROUP(t)
	time.Sleep(1 * time.Second)
	Test_STARTGROUP(t)
	time.Sleep(3 * time.Second)
	Test_SENDMESSAGE(t)
}

func Test_CREATE_GROUP(t *testing.T) {
	url := "http://localhost:" + host + "/multicast/v1/groups"
	method := "POST"

	payload := strings.NewReader(`{"multicast_type":"BMULTICAST","multicast_id":"PROVA"}`)

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

func Test_STARTGROUP(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/:MId"
	method := "PUT"

	payload := strings.NewReader(`{"multicast_id":"PROVA"}`)

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

func Test_SENDMESSAGE(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/messaging/:Mid"
	method := "POST"

	payload := strings.NewReader(`{"multicast_id":"PROVA"}`)

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
