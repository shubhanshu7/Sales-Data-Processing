package models

import "time"

type Order struct {
	OrderID       string    `bson:"order_id"`
	ProductID     string    `bson:"product_id"`
	CustomerID    string    `bson:"customer_id"`
	DateOfSale    time.Time `bson:"date_of_sale"`
	QuantitySold  int       `bson:"quantity_sold"`
	UnitPrice     float64   `bson:"unit_price"`
	Discount      float64   `bson:"discount"`
	ShippingCost  float64   `bson:"shipping_cost"`
	PaymentMethod string    `bson:"payment_method"`
	TotalCost     float64   `bson:"total_cost"`
}

type Product struct {
	ProductID   string `bson:"product_id"`
	ProductName string `bson:"product_name"`
	Category    string `bson:"category"`
}

type Customer struct {
	CustomerID      string `bson:"customer_id"`
	CustomerName    string `bson:"customer_name"`
	CustomerEmail   string `bson:"customer_email"`
	CustomerAddress string `bson:"customer_address"`
	Region          string `bson:"region"`
}
