package model

// ==================== PPOB Models ====================

type PPOBOrder struct {
	Id            string             `json:"id" validate:"required"`
	ReferenceId   string             `json:"reference_id" validate:"required"`
	Customer      PPOBCustomer       `json:"customer" validate:"required"`
	Product       PPOBProduct        `json:"product"`
	Category      PPOBCategory       `json:"category"`
	Payment       PPOBPayment        `json:"payment" validate:"required"`
	CommissionFee PPOBCommissionFee  `json:"commission_fee"`
	Amount        int                `json:"amount" validate:"required"`
	ServiceFee    int                `json:"service_fee"`
	PaymentFee    int                `json:"payment_fee"`
	TotalTax      int                `json:"total_tax"`
	TotalFee      int                `json:"total_fee"`
	Total         int                `json:"total"`
	PaymentStatus string             `json:"payment_status"`
	OrderStatus   string             `json:"order_status"`
	ExpiredAt     string             `json:"expired_at"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
	Source        string             `json:"source"`
	Details       []PPOBOrderDetail  `json:"details"`
	History       []PPOBOrderHistory `json:"history"`
	Metadata      interface{}        `json:"metadata"`
	MetadataRaw   string             `json:"-"` // Raw JSON string for storage
}

type PPOBCustomer struct {
	Id    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Type  string `json:"type" validate:"required"`
	Phone string `json:"phone"`
}
type PPOBProduct struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type PPOBCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type PPOBPayment struct {
	ReferenceId    string      `json:"reference_id" validate:"required"`
	Channel        string      `json:"channel" validate:"required"`
	PaymentCode    string      `json:"payment_code" validate:"required"`
	PaymentGateway string      `json:"payment_gateway" validate:"required"`
	Metadata       interface{} `json:"metadata"`
}

type PPOBCommissionFee struct {
	Jatis      int `json:"jatis"`
	Aggregator int `json:"aggregator"`
	Biller     int `json:"biller"`
}

type PPOBOrderDetail struct {
	Id                string      `json:"id" validate:"required"`
	ItemId            string      `json:"item_id" validate:"required"`
	ItemParentId      string      `json:"item_parent_id"`
	Inquiry           PPOBInquiry `json:"inquiry"`
	ItemName          string      `json:"item_name" validate:"required"`
	Nominal           int         `json:"nominal"`
	TransactionStatus string      `json:"transaction_status"`
	Price             int         `json:"price"`
	Tax               int         `json:"tax"`
	Fee               int         `json:"fee"`
	Reference         string      `json:"reference"`
	Metadata          interface{} `json:"metadata"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
	UpdatedBy         string      `json:"updated_by"`
	ProductCode       string      `json:"product_code"`
}

type PPOBInquiry struct {
	Id         string      `json:"id"`
	AccountRef string      `json:"account_ref"`
	Expiry     string      `json:"expiry"`
	Response   interface{} `json:"response"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  *string     `json:"updated_at"`
	DeletedAt  *string     `json:"deleted_at"`
}

type PPOBOrderHistory struct {
	Id        string      `json:"id" validate:"required"`
	Status    string      `json:"status" validate:"required"`
	Metadata  interface{} `json:"metadata"`
	CreatedAt string      `json:"created_at"`
	CreatedBy string      `json:"created_by"`
}

// ==================== Aggregator Metadata ====================

type AggregatorMetadata struct {
	Aggregator *Aggregator `json:"aggregator,omitempty"`
}

type Aggregator struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}
