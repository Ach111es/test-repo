package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	httpmodel "git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/handler/model"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/utility"
	"github.com/go-playground/validator/v10"
)

type PPOBUsecase interface {
	Create(order *httpmodel.PPOBOrder) error
}

type PPOBUsecaseImpl struct {
	validator     *validator.Validate
	configuration utility.Configuration
}

func extractAggregatorMetadata(metadata interface{}) (string, string, string) {
	if metadata == nil {
		return "", "", ""
	}

	metadataMap, ok := metadata.(map[string]interface{})
	if !ok {
		return "", "", ""
	}

	aggData, exists := metadataMap["aggregator"]
	if !exists {
		return "", "", ""
	}

	agg, ok := aggData.(map[string]interface{})
	if !ok {
		return "", "", ""
	}

	var aggregatorId, aggregatorName, aggregatorPhone string
	if id, ok := agg["id"].(string); ok {
		aggregatorId = id
	}
	if name, ok := agg["name"].(string); ok {
		aggregatorName = name
	}
	if contact, ok := agg["contact"].(string); ok {
		aggregatorPhone = contact
	}

	return aggregatorId, aggregatorName, aggregatorPhone
}

func NewPPOBUsecase(configuration utility.Configuration) PPOBUsecase {
	return &PPOBUsecaseImpl{
		validator:     validator.New(),
		configuration: configuration,
	}
}

func (u *PPOBUsecaseImpl) Create(order *httpmodel.PPOBOrder) error {
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

	// Extract aggregator metadata
	aggregatorId, aggregatorName, aggregatorPhone := extractAggregatorMetadata(order.Metadata)

	// Extract product info from payload
	var productId, productName, productCategory string
	if order.Product.Id != "" {
		productId = order.Product.Id
		productName = order.Product.Name
		productCategory = order.Product.Category
	}

	// Extract category info from payload
	var categoryId, categoryName string
	if order.Category.Id != "" {
		categoryId = order.Category.Id
		categoryName = order.Category.Name
	}

	// Extract account ref from details (first item)
	var accountRef string
	if len(order.Details) > 0 {
		accountRef = order.Details[0].Inquiry.AccountRef
	}

	// Marshal metadata objects to JSON strings
	metadataJson, _ := json.Marshal(order.Metadata)
	paymentMetadataJson, _ := json.Marshal(order.Payment.Metadata)

	strSql = fmt.Sprintf("INSERT INTO orders_ppob (id, reference_id, customer_id, customer_name, customer_type, customer_phone, product_id, product_name, product_category, category_id, category_name, account_ref, payment_reference, payment_channel, payment_code, payment_gateway, payment_metadata, amount, commission_jatis, commission_biller, commission_aggregator, service_fee, payment_fee, total_tax, total_fee, total, payment_status, order_status, source, expired_at, created_at, updated_at, aggregator_id, aggregator_name, aggregator_phone, metadata, details, history, metadata_source) VALUES(%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v);",
		utility.NullIfEmpty(order.Id),
		utility.NullIfEmpty(order.ReferenceId),
		utility.NullIfEmpty(order.Customer.Id),
		utility.NullIfEmpty(utility.SafeString(order.Customer.Name)),
		utility.NullIfEmpty(utility.SafeString(order.Customer.Type)),
		utility.NullIfEmpty(order.Customer.Phone),
		utility.NullIfEmpty(productId),
		utility.NullIfEmpty(utility.SafeString(productName)),
		utility.NullIfEmpty(utility.SafeString(productCategory)),
		utility.NullIfEmpty(categoryId),
		utility.NullIfEmpty(utility.SafeString(categoryName)),
		utility.NullIfEmpty(accountRef),
		utility.NullIfEmpty(order.Payment.ReferenceId),
		utility.NullIfEmpty(utility.SafeString(order.Payment.Channel)),
		utility.NullIfEmpty(utility.SafeString(order.Payment.PaymentCode)),
		utility.NullIfEmpty(utility.SafeString(order.Payment.PaymentGateway)),
		utility.SafeJsonString(string(paymentMetadataJson)),
		order.Amount,
		order.CommissionFee.Jatis,
		order.CommissionFee.Biller,
		order.CommissionFee.Aggregator,
		order.ServiceFee,
		order.PaymentFee,
		order.TotalTax,
		order.TotalFee,
		order.Total,
		utility.NullIfEmpty(order.PaymentStatus),
		utility.NullIfEmpty(order.OrderStatus),
		utility.NullIfEmpty(order.Source),
		utility.NullIfEmpty(order.ExpiredAt),
		utility.NullIfEmpty(order.CreatedAt),
		utility.NullIfEmpty(order.UpdatedAt),
		utility.NullIfEmpty(aggregatorId),
		utility.NullIfEmpty(utility.SafeString(aggregatorName)),
		utility.NullIfEmpty(aggregatorPhone),
		utility.SafeJsonString(string(metadataJson)),
		utility.SafeJsonString(string(detailsJson)),
		utility.SafeJsonString(string(historyJson)),
		utility.SafeJsonString(order.MetadataRaw),
	)

	//utility.PrintConsole(fmt.Sprintf("Generated query PPOB: %v", strSql), "info")

	// Send SQL string to queue for execution
	err = utility.SendMessageToQueue(u.configuration, strSql)
	if err != nil {
		return err
	}

	return nil
}
