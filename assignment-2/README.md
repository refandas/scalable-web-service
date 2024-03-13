# Order Item Management 

This Go program manages orders and their items. It allows users to create, update, delete, and retrieve orders and their associated items. It demonstrates basic Gin framework, CRUD operations, and RESTful API design.

## Features

- Create a new order with items.
- Update an existing order, including adding or removing items.
- Delete an order and its associated items.
- Retrieve information about orders and their items.

## Usage

Install dependencies

```bash
go mod tidy
```

To run the program, use the following command:

```bash
go run main.go
```

## Database Schema
This project uses a PostgreSQL database to store orders and items. You can find the database schema in `schema.sql`.

## API Documentation
The API documentation is provided in `oas.yml` (OpenAPI Specification). You can use tools like Swagger UI to visualize and interact with the API.

## Configuration
Copy the `.env.example` file to `.env` and configure the database connection settings according to your environment.