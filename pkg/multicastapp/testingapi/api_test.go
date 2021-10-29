package test

import (
	"encoding/json"
	"fmt"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting/test"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

//PHASE 1 : creation of groups
func Test_create_a_group(t *testing.T) {
	//CREATEORJOIN_GROUP("8080", "BMULTICAST", "BMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8081", "BMULTICAST", "BMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8082", "BMULTICAST", "BMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8083", "BMULTICAST", "BMULTICASTGROUP")
	//	time.Sleep(5 * time.Second)
	//////////////////////////////////////////////////////////////////////////
	//CREATEORJOIN_GROUP("8080", "COMULTICAST", "COMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8081", "COMULTICAST", "COMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8082", "COMULTICAST", "COMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8083", "COMULTICAST", "COMULTICASTGROUP")
	//	time.Sleep(5 * time.Second)
	////////////////////////////////////////////////////////////////////////////
	//CREATEORJOIN_GROUP("8081", "TOCMULTICAST", "TOCMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8080", "TOCMULTICAST", "TOCMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8082", "TOCMULTICAST", "TOCMULTICASTGROUP")
	//CREATEORJOIN_GROUP("8083", "TOCMULTICAST", "TOCMULTICASTGROUP")
	//	time.Sleep(5 * time.Second)
	//////////////////////////////////////////////////////////////////////////////
	CREATEORJOIN_GROUP("8080", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8081", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8082", "TODMULTICAST", "TODMULTICASTGROUP")
	CREATEORJOIN_GROUP("8083", "TODMULTICAST", "TODMULTICASTGROUP")
	//	time.Sleep(5 * time.Second)
}

//PHASE 2 : One of the member start group
func Test_start_group(t *testing.T) {

	//STARTGROUPTEST("8082", "COMULTICASTGROUP")
	//time.Sleep(10 * time.Second)
	STARTGROUPTEST("8082", "TODMULTICASTGROUP")
	//	time.Sleep(5 * time.Second)
	//STARTGROUPTEST("8083", "TOCMULTICASTGROUP")
}

// PHASE 3 : start messaging
func Test_send_message(t *testing.T) {
	SENDMESSAGETEST("8083", "COMULTICASTGROUP")
	//	time.Sleep(3 * time.Second)
	//	SENDMESSAGETEST("8080", "TOCMULTICASTGROUP")
	//time.Sleep(10 * time.Second)
	//	SENDMESSAGETEST("8081", "TOCMULTICASTGROUP")
	//time.Sleep(10 * time.Second)
	//SENDMESSAGETEST("8080", "TOCMULTICASTGROUP")
	//time.Sleep(10 * time.Second)
	//	SENDMESSAGETEST("8081", "BMULTICASTGROUP")
	//	time.Sleep(3 * time.Second)
	//	SENDMESSAGETEST("8083", "TODMULTICASTGROUP")
	//	time.Sleep(3 * time.Second)
	//SENDMESSAGETEST("8081", "TOCMULTICASTGROUP")
	//	time.Sleep(3 * time.Second)
	//	SENDMESSAGETEST("8085", "TOCMULTICASTGROUP")
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
	DELETEGROUP("8081", "TODMULTICASTGROUP")
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

func Test_selector(t *testing.T) {
	ex1 := "192.168.96.3:90"
	splitted_strings := strings.Split(strings.Split(ex1, ":")[0], ".")
	log.Println(splitted_strings)
	log.Println(splitted_strings[len(splitted_strings)-1])
}
