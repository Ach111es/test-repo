package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	httpmodel "git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/handler/model"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/utility"
	"github.com/go-playground/validator/v10"
)

type NonPPOBUsecase interface {
	Create(order *httpmodel.NonPPOBOrder) error
}

type NonPPOBUsecaseImpl struct {
	validator     *validator.Validate
	configuration utility.Configuration
}

func NewNonPPOBUsecase(configuration utility.Configuration) NonPPOBUsecase {
	return &NonPPOBUsecaseImpl{
		validator:     validator.New(),
		configuration: configuration,
	}
}

func (u *NonPPOBUsecaseImpl) Create(order *httpmodel.NonPPOBOrder) error {
	var strSql string

	err := u.validator.Struct(order)
	if err != nil {
		msg := utility.ValidationErrorHandle(err)
		return errors.New(msg)
	}

	// Marshal nested objects to JSON strings
	detailsJson, err := json.Marshal(order.Details)
	if err != nil {
		return err
	}

	historyJson, err := json.Marshal(order.History)
	if err != nil {
		return err
	}

	// Extract delivery URL (handle pointer)
	var deliveryUrl string
	if order.Delivery.DeliveryUrl != nil {
		deliveryUrl = *order.Delivery.DeliveryUrl
	}

	strSql = fmt.Sprintf("INSERT INTO orders_nonppob (id, reference_id, customer_id, customer_name, customer_phone, customer_address, tenant_id, tenant_name, store_id, store_name, delivery_method, delivery_reference, delivery_fee, delivery_url, payment_reference, payment_channel, payment_code, payment_gateway, payment_fee, item_qty, amount, voucher_code, voucher_amount, discount, service_fee, insurance_fee, total_before_tax, total_tax, total, commission_jatis, commission_aggregator, commission_biller, referral_code, is_paid, payment_status, order_status, expired_at, created_at, updated_at, details, history, metadata_source) VALUES(%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v);",
		utility.NullIfEmpty(order.Id),
		utility.NullIfEmpty(order.ReferenceId),
		utility.NullIfEmpty(order.Customer.Id),
		utility.NullIfEmpty(utility.SafeString(order.Customer.Name)),
		utility.NullIfEmpty(order.Customer.Phone),
		utility.NullIfEmpty(order.Customer.Address),
		utility.NullIfEmpty(order.Tenant.Id),
		utility.NullIfEmpty(utility.SafeString(order.Tenant.Name)),
		utility.NullIfEmpty(order.Store.Id),
		utility.NullIfEmpty(utility.SafeString(order.Store.Name)),
		utility.NullIfEmpty(utility.SafeString(order.Delivery.Method)),
		utility.NullIfEmpty(order.Delivery.ReferenceId),
		order.DeliveryFee,
		utility.NullIfEmpty(utility.SafeString(deliveryUrl)),
		utility.NullIfEmpty(order.Payment.ReferenceId),
		utility.NullIfEmpty(utility.SafeString(order.Payment.Channel)),
		utility.NullIfEmpty(utility.SafeString(order.Payment.PaymentCode)),
		utility.NullIfEmpty(utility.SafeString(order.Payment.PaymentGateway)),
		order.PaymentFee,
		order.ItemQty,
		order.Amount,
		utility.NullIfEmpty(order.VoucherCode),
		order.VoucherAmount,
		order.Discount,
		order.ServiceFee,
		order.InsuranceFee,
		order.TotalBeforeTax,
		order.TotalTax,
		order.Total,
		order.CommissionFee.Jatis,
		order.CommissionFee.Aggregator,
		order.CommissionFee.Biller,
		utility.NullIfEmpty(order.ReferralCode),
		order.IsPaid,
		utility.NullIfEmpty(order.PaymentStatus),
		utility.NullIfEmpty(order.OrderStatus),
		utility.NullIfEmpty(order.ExpiredAt),
		utility.NullIfEmpty(order.CreatedAt),
		utility.NullIfEmpty(order.UpdatedAt),
		utility.SafeJsonString(string(detailsJson)),
		utility.SafeJsonString(string(historyJson)),
		utility.SafeJsonString(order.MetadataRaw),
	)

	//utility.PrintConsole(fmt.Sprintf("Generated query NonPPOB: %v", strSql), "info")

	// Send SQL string to queue for execution
	err = utility.SendMessageToQueue(u.configuration, strSql)
	if err != nil {
		return err
	}

	return nil
}
