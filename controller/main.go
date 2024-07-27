package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"lumel/models"
	"os"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection

func main() {
	r := gin.Default()

	r.POST("/refresh", refreshSalesData)
	r.GET("/revenue", getTotalRevenue)
	r.GET("/revenue/product/:product_id", getTotalRevenueByProduct)

	r.Run(":8080")
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	ordersCollection = client.Database("sales").Collection("orders")
}

func refreshSalesData(c *gin.Context) {

	refreshData()
	fmt.Println("====data refreshing completed============")
	c.JSON(http.StatusOK, gin.H{"status": "Data refresh completed"})
}

func getTotalRevenue(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))
	fmt.Println("======startDate===", startDate)
	fmt.Println("======endDate===", endDate)

	filter := bson.M{"date_of_sale": bson.M{"$gte": startDate, "$lte": endDate}}
	cursor, err := ordersCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var totalRevenue float64
	for cursor.Next(context.Background()) {
		var order models.Order
		cursor.Decode(&order)
		totalRevenue += order.TotalCost
	}

	c.JSON(http.StatusOK, gin.H{"total_revenue": totalRevenue})
}
func getTotalRevenueByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))

	fmt.Println("======startDate===", startDate)
	fmt.Println("======endDate===", endDate)

	filter := bson.M{"product_id": productID, "date_of_sale": bson.M{"$gte": startDate, "$lte": endDate}}
	cursor, err := ordersCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var totalRevenue float64
	for cursor.Next(context.Background()) {
		var order models.Order
		cursor.Decode(&order)
		totalRevenue += order.TotalCost
	}

	c.JSON(http.StatusOK, gin.H{"total_revenue": totalRevenue})
}

func refreshData() {
	file, err := os.Open("new_sales_data.csv")
	if err != nil {
		log.Fatalf("File not found: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading CSV: %s", err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}
	defer client.Disconnect(context.Background())

	ordersCollection := client.Database("sales").Collection("orders")
	productsCollection := client.Database("sales").Collection("products")
	customersCollection := client.Database("sales").Collection("customers")

	for _, record := range records[1:] {
		orderID := record[0]
		productID := record[1]
		customerID := record[2]
		dateOfSale, _ := time.Parse("2006-01-02", record[6])
		quantitySold, _ := strconv.Atoi(record[7])
		unitPrice, _ := strconv.ParseFloat(record[8], 64)
		discount, _ := strconv.ParseFloat(record[9], 64)
		shippingCost, _ := strconv.ParseFloat(record[10], 64)
		paymentMethod := record[11]
		totalCost := (unitPrice * float64(quantitySold)) - discount + shippingCost

		order := models.Order{
			OrderID:       orderID,
			ProductID:     productID,
			CustomerID:    customerID,
			DateOfSale:    dateOfSale,
			QuantitySold:  quantitySold,
			UnitPrice:     unitPrice,
			Discount:      discount,
			ShippingCost:  shippingCost,
			PaymentMethod: paymentMethod,
			TotalCost:     totalCost,
		}

		_, err = ordersCollection.UpdateOne(context.Background(), bson.M{"order_id": orderID}, bson.M{"$set": order}, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatalf("Failed to upsert order: %s", err)
		}

		product := models.Product{
			ProductID:   productID,
			ProductName: record[3],
			Category:    record[4],
		}

		_, err = productsCollection.UpdateOne(context.Background(), bson.M{"product_id": productID}, bson.M{"$set": product}, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatalf("Failed to upsert product: %s", err)
		}

		customer := models.Customer{
			CustomerID:      customerID,
			CustomerName:    record[12],
			CustomerEmail:   record[13],
			CustomerAddress: record[14],
			Region:          record[5],
		}
		fmt.Println("customerid", customerID)

		_, err = customersCollection.UpdateOne(context.Background(), bson.M{"customer_id": customerID}, bson.M{"$set": customer}, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatalf("Failed to upsert customer: %s", err)
		}
	}

	log.Println("Data refresh completed successfully")
}