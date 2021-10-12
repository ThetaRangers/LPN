package cloud
import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type Item struct {
	Key string
	Value []string
}


func exitErrorf(msg string/*, args ...interface{}*/) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}

func SetupClient() *dynamodb.DynamoDB{
	//TODO credentials as env variables
	//start session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	//create S3 service client
	svc := dynamodb.New(sess)

	return svc
}

func listTables(svc *dynamodb.DynamoDB){
	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to aws error.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

//key value table
func createTable(svc *dynamodb.DynamoDB, tableName string ){
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
}

//Override value for existing key
func PutItem(svc *dynamodb.DynamoDB, tableName, key string, value []string ){
	newItem := Item{key, value}
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

	fmt.Print("Successfully added ( " + newItem.Key)
	for i:=0; i< len(newItem.Value); i++  {
		fmt.Print("'" + newItem.Value[i]+"' ")
	}
	fmt.Println(") to table " + tableName)
}

func GetItem(svc *dynamodb.DynamoDB, tableName, key string) (Item, bool){
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
		fmt.Println( "Could not find '" + key + "'")
		return item, false
	}



	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println("Found item:")
	fmt.Println("Key:  ", item.Key)
	fmt.Println("Value: ", item.Value)

	return item, true
}

func DeleteItem(svc *dynamodb.DynamoDB, tableName, key string){
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

func AppendValue(svc *dynamodb.DynamoDB, tableName, key string, newValues []string){
	item, ret := GetItem(svc, tableName, key)
	if ret {
		item.Value = append(item.Value, newValues...)
		PutItem(svc, tableName, key, item.Value)
	}
}


/*func main(){
	svc := SetupClient()
	createTable(svc, "testTable")
	listTables(svc)
	PutItem(svc, "testTable", "put", []string{"g", "h", "i"})
	GetItem(svc, "testTable", "put")
	AppendValue(svc, "testTable", "put", []string{"l", "m", "n"})
	GetItem(svc, "testTable", "put")
	DeleteItem(svc, "testTable", "put")

}*/
