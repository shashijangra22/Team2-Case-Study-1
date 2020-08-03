package Services

import (
	RestaurantModels "Team2CaseStudy1/pkg/Restaurant/Models"

	"Team2CaseStudy1/pkg/OrderProto/orderpb"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func FetchRestaurantTable(db *dynamodb.DynamoDB) []*orderpb.Restaurant {

	var allRestaurants []*orderpb.Restaurant

	// Create the Expression to fill the input struct with.
	filt := expression.Name("ID").GreaterThan(expression.Value(0))

	proj := expression.NamesList(expression.Name("ID"),
		expression.Name("Name"),
		expression.Name("Availability"),
		expression.Name("Items"),
		expression.Name("Category"),
		expression.Name("Rating"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Got error building expression for getting all Restaurants")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("T2-Restaurants"),
	}

	// Make the DynamoDB Query API call
	result, err := db.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed for Restaurent table fetching")
		fmt.Println((err.Error()))
		os.Exit(1)
	}

	for _, i := range result.Items {
		restaurantItem := RestaurantModels.Rest{}
		var itemLine []*orderpb.Item

		err = dynamodbattribute.UnmarshalMap(i, &restaurantItem)

		if err != nil {
			fmt.Println("Got error unmarshalling restaurant table")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, item := range restaurantItem.Items {
			itemLine = append(itemLine, &orderpb.Item{Name: item.Name, Price: item.Price})
		}

		allRestaurants = append(allRestaurants, &orderpb.Restaurant{
			Id: restaurantItem.ID,
			Name: restaurantItem.Name,
			Category: restaurantItem.Category,
			Availability: restaurantItem.Availability,
			ItemLine: itemLine,
			Rating: restaurantItem.Rating})
	}

	return allRestaurants

}


func GetSpecificRestaurantDetails(db *dynamodb.DynamoDB, Id int64) *orderpb.Restaurant {

	keyCond := expression.Key("ID").Equal(expression.Value(Id))
	proj := expression.NamesList(expression.Name("ID"),
		expression.Name("Name"),
		expression.Name("Availability"),
		expression.Name("Items"),
		expression.Name("Category"),
		expression.Name("Rating"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression for getting specific Customer")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("T2-Restaurants"),
	}

	res, err := db.Query(params)

	var restDetails []RestaurantModels.Rest

	_ = dynamodbattribute.UnmarshalListOfMaps(res.Items, &restDetails)

	var rest *orderpb.Restaurant

	if len(restDetails) == 0 {
		rest = &orderpb.Restaurant{}
	} else {
		restaurantItem := restDetails[0]
		var itemline []*orderpb.Item
		for _, item := range restaurantItem.Items {
			itemline = append(itemline, &orderpb.Item{Name: item.Name, Price: item.Price})
		}
		rest = &orderpb.Restaurant{
			Id: restaurantItem.ID,
			Name: restaurantItem.Name,
			Category: restaurantItem.Category,
			Availability: restaurantItem.Availability,
			ItemLine: itemline,
			Rating: restaurantItem.Rating}
	}
	return rest
}

func AddRstDetails(db *dynamodb.DynamoDB, rest RestaurantModels.Rest) {

	orderDynAttr, err := dynamodbattribute.MarshalMap(rest)

	if err != nil {
		panic("Cannot map the values given in Order struct for post request...")
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("T2-Restaurants"),
		Item:      orderDynAttr,
	}

	_, err = db.PutItem(params)

	if err != nil {
		panic("Error in putting the rest item")
	}

}
