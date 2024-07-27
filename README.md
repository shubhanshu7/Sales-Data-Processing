# Sales Data Processing 

## Introduction

This project processes and analyzes sales data from a CSV file, stores it in a MongoDB database, and provides a RESTful API to trigger data refresh and perform various analyses.

## Prerequisites

- Go 1.16 or later
- MongoDB 4.4 or later
- Gin web framework
- CSV file with sales data 
      

### Running the Application

1. **Ensure MongoDB is running:**

    mongod

2. **Run the application:**

    For running go to the following paths and run the executables
    These are different exe's for all the requirements
    
    1.DataLoading  (Data Loading from csv file)
        dataloading/loadData.exe
        ## please put csv file in this path before running the exe

    2.RefreshMechanism (Data Refresh Mechanism)
        refreshMechanism/refreshMechanism.exe

    3.API Server (RESTful API for Analysis)
        controller/server.exe
    

### API Endpoints

1. Trigger Data Refresh

    curl -X POST http://localhost:8080/refresh

2. Retrieve Total Revenue

   curl -X GET "http://localhost:8080/revenue" -H "Content-Type: application/json" -d '{
    "start_date": "2023-01-01",
    "end_date": "2023-12-31"
     }'

3.  Retrieve Total Revenue by Product

   curl -X GET "http://localhost:8080/revenue/product" -H "Content-Type: application/json" -d '{
    "start_date": "2023-01-01",
    "end_date": "2023-12-31",
    "product_id": "P123"
}'
