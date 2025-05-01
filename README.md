
# API Documentation

This section provides details on the available API endpoints in the Go PostgreSQL Project.

---

## 1. Upload Sales Data

- **Endpoint:** `/salesdata/read`
- **Method:** `POST`
- **Description:** Uploads a CSV file. Data is inserted into the `sales_file_data` table.
- **Headers:**
  - `Content-Type: multipart/form-data`
- **Request:** Form-data with a file input named `file`
- **Response:**
```json
{
  "status": "success",
  "message": "File uploaded and data inserted"
}
```

---

## 2. Get Total Revenue by Date

- **Endpoint:** `/revenue/total`
- **Method:** `GET`
- **Description:** Returns the total revenue between a given date range.
- **Headers:**
  - `start: YYYY-MM-DD`
  - `end: YYYY-MM-DD`
- **Response:**
```json
{
  "totalRevenue": 12345.67
}
```

---

## 3. Get Revenue by Product

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
  "productName": "Product A",
  "totalRevenue": 4567.89
}
```

---

## 4. Get Revenue by Category

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
  "category": "Electronics",
  "totalRevenue": 7890.12
}
```

---

## 5. Get Revenue by Region

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

---

## Notes

- Dates should be in `YYYY-MM-DD` format.
- Ensure correct content types for all requests.
