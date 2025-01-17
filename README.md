# frappuccino 🥤

`frappuccino` - It is a RESTful API for managing inventory, menus, orders and reports for a coffee shop.

## Installation 🛠️

### Requirements

- [Go 1.23](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/docs/)
- [Docker](https://www.docker.com/get-started)

### Installation Steps

1. Clone the repository:

```bash
git clone git@github.com:ChivononGeko/frappuccino.git
cd frappuccino
```

2. Create a file `.env` in the root of the project and fill it with the following values:

```markdown
DB_HOST=<your_host>
DB_PORT=<your_db_port>
DB_USER=<your_db_user>
DB_PASSWORD=<your_db_password>
DB_NAME=<your_db_name>
API_PORT=<your_api_port>
```

3. Run `Docker Compose` to create and launch containers:

```bash
docker-compose up --build
```

## Strcuture 🗂️ 

```bash
frappuccino/
├── cmd/
│  └── main.go  
├── internal/
│  ├── apperrors/
│  │  └── apperrors.go
│  ├── db/
│  │  └── db.go
│  ├── handlers/ 			
│  │  ├── inventory_handler.go  
|  |  ├── menu_handler.go  
│  │  ├── order_handler.go  
│  │  ├── reports_handler.go  
│  │  └── utils.go  
│  ├── models/ 					
│  │  ├── customer.go  
│  │  ├── inventory_item_request.go  
│  │  ├── inventory_item.go  
│  │  ├── menu_item_request.go 
│  │  ├── menu_item.go  
│  │  ├── order_request.go 
│  │  ├── order.go  
│  │  ├── report.go 
│  │  └── utils.go  
│  ├── repositories/ 
│  │  ├── customer_repository.go  
│  │  ├── inventory_repository.go  
│  │  ├── menu_repository.go  
│  │  ├── order_repository.go  
│  │  ├── reports_repository.go 
│  │  └── utils.go  
│  ├── router/  
│  │  ├── inventory_router.go  
│  │  ├── menu_router.go  
│  │  ├── order_router.go  
│  │  ├── reports_router.go 
│  │  └── router.go 
│  └──services/ 
│     ├── inventory_service.go  
│     ├── menu_service.go  
|     ├── order_service.go  
│     └── report_service.go  
├── docker-compose.yml
├── Dockerfile
├── ERD.png
├── go.mod
├── go.sum
├── init.sql
├── insert.sql
└── README.md
```

## Using ⚙️

The API will be available at `http://localhost:8080`.

### Query examples

#### Inventory

- Create an inventory item:

  ```json
  POST /inventory
  Content-Type: application/json

  {
  "name": "Coffee Beans",
  "stock_level": 100,
  "price": 15.0,
  "unit_type": "kg"
  }
  ```

- Get all inventory items:

  ```json
  GET /inventory
  ```

- Get an inventory item by ID:

  ```json
  GET /inventory/{id}
  ```

- Update an inventory item:

  ```json
  PUT /inventory/{id}
  Content-Type: application/json

  {
  "name": "Coffee Beans",
  "stock_level": 150,
  "price": 15.0,
  "unit_type": "kg"
  }
  ```

- Delete an inventory item:

  ```json
  DELETE /inventory/{id}
  ```

- Get leftovers

  ```json
  GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}
  ```
  
  - `sortBy` (optional): Determines the sorting method. Can be either:
    - `price`: Sort by item price.
    - `quantity`: Sort by item quantity.
  - `page` (optional): Current page number, starting from 1.
  - `pageSize` (optional): Number of items per page. Default value: 10.


#### Menu

- Create a menu item:

  ```json
  POST /menu
  Content-Type: application/json

  {
    "name": "Espresso",
    "description": "Strong and bold coffee",
    "price": 3.50,
    "size": "small",
    "ingredients": [
      {
        "ingredient_id": "coffee_beans",
        "quantity": 0.02
      }
    ]
  }
  ```

- Get all the menu items:

  ```markwon
  GET /menu
  ```

- Get a menu item by ID:

  ```json
  GET /menu/{id}
  ```

- Update a menu item:

  ```json
  PUT /menu/{id}
  Content-Type: application/json

  {
    "name": "Espresso",
    "description": "Strong and bold coffee",
    "price": 3.50,
    "size": "small",
    "ingredients": [
      {
        "ingredient_id": "coffee_beans",
        "quantity": 0.02
      }
    ]
  }
  ```

- Delete a menu item:

  ```json
  DELETE /menu/{id}
  ```

#### Orders

- Create an order:

  ```json
  POST /orders
  Content-Type: application/json

  {
    "customer_name": "John Doe",
    "payment_method": "cash",
    "items": [
      {
        "product_id": "espresso",
        "quantity": 2
      }
    ],
    "instructions": {
    "note": "No sugar"
    }
  }
  ```

  -Create batch orders

  ```json
  POST /orders/batch-process
  Content-Type: application/json

  [
    {
      "customer_name": "John Doe",
      "payment_method": "cash",
      "items": [
        {
          "product_id": "espresso",
          "quantity": 2
        }
      ],
      "instructions": {
      "note": "No sugar"
      }
    },
    {
      "customer_name": "Emily Johnson",
      "payment_method": "cash",
      "items": [
        {
          "product_id": "latte",
          "quantity": 1
        }
      ],
      "instructions": {
      "note": "With caramel syrop"
      }
    },
    {
      "customer_name": "Sarah Brown",
      "payment_method": "cash",
      "items": [
        {
          "product_id": "americano",
          "quantity": 1
        },
        {
          "product_id": "chocolate_croissant",
          "quantity": 1
        }
      ],
     "instructions": {
      "note": "Warm up the croissant"
      }
    }
  ]
  ```

- Get all orders:

  ```json
  GET /orders
  ```

- Get an order by ID:

  ```json
  GET /orders/{id}
  ```

- Update the order:

  ```json
  PUT /orders/{id}
  Content-Type: application/json

  {
    "customer_name": "John Doe",
    "payment_method": "cash",
    "items": [
      {
        "product_id": "espresso",
        "quantity": 2
      }
    ],
    "instructions": {
    "note": "No sugar"
    }
  }
  ```

- Delete an order:

  ```json
  DELETE /orders/{id}
  ```

- Close the order:

  ```json
  POST /orders/{id}/close
  ```

- Number of ordered items
  
  ```json
  GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}
  ```

  - `startDate` (optional): The start date of the period in YYYY-MM-DD format.
  - `endDate` (optional): The end date of the period in YYYY-MM-DD format.


#### Reports

- Get a general sales report:

  ```json
  GET /reports/total-sales
  ```

- Get a report on popular products:

  ```json
  GET /reports/popular-items
  ```

- Menu and order search:

  ```json
  GET /reports/search?q=espresso&filter={orders|menu|all}&minPrice={minPrice}&maxPrice={maxPrice}
  ```

  - `q` (required): Search query string
  - `filter` (optional): What to search through, can be multiple values comma-separated:
      - `orders` (search in customer names and order details)
      - `menu` (search in item names and descriptions)
      - `all` (default, search everywhere)
  - `minPrice` (optional): Minimum order/item price to include
  - `maxPrice` (optional): Maximum order/item price to include


- Get a report on orders for the period:

  ```json
  GET /reports/orderedItemsByPeriod?period={day|month}&month={month}
  ```

  - `period` (required):
      - `day`: Groups data by day within the specified month.
      - `month`: Groups data by month within the specified year.
  - `month` (optional): Specifies the month (e.g., october). Used only if period=day.
  - `year` (optional): Specifies the year. Used only if `period`=`month`.

## License 📜

This project is licensed under the alem.school license

## Authors ✍🏻

- @maabylka
- @dausetov
