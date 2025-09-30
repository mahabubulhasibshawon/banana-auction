
---

# Auction System

Welcome to the Banana Auction System, a RESTful API application designed to facilitate the creation, management, and bidding on banana lots. This project allows sellers to list lots, start auctions, and buyers to place bids, with role-based access control (RBAC) enforced via JWT authentication.

## Tech Stack

- **Language**: Go (1.24)
- **Database**: PostgreSQL
- **Dependencies**:
  - `github.com/golang-jwt/jwt/v5` for JWT authentication
  - `github.com/joho/godotenv` for environment variable management
  - `github.com/lib/pq` for PostgreSQL driver
  - `golang.org/x/crypto` for password hashing
- **Configuration**: Environment variables loaded via `.env` file
- **HTTP Server**: Built with `net/http`
- **Development Tools**: Go modules for dependency management

## Prerequisites

- Go 1.24 or higher
- PostgreSQL installed and running
- Git for cloning the repository

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/banana-auction.git
   cd banana-auction
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following variables:
   ```
   VERSION=1.0.0
   SERVICE_NAME=banana-auction
   HTTP_PORT=8080
   JWT_SECRET_KEY=your-secure-jwt-secret-key
   JWT_REFRESH_KEY=your-secure-jwt-refresh-key
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=bananaauction
   ```

4. Set up the database:
   - Create a database named `bananaauction` in PostgreSQL.
   - The application will create the necessary tables (`users`, `lots`, `auctions`, `bids`) on startup.

