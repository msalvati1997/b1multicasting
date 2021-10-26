package test

import (
	"encoding/json"
	"fmt"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_create_a_group(t *testing.T) {
	CREATEORJOIN_GROUP("8080", "COMULTICAST", "MYGROUP")
	CREATEORJOIN_GROUP("8081", "COMULTICAST", "MYGROUP")
	CREATEORJOIN_GROUP("8082", "COMULTICAST", "MYGROUP")
	CREATEORJOIN_GROUP("8083", "COMULTICAST", "MYGROUP")
}

func Test_start_group(t *testing.T) {
	StartGroupTest("8080", "MYGROUP")
	GetGroupTest("8080", "MYGROUP")
}

func Test_send_message(t *testing.T) {
	SENDMESSAGETEST("8081", "PROVA", "MYGROUP")
}

func Test_get_message(t *testing.T) {
	GETMESSAGESTEST("8080", "MYGROUP")
}

func Test_get_delivered_message(t *testing.T) {
	GETDELIVERMESSAGESTEST("8080", "MYGROUP")
}

func CREATEORJOIN_GROUP(host string, mtype string, mid string) {

	obj := multicastapp.MulticastReq{
		MulticastId:   mid,
		MulticastType: mtype,
	}
	jsn, _ := json.Marshal(obj)
	url := "http://localhost:" + host + "/multicast/v1/groups"
	method := "POST"
	payload := strings.NewReader(string(jsn))

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

func StartGroupTest(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/groups/" + group
	method := "PUT"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

func GetGroupTest(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/groups/" + group
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

type Message struct {
	Payload []byte
}

func SENDMESSAGETEST(host string, message string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/messaging/" + group
	method := "POST"
	m := []byte(message)
	obj := Message{m}
	json, _ := json.Marshal(obj)
	reader := strings.NewReader(string(json))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reader)

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

func GETMESSAGESTEST(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/messaging/" + group
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

func GETDELIVERMESSAGESTEST(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/deliver/" + group
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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
