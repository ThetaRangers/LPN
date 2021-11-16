package cloud

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Item struct {
	Key   string
	Value []string `dynamodbav:"ValueList,omitempty"`
}

func SetupClient(region string) *dynamodb.DynamoDB {
	//Region taken from config
	//start session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	//create dynamo service client
	svc := dynamodb.New(sess)

	return svc
}

//PutItem Override value for existing key
func PutItem(svc *dynamodb.DynamoDB, tableName, key string, value [][]byte) {
	val := make([]string, 0)
	for i := 0; i < len(value); i++ {
		val = append(val, string(value[i]))
	}

	newItem := Item{key, val}
	av, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	log.Println("Successfully added " + newItem.Key)

}

func GetItem(svc *dynamodb.DynamoDB, tableName, key string) [][]byte {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	item := Item{}

	if result.Item == nil {
		log.Println("Could not find '" + key + "'")
		return [][]byte{}
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	log.Println("Item found")

	retValue := make([][]byte, 0)

	for i := 0; i < len(item.Value); i++ {
		retValue = append(retValue, []byte(item.Value[i]))
	}

	return retValue
}

func DeleteItem(svc *dynamodb.DynamoDB, tableName, key string) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}

	fmt.Println("Deleted '" + key + " from table " + tableName)
}

func AppendValue(svc *dynamodb.DynamoDB, tableName, key string, newValues [][]byte) {
	var dynamoValues []*dynamodb.AttributeValue
	length := len(newValues)

	if length == 0 {
		dynamoValues = []*dynamodb.AttributeValue{}
	}

	for i := 0; i < length; i++ {
		av := &dynamodb.AttributeValue{
			S: aws.String(string(newValues[i])),
		}
		dynamoValues = append(dynamoValues, av)
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newVal": {
				L: dynamoValues,
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET ValueList = list_append(ValueList, :newVal)"),
		TableName:        aws.String(tableName),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	log.Println("Appended completed")
}

/*func main() {

	client := SetupClient(utils.AwsRegion)
	key := "key"
	value := make([][]byte, 0)
	value = append(value, []byte("str1"))
	value = append(value, []byte("bho"))
	value = append(value, []byte("str2"))

	PutItem(client, utils.DynamoTable, key, value)

	retValue := GetItem(client, utils.DynamoTable, key)
	length := len(retValue)
	for i:=0; i< length; i++{
		fmt.Println(string(retValue[i]))
	}

	value1 := make([][]byte, 0)
	value1 = append(value1, []byte("at"))
	value1 = append(value1, []byte("the"))
	value1 = append(value1, []byte("end"))

	AppendValue(client, utils.DynamoTable, key, value1)

	retValue = GetItem(client, utils.DynamoTable, key)
	length = len(retValue)
	for i:=0; i< length; i++{
		fmt.Println(string(retValue[i]))
	}



}*/
