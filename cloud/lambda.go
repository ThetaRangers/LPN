package cloud

import (
	pb "SDCC/registerServer"
	"SDCC/utils"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"google.golang.org/grpc"
	"log"
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
		log.Fatalln(err)
	}

	//TODO change name of the lambda function
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("RegService"), Payload: payload})
	if err != nil {
		log.Fatalln(err)
	}

	b := result.Payload
	resp := ReplicaSet{}
	if err := json.Unmarshal(b, &resp); err != nil {
		panic(err)
	}

	return resp
}

func RegisterStub(ip string, ipStr string, n int, region string) ReplicaSet {
	addr := utils.TestingServer + ":50052"

	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	c := pb.NewOperationsClient(conn)

	res, _ := c.Register(context.Background(), &pb.RegisterMessage{Ip: ip, NodeId: ipStr})
	crashed := res.GetCrashed()
	addresses := res.GetAddresses()
	nodeIds := res.GetNodeIdS()

	i := 0
	if crashed {
		i = 1
	}

	ipList := make([]IpStruct, 0)

	for index, k := range addresses {
		ipList = append(ipList, IpStruct{Ip: k, IpString: nodeIds[index]})
	}

	r := ReplicaSet{Valid: 0, Crashed: i, IpList: ipList}
	return r
}
