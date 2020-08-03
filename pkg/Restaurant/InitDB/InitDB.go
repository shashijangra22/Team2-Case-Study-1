package InitDB

import (
	"Team2CaseStudy1/pkg/Err"
	"Team2CaseStudy1/pkg/Restaurant/Models"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"io"
	"log"
	"os"
	"strconv"
)


func parseRecord(record []string) Models.Rest {
	ID, _ := strconv.ParseInt(record[0], 10, 64)
	Name := record[1]
	Availablity := true
	if record[2] == "Close" {
		Availablity = false
	}
	itemName := record[3]
	Price, _ := strconv.ParseFloat(record[4], 64)
	Rating, _ := strconv.ParseFloat(record[5], 64)
	Category := record[6]
	itemObj := Models.Item{itemName, float32(Price)}
	restObj := Models.Rest{ID, Name, Availablity, []Models.Item{itemObj}, float32(Rating),Category}
	return restObj
}

func clubRecords(records [][]string) []Models.Rest {
	var clubbedRecords []Models.Rest
	prev := "INF"
	var restObj Models.Rest
	var itemObj Models.Item
	flag := false
	for _, record := range records {
		tempObj := parseRecord(record)
		if record[0] != prev {
			if flag == true {
				clubbedRecords = append(clubbedRecords, restObj)
			}
			flag = true
			restObj = tempObj
			prev = record[0]
		} else {
			itemObj = tempObj.Items[0]
			restObj.Items = append(restObj.Items, itemObj)
		}
	}
	clubbedRecords = append(clubbedRecords, restObj)
	fmt.Printf("Records after clubbing: %v\n", len(clubbedRecords))
	return clubbedRecords
}

func readCSV(filename string) [][]string {
	csvFile, err := os.Open(filename)
	Err.CheckError(err)
	r := csv.NewReader(csvFile)
	_, _ = r.Read()
	var records [][]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		Err.CheckError(err)
		records = append(records, record)
	}
	fmt.Printf("Records processed: %v\n", len(records))
	return records
}


func initDB(db *dynamodb.DynamoDB)  {

	awsParams := &dynamodb.CreateTableInput{
		TableName: aws.String("T2-Restaurants"),

		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("ID"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("Name"), KeyType: aws.String("RANGE")},
		},

		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("ID"), AttributeType: aws.String("N")},
			{AttributeName: aws.String("Name"), AttributeType: aws.String("S")},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	response, err := db.CreateTable(awsParams)
	if err != nil {
		log.Fatalf("Sorry error creating table : %v", err)
	}
	// print the response
	fmt.Println(response)
}



func toDB(db *dynamodb.DynamoDB,restaurants []Models.Rest, filename string){


	for i, r := range restaurants{

		restMap, err := dynamodbattribute.MarshalMap(r)

		if err != nil {
			panic("Cannot map the values given in music struct... ")
		}

		params := &dynamodb.PutItemInput{
			TableName: aws.String("T2-Restaurants"),
			Item:      restMap,
		}


		resp, err := db.PutItem(params)

		if err != nil {
			fmt.Printf("While inserting %v th record following error occured %v", i, err.Error())
			panic(err.Error())
		}

		fmt.Printf("Record %v is Inserted... \n", i)
		fmt.Println(resp)

	}
}

func InitRestaurant(db *dynamodb.DynamoDB,filename string) {
	//var filename string
	//filename := "/Users/aadithya.md/Desktop/bootcamp/training/workspace/go-basics/casestudy1/c2/rest.csv"
	fmt.Println("Reading " + filename)
	records := readCSV(filename)
	restaurants := clubRecords(records)
	initDB(db)
	toDB(db,restaurants, filename)
}
