package test

import (
	"encoding/json"
	"fmt"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	CREATE_GROUP(t, "8081")
	CREATE_GROUP(t, "8082")
	CREATE_GROUP(t, "8083")
	CREATE_GROUP(t, "8080")
	time.Sleep(3 * time.Second)
	Test_STARTGROUP(t)
}

func CREATE_GROUP(t *testing.T, host string) {
	url := "http://localhost:" + host + "/multicast/v1/groups"
	method := "POST"

	payload := strings.NewReader(`{"multicast_type":"TOCMULTICAST","multicast_id":"TOCMULTICASTGROUP"}`)

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
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/TOCMULTICASTGROUP"
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

func Test_GETGROUP(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/PROVA"
	method := "GET"

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

type Message struct {
	Payload []byte
}

func SENDMESSAGE(t *testing.T, message string) {
	url := "http://localhost:8083/multicast/v1/messaging/PROVA"
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

func Test_SendMessage(t *testing.T) {
	url := "http://localhost:8083/multicast/v1/messaging/PROVA"
	method := "POST"
	m := []byte("PROVA")
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

func Test_GETMESSAGES(t *testing.T) {
	url := "http://localhost:8081/multicast/v1/messaging/PROVA"
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

func Test_src(t *testing.T) {
	src := "192.168.0.3"
	id := strings.Split(src, ".")
	id1 := strings.Join(id, "")
	ids, _ := strconv.Atoi(id1)
	utils.Myid = ids
	log.Println(ids)
}
