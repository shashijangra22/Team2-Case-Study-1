package main

import (
	CustomerModels "Team2CaseStudy1/pkg/Customer/Models"
	OrderModels "Team2CaseStudy1/pkg/Order/Models"
	RestaurantModels "Team2CaseStudy1/pkg/Restaurant/Models"

	"Team2CaseStudy1/pkg/OrderProto/orderpb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var queryServiceClient orderpb.QueryServiceClient

func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Team 2": "Hello from Aadithya, Abhishek, Priya, Shashi!",
	})
}

// To get all the orders
func GetAllOrders(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetOrders(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get all the customers
func GetAllCustomers(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetCustomers(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get all the restaurants
func GetAllRestaurants(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetRestaurants(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get specific customer
func GetSpecificCustomer(c *gin.Context) {
	customerid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &orderpb.SpecificCustomerRequest{CustId: customerid}

	res, err := queryServiceClient.GetACustomer(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get specific order
func GetSpecificOrder(c *gin.Context) {
	orderid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &orderpb.SpecificOrderRequest{OrderId: orderid}

	res, err := queryServiceClient.GetAnOrder(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get specific Rest
func GetSpecificRest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &orderpb.SpecificRestaurantRequest{Id: id}

	res, err := queryServiceClient.GetARestaurant(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}


// To place a new order
func PostOrder(c *gin.Context) {

	body := c.Request.Body
	byteContent, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Sorry no content found: ", err.Error())
	}

	var NewOrder OrderModels.Order
	_ = json.Unmarshal(byteContent, &NewOrder)

	orderid := NewOrder.OrderId
	customerid := NewOrder.CustomerId
	restaurantid := NewOrder.RestaurantId
	itemlist := NewOrder.ItemLine
	price := NewOrder.Price
	discount := NewOrder.Discount

	var itemline []*orderpb.Item

	for i := range itemlist {
		itemline = append(itemline, &orderpb.Item{
			Name:  itemlist[i].Name,
			Price: itemlist[i].Price,
		})
	}

	req := &orderpb.OrderRequest{Ord: &orderpb.Order{
		OrderId:      orderid,
		CustomerId:   customerid,
		RestaurantId: restaurantid,
		ItemLine:     itemline,
		Price:        price,
		Discount:     discount,
	}}

	res, err := queryServiceClient.AddOrder(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})
}

// To add new customer
func PostCustomer(c *gin.Context) {

	body := c.Request.Body
	byteContent, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Sorry no content found: ", err.Error())
	}

	var NewCustomer CustomerModels.Customer

	_ = json.Unmarshal(byteContent, &NewCustomer)

	req := &orderpb.CustomerRequest{Cust: &orderpb.Customer{
		CustomerId: NewCustomer.CustomerId,
		Name:       NewCustomer.Name,
		Address:    NewCustomer.Address,
		Phone:      NewCustomer.Phone,
	}}

	res, err := queryServiceClient.AddCustomer(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})
}

// To place a new order
func PostRestaurant(c *gin.Context) {

	body := c.Request.Body
	byteContent, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Sorry no content found: ", err.Error())
	}

	var NewRest RestaurantModels.Rest
	err = json.Unmarshal(byteContent, &NewRest)

	if err != nil {
		fmt.Println("Sorry no content not in correct format: ", err.Error())
	}

	id := NewRest.ID
	name := NewRest.Name
	avail := NewRest.Availability
	itemlist := NewRest.Items
	rating := NewRest.Rating
	categ := NewRest.Category

	var itemline []*orderpb.Item

	for i := range itemlist {
		itemline = append(itemline, &orderpb.Item{
			Name:  itemlist[i].Name,
			Price: itemlist[i].Price,
		})
	}

	req := &orderpb.RestaurantRequest{Rest: &orderpb.Restaurant{
		Id: id,
		Name: name,
		Availability: avail,
		ItemLine: itemline,
		Rating: rating,
		Category: categ,
	}}

	res, err := queryServiceClient.AddRestaurant(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})

}

func main() {

	fmt.Println("hello from API INIT function")
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Sorry client cannot talk to server: %v", err)
	}

	defer conn.Close()

	queryServiceClient = orderpb.NewQueryServiceClient(conn)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiRouter := router.Group("/api")

	apiRouter.GET("/", GetIndex)
	apiRouter.GET("/orders", GetAllOrders)
	apiRouter.GET("/customers", GetAllCustomers)
	apiRouter.GET("/restaurants", GetAllRestaurants)
	apiRouter.GET("/order/:id", GetSpecificOrder)
	apiRouter.GET("/customer/:id", GetSpecificCustomer)
	apiRouter.GET("/restaurant/:id", GetSpecificRest)
	apiRouter.POST("/new-order", PostOrder)
	apiRouter.POST("/new-customer", PostCustomer)
	apiRouter.POST("/new-restaurant", PostRestaurant)

	router.Run("localhost:9001")
}
