package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Key   string
	Value []string `dynamodbav:"ValueList,omitempty"`
}

func SetupClient(region string) (*dynamodb.DynamoDB, error) {
	//Region taken from config
	//start session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return nil, err
	}

	//create dynamo service client
	svc := dynamodb.New(sess)

	return svc, nil
}

//PutItem Override value for existing key
func PutItem(svc *dynamodb.DynamoDB, tableName, key string, value [][]byte) error {
	val := make([]string, 0)
	for i := 0; i < len(value); i++ {
		val = append(val, string(value[i]))
	}

	newItem := Item{key, val}
	av, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetItem(svc *dynamodb.DynamoDB, tableName, key string) ([][]byte, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return [][]byte{}, err
	}

	item := Item{}

	if result.Item == nil {
		return [][]byte{}, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return [][]byte{}, err
	}

	retValue := make([][]byte, 0)

	for i := 0; i < len(item.Value); i++ {
		retValue = append(retValue, []byte(item.Value[i]))
	}

	return retValue, nil
}

func DeleteItem(svc *dynamodb.DynamoDB, tableName, key string) error {
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
		return err
	}

	return nil
}

func AppendValue(svc *dynamodb.DynamoDB, tableName, key string, newValues [][]byte) error {
	var dynamoValues []*dynamodb.AttributeValue
	length := len(newValues)

	if length == 0 {
		dynamoValues = []*dynamodb.AttributeValue{}
		//avoid dynamo update
		return nil
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
	if err!= nil {
		return err
	}
	/*for err != nil {
		_, err = svc.UpdateItem(input)
	}*/


	return nil
}