5. Run the application:
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8080`.

## Endpoints

All endpoints except `/signup` and `/login` require a valid JWT token in the `Authorization` header (e.g., `Bearer <token>`). Use the `/login` endpoint to obtain a token.

### Authentication Endpoints

- **Signup**
  - **Method**: `POST`
  - **URL**: `/signup`
  - **Description**: Register a new user (seller or buyer).
  - **Request Payload**:
    ```json
    {
      "username": "testuser",
      "password": "password123",
      "name": "Test User",
      "role": "seller"
    }
    ```
  - **Response** (Success, 201 Created):
    ```json
    {
      "id": 1
    }
    ```
  - **Response** (Failure, 400 Bad Request):
    ```json
    {
      "error": "Username, password, name, and role are required"
    }
    ```
  - **Response** (Failure, 409 Conflict):
    ```json
    {
      "error": "Username already exists"
    }
    ```

- **Login**
  - **Method**: `POST`
  - **URL**: `/login`
  - **Description**: Authenticate a user and receive a JWT token.
  - **Request Payload**:
    ```json
    {
      "username": "testuser",
      "password": "password123"
    }
    ```
  - **Response** (Success, 200 OK):
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
  - **Response** (Failure, 401 Unauthorized):
    ```json
    {
      "error": "Invalid username or password"
    }
    ```

### Lot Management Endpoints (Seller Only)

- **Create Lot**
  - **Method**: `POST`
  - **URL**: `/lots`
  - **Description**: Create a new banana lot.
  - **Request Payload**:
    ```json
    {
      "cultivar": "Cavendish",
      "planted_country": "Ecuador",
      "harvest_date": "2025-10-01",
      "total_weight_kg": 1500
    }
    ```
  - **Response** (Success, 201 Created):
    ```json
    {
      "id": 1
    }
    ```
  - **Response** (Failure, 400 Bad Request):
    ```json
    {
      "error": "Minimum weight allowed is 1000 kg"
    }
    ```

- **Update Lot**
  - **Method**: `PATCH`
  - **URL**: `/lots/{id}`
  - **Description**: Update the harvest date of a lot (seller-owned only).
  - **Request Payload**:
    ```json
    {
      "harvest_date": "2025-10-02"
    }
    ```
  - **Response** (Success, 200 OK): No content.
  - **Response** (Failure, 404 Not Found):
    ```json
    {
      "error": "Lot not found"
    }
    ```

- **Delete Lot**
  - **Method**: `DELETE`
  - **URL**: `/lots/{id}`
  - **Description**: Delete a lot and its associated auctions/bids (seller-owned only).
  - **Response** (Success, 204 No Content): No content.
  - **Response** (Failure, 404 Not Found):
    ```json
    {
      "error": "Lot not found"
    }
    ```

- **List Lots**
  - **Method**: `GET`
  - **URL**: `/lots`
  - **Description**: List all lots (seller-owned only).
  - **Response** (Success, 200 OK):
    ```json
    [
      {
        "id": 1,
        "seller_id": 1,
        "cultivar": "Cavendish",
        "planted_country": "Ecuador",
        "harvest_date": "2025-10-01",
        "total_weight_kg": 1500
      }
    ]
    ```
  - **Response** (Failure, 403 Forbidden):
    ```json
    {
      "error": "Only sellers can list lots"
    }
    ```

### Auction Management Endpoints (Seller Only)

- **Create Auction**
  - **Method**: `POST`
  - **URL**: `/auctions`
  - **Description**: Start a new auction for a lot.
  - **Request Payload**:
    ```json
    {
      "lot_id": 1,
      "start_date": "2025-10-01",
      "duration_days": 7,
      "initial_price_per_kg": 0.5
    }
    ```
  - **Response** (Success, 201 Created):
    ```json
    {
      "id": 1
    }
    ```
  - **Response** (Failure, 400 Bad Request):
    ```json
    {
      "error": "Auction already exists for this lot"
    }
    ```

- **List Bids**
  - **Method**: `GET`
  - **URL**: `/auctions/{id}/bids`
  - **Description**: List all bids for an auction (seller-owned lot only).
  - **Response** (Success, 200 OK):
    ```json
    [
      {
        "id": 1,
        "auction_id": 1,
        "buyer_id": 2,
        "bid_price_per_kg": 0.6
      }
    ]
    ```
  - **Response** (Failure, 404 Not Found):
    ```json
    {
      "error": "Auction not found"
    }
    ```

### Bid Management Endpoints (Buyer Only)

- **Create Bid**
  - **Method**: `POST`
  - **URL**: `/auctions/{id}/bids`
  - **Description**: Place a bid on an auction.
  - **Request Payload**:
    ```json
    {
      "bid_price_per_kg": 0.6
    }
    ```
  - **Response** (Success, 201 Created):
    ```json
    {
      "id": 1
    }
    ```
  - **Response** (Failure, 404 Not Found):
    ```json
    {
      "error": "Auction not found"
    }
    ```

![alt text](image-1.png)
## Relationships

 * users to lots: One-to-many relationship where a user (seller) can create multiple lots (seller_id FK).  
 * lots to auctions: One-to-many relationship where a lot can have multiple auctions (lot_id FK), though the application enforces a single active auction per lot via    
 * AuctionExistsForLot.  
 * auctions to bids: One-to-many relationship where an auction can have multiple bids (auction_id FK).  
 * users to bids: One-to-many relationship where a user (buyer) can place multiple bids (buyer_id FK).  


## Usage

1. **Start the Server**: Run `go run main.go` after setting up the `.env` file.
2. **Test Endpoints**: Use a tool like `curl`, Postman, or a similar API client.
   - Example login request:
     ```bash
     curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"testuser","password":"password123"}'
     ```
   - Example protected request with token:
     ```bash
     curl -X POST http://localhost:8080/lots -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." -H "Content-Type: application/json" -d '{"cultivar":"Cavendish","planted_country":"Ecuador","harvest_date":"2025-10-01","total_weight_kg":1500}'
     ```

## Project Structure

```
banana-auction/
├── cmd/                # Main application entry point
├── config/             # Configuration loading
├── database/           # Database operations
├── rest/               # REST API handlers, middlewares, and routes
├── utils/              # Utility functions (e.g., JWT, password hashing)
├── models/             # Data models
├── go.mod              # Go module file
└── .env                # Environment variables
```

## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit changes (`git commit -m "Add feature"`).
4. Push to the branch (`git push origin feature-name`).
5. Open a Pull Request.


---