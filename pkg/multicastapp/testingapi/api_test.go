package test

import (
	"encoding/json"
	"fmt"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting/test"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_create_a_group(t *testing.T) {
	CREATEORJOIN_GROUP("8080", "COMULTICAST", "COMULTICASTGROUP")
	CREATEORJOIN_GROUP("8081", "COMULTICAST", "COMULTICASTGROUP")
	CREATEORJOIN_GROUP("8082", "COMULTICAST", "COMULTICASTGROUP")
	CREATEORJOIN_GROUP("8083", "COMULTICAST", "COMULTICASTGROUP")
	////////////////////////////////////////////////////////////////////////////
	CREATEORJOIN_GROUP("8080", "TOCMULTICAST", "TOCMULTICASTGROUP")
	CREATEORJOIN_GROUP("8081", "TOCMULTICAST", "TOCMULTICASTGROUP")
	CREATEORJOIN_GROUP("8082", "TOCMULTICAST", "TOCMULTICASTGROUP")
	CREATEORJOIN_GROUP("8083", "TOCMULTICAST", "TOCMULTICASTGROUP")
	//////////////////////////////////////////////////////////////////////////////
	CREATEORJOIN_GROUP("8080", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8081", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8082", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8083", "TODMULTICAST", "TODMULTICASTGROUP")
	//////////////////////////////////////////////////////////////////////////////
	CREATEORJOIN_GROUP("8080", "BMULTICAST", "BMULTICASTGROUP")
	CREATEORJOIN_GROUP("8081", "BMULTICAST", "BMULTICASTGROUP")
	CREATEORJOIN_GROUP("8082", "BMULTICAST", "BMULTICASTGROUP")
	CREATEORJOIN_GROUP("8083", "BMULTICAST", "BMULTICASTGROUP")
}

func Test_start_group(t *testing.T) {
	STARTGROUPTEST("8080", "COMULTICASTGROUP")
	time.Sleep(2 * time.Second)
	STARTGROUPTEST("8081", "TOCMULTICASTGROUP")
	time.Sleep(2 * time.Second)
	STARTGROUPTEST("8082", "TODMULTICASTGROUP")
	time.Sleep(2 * time.Second)
	STARTGROUPTEST("8080", "BMULTICASTGROUP")
}

func Test_send_message(t *testing.T) {
	SENDMESSAGETEST("8081", "COMULTICASTGROUP")
	SENDMESSAGETEST("8082", "TOCMULTICASTGROUP")
	SENDMESSAGETEST("8081", "BMULTICASTGROUP")
	SENDMESSAGETEST("8083", "TODMULTICASTGROUP")
	SENDMESSAGETEST("8081", "TODMULTICASTGROUP")
	SENDMESSAGETEST("8085", "TOCMULTICASTGROUP")
}

func Test_get_message(t *testing.T) {
	GETMESSAGESTEST("8080", "MYGROUP")
}

func Test_get_groups(t *testing.T) {
	GETGROUPTEST("8080", "MYGROUP")
}
func Test_get_delivered_message(t *testing.T) {
	GETDELIVERMESSAGESTEST("8081", "MYGROUP")
}

func Test_Delete_group(t *testing.T) {
	DELETEGROUP("8081", "MYGROUP")
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

func STARTGROUPTEST(host string, group string) {
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

func GETGROUPTEST(host string, group string) {
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

func DELETEGROUP(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/groups/" + group
	method := "DELETE"

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

func SENDMESSAGETEST(host string, group string) {
	url := "http://localhost:" + host + "/multicast/v1/messaging/" + group
	method := "POST"
	m := []byte(test.RandomString(16))
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
