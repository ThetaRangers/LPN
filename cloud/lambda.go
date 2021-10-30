package cloud

import (
	"SDCC/utils"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/opentracing/opentracing-go/log"
	"strconv"
	_ "strconv"
)

type RequestNetwork struct {
	Ip string `json:"ip"`
	N  int    `json:"n"`
}

type ReplicaSet struct {
	Crashed int      `json:"crashed"`
	Valid   int      `json:"valid"`
	IpList  []string `json:"ipList"`
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

func RegisterToTheNetwork(ip string, n int, region string) ReplicaSet {
	client := setupClient(region)

	x := RequestNetwork{Ip: ip, N: n}

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

func main() {

	ret := RegisterToTheNetwork("test12", 2, "us-east-1")
	fmt.Printf("crashed %d valid %d list %s\n", ret.Crashed, ret.Valid, ret.IpList)
	/*TODO order to parse:
	-crashed:
		-1 -> ipList old replicas
		-0 -> check valid
	-valid:
		-1 -> ipList new replicas
		-0 -> ipList empty. Not enough nodes to get n replicas. Retry.
	*/
}
