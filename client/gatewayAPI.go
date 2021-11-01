package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	//TODO set as terraform output result
	url = ""
)

type IpListStruct struct {
	IpList []string `json:"ipList"`
}

func GetAllNodes() []string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

	//Close Body after read
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(response.Body)

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
	}

	ipList := IpListStruct{}

	if err := json.Unmarshal(b, &ipList); err != nil {
		panic(err)
	}

	return ipList.IpList
}

func main() {

	fmt.Println(GetAllNodes())

}
