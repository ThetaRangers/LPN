package cloud

import (
	"SDCC/utils"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/opentracing/opentracing-go/log"
	"strconv"
)

type RequestNetwork struct {
	Ip    string `json:"ip"`
	N     int    `json:"n"`
	IpStr string `json:"ipStr"`
}

type IpStruct struct {
	Ip       string `json:"ip"`
	IpString string `json:"strIp"`
}

type ReplicaSet struct {
	Crashed int        `json:"crashed"`
	Valid   int        `json:"valid"`
	IpList  []IpStruct `json:"ipList"`
}

func setupClient(region string) *lambda.Lambda {
	//Region taken from config
	//start session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	//create lambda service client
	svc := lambda.New(sess)

	return svc
}

func RegisterToTheNetwork(ip string, ipStr string, n int, region string) ReplicaSet {
	client := setupClient(region)

	x := RequestNetwork{Ip: ip, IpStr: ipStr, N: n}

	payload, err := json.Marshal(&x)
	if err != nil {
		log.Error(err)
	}

	//TODO change name of the lambda function
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("RegService"), Payload: payload})
	if err != nil {
		log.Error(err)
	}

	b := []byte(result.Payload)
	resp := ReplicaSet{}
	if err := json.Unmarshal(b, &resp); err != nil {
		panic(err)
	}

	return resp
}

func RegisterStub(ip, network string, n int, region string) []string {
	address := "172.17.0."
	set := make([]string, 0)
	for i := 2; i < utils.Replicas+3; i++ {
		abba := strconv.Itoa(i)
		tmpAddr := address + abba
		if tmpAddr != ip {
			set = append(set, tmpAddr)
		}
	}
	return set
}

func RegisterStub2(ip, network string, n int, region string) ([]string, bool, bool) {
	address := "172.17.0."
	set := make([]string, 0)

	if ip == "172.17.0.2" {
		for i := 3; i < utils.Replicas+3; i++ {
			abba := strconv.Itoa(i)
			tmpAddr := address + abba

			set = append(set, tmpAddr)
		}
	} else if ip == "172.17.0.11" {
		set = []string{"172.17.0.7", "172.17.0.8", "172.17.0.9", "172.17.0.10"}
	}

	if ip == "172.17.0.7" || ip == "172.17.0.8" || ip == "172.17.0.9" || ip == "172.17.0.10" {
		return append(set, "172.17.0.2"), true, false
	}

	return set, false, false
}

func main() {
	var ip string
	var ipStr string
	for i := 7; i < 10; i++ {
		ip = fmt.Sprintf("ip%d", i)
		ipStr = fmt.Sprintf("ip%dstr", i)
		ret := RegisterToTheNetwork(ip, ipStr, 2, "us-east-1")
		fmt.Printf("ip %s ipStr %s\n", ip, ipStr)
		fmt.Printf("crashed %d valid %d list\n", ret.Crashed, ret.Valid)
		for j := 0; j < len(ret.IpList); j++ {
			r := ret.IpList[j]
			fmt.Printf("repIp %s RepIpString %s\n", r.Ip, r.IpString)
		}
	}
}
