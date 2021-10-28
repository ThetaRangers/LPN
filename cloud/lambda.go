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
	"strings"
)

type RequestNetwork struct {
	Ip      string `json:"ip"`
	Network string `json:"network"`
	N       int    `json:"n"`
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

func RegisterToTheNetwork(ip, network string, n int, region string) []string {
	client := setupClient("region")

	x := RequestNetwork{Ip: ip, Network: network, N: n}

	payload, err := json.Marshal(&x)
	if err != nil {
		log.Error(err)
	}

	//TODO change name of the lambda function
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("registryService"), Payload: payload})
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(result.Payload))
	res := strings.Replace(string(result.Payload[1:len(result.Payload)-1]), "\\", "", -1)

	var dat map[string][]string
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		panic(err)
	}

	return dat["results"]

}

func RegisterStub(ip, network string, n int, region string) []string {
	address := "172.17.0."
	set := make([]string, 0)
	for i := 2; i < utils.Replicas + 3; i++ {
		abba := strconv.Itoa(i)
		tmpAddr := address + abba
		if tmpAddr != ip {
			set = append(set, tmpAddr)
		}
	}
	return set
}

func main() {

	fmt.Println(RegisterToTheNetwork("test12", "tabellone", 2, "us-east-1"))
}
