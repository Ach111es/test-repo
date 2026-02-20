package model

// ==================== NON-PPOB Models ====================

type NonPPOBOrder struct {
	Id             string                `json:"id" validate:"required"`
	ReferenceId    string                `json:"reference_id" validate:"required"`
	Tenant         Tenant                `json:"tenant" validate:"required"`
	Store          Store                 `json:"store" validate:"required"`
	Customer       NonPPOBCustomer       `json:"customer" validate:"required"`
	Delivery       Delivery              `json:"delivery" validate:"required"`
	Payment        NonPPOBPayment        `json:"payment" validate:"required"`
	ItemQty        int                   `json:"item_qty" validate:"required"`
	Amount         int                   `json:"amount" validate:"required"`
	VoucherCode    string                `json:"voucher_code"`
	VoucherAmount  int                   `json:"voucher_amount"`
	ServiceFee     int                   `json:"service_fee"`
	PaymentFee     int                   `json:"payment_fee"`
	DeliveryFee    int                   `json:"delivery_fee"`
	Discount       int                   `json:"discount"`
	InsuranceFee   int                   `json:"insurance_fee"`
	TotalBeforeTax int                   `json:"total_before_tax"`
	TotalTax       int                   `json:"total_tax"`
	Total          int                   `json:"total"`
	IsPaid         bool                  `json:"is_paid"`
	PaymentStatus  string                `json:"payment_status"`
	OrderStatus    string                `json:"order_status"`
	ExpiredAt      string                `json:"expired_at"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
	ReferralCode   string                `json:"referral_code"`
	CommissionFee  CommissionFee         `json:"commission_fee"`
	Details        []NonPPOBOrderDetail  `json:"details"`
	History        []NonPPOBOrderHistory `json:"history"`
	Metadata       interface{}           `json:"metadata"`
	MetadataRaw    string                `json:"-"` // Raw JSON string for storage
}

type Tenant struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Store struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type NonPPOBCustomer struct {
	Id      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type Delivery struct {
	Method      string      `json:"method" validate:"required"`
	ReferenceId string      `json:"reference_id"`
	DeliveryUrl *string     `json:"delivery_url"`
	Metadata    interface{} `json:"metadata"`
}

type NonPPOBPayment struct {
	ReferenceId    string      `json:"reference_id" validate:"required"`
	Channel        string      `json:"channel" validate:"required"`
	PaymentCode    string      `json:"payment_code" validate:"required"`
	PaymentGateway string      `json:"payment_gateway" validate:"required"`
	Metadata       interface{} `json:"metadata"`
}

type CommissionFee struct {
	Jatis      int `json:"jatis"`
	Aggregator int `json:"aggregator"`
	Biller     int `json:"biller"`
}

type NonPPOBOrderDetail struct {
	Id        int    `json:"id" validate:"required"`
	Count     int    `json:"count" validate:"required"`
	Deducted  bool   `json:"deducted"`
	Note      string `json:"note"`
	ItemId    string `json:"item_id" validate:"required"`
	ItemName  string `json:"item_name" validate:"required"`
	BasePrice int    `json:"base_price"`
	Price     int    `json:"price"`
	Discount  int    `json:"discount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NonPPOBOrderHistory struct {
	Id          int         `json:"id" validate:"required"`
	OrderId     int         `json:"order_id" validate:"required"`
	OrderStatus string      `json:"order_status" validate:"required"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata"`
	CreatedAt   string      `json:"created_at"`
}
