package main

import (
	"context"
	"encoding/csv"
	"log"
	"lumel/models"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %s", err)
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
		_, err = ordersCollection.InsertOne(context.Background(), order)
		if err != nil {
			log.Fatalf("Failed to insert order: %s", err)
		}

		product := models.Product{
			ProductID:   productID,
			ProductName: record[3],
			Category:    record[4],
		}
		_, err = productsCollection.InsertOne(context.Background(), product)
		if err != nil {
			log.Fatalf("Failed to insert product: %s", err)
		}

		customer := models.Customer{
			CustomerID:      customerID,
			CustomerName:    record[12],
			CustomerEmail:   record[13],
			CustomerAddress: record[14],
			Region:          record[5],
		}
		_, err = customersCollection.InsertOne(context.Background(), customer)
		if err != nil {
			log.Fatalf("Failed to insert customer: %s", err)
		}
	}
}
