# Go PostgreSQL Project with Auto Table Creation

This Go project automatically creates tables in PostgreSQL on startup using configuration from a JSON file specified by the user.

## Tables
- customer_details
- product_details
- order_details
- sales_file_data

## How to Run

### 1. Set Up PostgreSQL Database

Before running the project, **create the database manually**.

```bash
psql -U postgres -c "CREATE DATABASE postgres;"
```
### 2. Mention the database credentials in dev.json file

```json
"database": {
      "type": "POSTGRES",
      "uri": "localhost",
      "name": "your-database-name(postgres)",
      "username": "your-username",
      "password": "your-password",
      "timeout": 90,
      "profile": true,
      "runScript": true
 }
```

### 3. Run with following Command

```bash
go run main.go -config=config/dev.json
```

### 4. Upload Sales Data

- **Endpoint:** `/salesdata/read`
- **Method:** `POST`
- **Description:** Uploads a CSV file. Data is inserted into the `sales_file_data` table.
- **Headers:**
  - `Content-Type: multipart/form-data`
- **Request:** Form-data with a file input named `files`
- **Response:**
```json
{
  "status":"S",
  "errMsg":"",
}
```

---

### 5. Get Total Revenue by Date

- **Endpoint:** `/revenue/total`
- **Method:** `GET`
- **Description:** Returns the total revenue between a given date range.
- **Headers:**
  - `start: YYYY-MM-DD`
  - `end: YYYY-MM-DD`
- **Response:**
```json
{
  "status":"S",
  "errMsg":"",
  "totalRevenue": 12345.67
}
```

---

### 6. Get Revenue by Product

- **Endpoint:** `/revenue/byProduct`
- **Method:** `POST`
- **Description:** Returns total revenue for a specific product between two dates.
- **Request Body:**
```json
{
  "productName": "Product A",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"
}
```
- **Response:**
```json
{
  "status":"S",
  "errMsg":"",
  "totalRevenue": 4567.89
}
```

---

### 7. Get Revenue by Category

- **Endpoint:** `/revenue/byCategory`
- **Method:** `POST`
- **Description:** Returns total revenue for a specific category between two dates.
- **Request Body:**
```json
{
  "category": "Electronics",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"
}
```
- **Response:**
```json
{
  "status":"S",
  "errMsg":"",
  "totalRevenue": 7890.12
}
```

---

### 8. Get Revenue by Region

- **Endpoint:** `/revenue/byRegion`
- **Method:** `POST`
- **Description:** Returns total revenue for a specific region between two dates.
- **Request Body:**
```json
{
  "region": "North",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"
}
```
- **Response:**
```json
{
  "status":"S",
  "errMsg":"",
  "totalRevenue": 3456.78
}
```

Refer config folder for example json