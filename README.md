# **README**

## **Description**

This App is a simple gRPC-based application that allows users to create, read, update, and delete orders. The app uses a database to store order information and provides a gRPC interface for clients to interact with the database.

## **Installation**

To install and run the app, follow these steps:

1. Clone the repository: `git clone https://github.com/glebkasilov/Order-grpc`
2. Change into the repository directory: `cd application`
3. Build the app: `make build`
4. Run the app: `make run`

## **Usage**

To use the app, you can use a gRPC client to interact with the gRPC interface. The app provides the following endpoints:

- `CreateOrder`: Create a new order
- `GetOrder`: Get an existing order by ID
- `UpdateOrder`: Update an existing order
- `DeleteOrder`: Delete an existing order
- `ListOrders`: List all existing orders

## **Authors**

- Gleb Kasilov
