package models

type SalesFileData struct {
	OrderId         string  `gorm:"column:order_id"`
	ProductId       string  `gorm:"column:product_id"`
	CustomerId      string  `gorm:"column:customer_id"`
	ProductName     string  `gorm:"column:product_name"`
	Category        string  `gorm:"column:category"`
	Region          string  `gorm:"column:region"`
	DateOfSale      string  `gorm:"column:date_of_sale"`
	QuantitySold    int     `gorm:"column:quantity_sold"`
	UnitPrice       float64 `gorm:"column:unit_price"`
	Discount        float64 `gorm:"column:discount"`
	ShippingCost    float64 `gorm:"column:shipping_cost"`
	PaymentMethod   string  `gorm:"column:payment_method"`
	CustomerName    string  `gorm:"column:customer_name"`
	CustomerEmail   string  `gorm:"column:customer_email"`
	CustomerAddress string  `gorm:"column:customer_address"`
	UniqueId        string  `gorm:"column:unique_id"`
	CreatedDate     string  `gorm:"column:created_date"`
	CreatedBy       string  `gorm:"column:created_by"`
}

type ReadDataResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type CustomerDetails struct {
	CustomerId      string `gorm:"column:customer_id;primaryKey"`
	CustomerName    string `gorm:"column:customer_name"`
	CustomerEmail   string `gorm:"column:customer_email"`
	CustomerAddress string `gorm:"column:customer_address"`
	CreatedDate     string `gorm:"column:created_date"`
	CreatedBy       string `gorm:"column:created_by"`
}

type ProductDetails struct {
	ProductId   string  `gorm:"column:product_id;primaryKey"`
	ProductName string  `gorm:"column:product_name"`
	Category    string  `gorm:"column:category"`
	UnitPrice   float64 `gorm:"column:unit_price"`
	CreatedDate string  `gorm:"column:created_date"`
	CreatedBy   string  `gorm:"column:created_by"`
}

type OrderDetails struct {
	OrderId       string  `gorm:"column:order_id;primaryKey"`
	ProductId     string  `gorm:"column:product_id"`
	CustomerId    string  `gorm:"column:customer_id"`
	QuantitySold  int     `gorm:"quantity_sold"`
	UnitPrice     float64 `gorm:"unit_price"`
	Discount      float64 `gorm:"discount"`
	Region        string  `json:"region"`
	ShippingCost  float64 `gorm:"shipping_cost"`
	PaymentMethod string  `gorm:"payment_method"`
	DateOfSale    string  `gorm:"date_of_sale"`
	CreatedDate   string  `gorm:"column:created_date"`
	CreatedBy     string  `gorm:"column:created_by"`
}

type CustomerReq struct {
	ReqArr []CustomerDetails `json:"reqArr"`
	User   string            `json:"user"`
}

type OrderReq struct {
	ReqArr []OrderDetails `json:"reqArr"`
	User   string         `json:"user"`
}

type ProductReq struct {
	ReqArr []ProductDetails `json:"reqArr"`
	User   string           `json:"user"`
}

type Config struct {
	Port    int      `json:"port"`
	Env     string   `json:"env"`
	Version string   `json:"version"`
	Profile bool     `json:"profile"`
	DB      DBStruct `json:"database"`
}

type DBStruct struct {
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Timeout   int    `json:"timeout"`
	Profile   bool   `json:"profile"`
	RunScript bool   `json:"runScript"`
}

type RevenueRespStruct struct {
	TotRev float64 `json:"totRev"`
	ErrMsg string  `json:"errMsg"`
	Status string  `json:"status"`
}

type ProductRevReq struct {
	ProductName string `json:"productName"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type CategoryRevReq struct {
	Category  string `json:"category"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type RegionRevReq struct {
	Region    string `json:"region"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
