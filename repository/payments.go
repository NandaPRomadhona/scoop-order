package repository

import (
	"context"
	"database/sql"
	"scoop-order/internal/databases"
	"time"
)

const selectPaymentGateway = `-- name: selectPaymentGateways :one
 SELECT id, name, is_active, base_currency_id, minimal_amount, is_renewal, payment_group FROM core_paymentgateways WHERE id = $1`

type SelectPaymentGatewaysRow struct {
	PaymentGatewayID int32           `json:"id"`
	Name             string          `json:"name"`
	IsActive         bool            `json:"is_active"`
	BaseCurrencyID   int32           `json:"base_currency_id"`
	MinimalAmount    sql.NullFloat64 `json:"minimal_amount"`
	IsRenewal        bool            `json:"is_renewal"`
	PaymentGroup     string          `json:"payment_group"`
}

func (q *Queries) SelectPaymentGateway(ctx context.Context, paymentGatewayID int32) (SelectPaymentGatewaysRow, error) {
	row := q.db.QueryRowContext(ctx, selectPaymentGateway, paymentGatewayID)
	var i SelectPaymentGatewaysRow
	err := row.Scan(
		&i.PaymentGatewayID,
		&i.Name,
		&i.IsActive,
		&i.BaseCurrencyID,
		&i.MinimalAmount,
		&i.IsRenewal,
		&i.PaymentGroup,
	)
	return i, err
}

const selectPaymentGateways = `-- name: selectPaymentGateways :many
 SELECT id, name, is_active, base_currency_id, minimal_amount, is_renewal, payment_group FROM core_paymentgateways WHERE is_active = true`

func (q *Queries) SelectPaymentGateways(ctx context.Context) ([]SelectPaymentGatewaysRow, error) {
	rows, err := q.db.QueryContext(ctx, selectPaymentGateways)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectPaymentGatewaysRow
	for rows.Next() {
		var i SelectPaymentGatewaysRow
		if err := rows.Scan(
			&i.PaymentGatewayID,
			&i.Name,
			&i.IsActive,
			&i.BaseCurrencyID,
			&i.MinimalAmount,
			&i.IsRenewal,
			&i.PaymentGroup,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createPayment = `-- name: CreatePayment :one
INSERT INTO core_payments(created, modified, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
RETURNING created, modified, id, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial, merchant_params
`

type CreatePaymentParams struct {
	Created              time.Time `json:"created"`
	Modified             time.Time `json:"modified"`
	OrderID              int32     `json:"order_id"`
	UserID               int32     `json:"user_id"`
	PaymentgatewayID     int32     `json:"paymentgateway_id"`
	CurrencyCode         string    `json:"currency_code"`
	Amount               float64   `json:"amount"`
	PaymentStatus        int32     `json:"payment_status"`
	IsActive             bool      `json:"is_active"`
	IsTestPayment        bool      `json:"is_test_payment"`
	PaymentDatetime      time.Time `json:"payment_datetime"`
	FinancialArchiveDate time.Time `json:"financial_archive_date"`
	IsTrial              bool      `json:"is_trial"`
	//MerchantParams       json.RawMessage `json:"merchant_params"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (databases.CorePayment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.Created,
		arg.Modified,
		arg.OrderID,
		arg.UserID,
		arg.PaymentgatewayID,
		arg.CurrencyCode,
		arg.Amount,
		arg.PaymentStatus,
		arg.IsActive,
		arg.IsTestPayment,
		arg.PaymentDatetime,
		arg.FinancialArchiveDate,
		arg.IsTrial,
		//arg.MerchantParams,
	)
	var i databases.CorePayment
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.PaymentgatewayID,
		&i.CurrencyCode,
		&i.Amount,
		&i.PaymentStatus,
		&i.IsActive,
		&i.IsTestPayment,
		&i.PaymentDatetime,
		&i.FinancialArchiveDate,
		&i.IsTrial,
		&i.MerchantParams,
	)
	return i, err
}

const updatePayment = `-- name: UpdatePAYMENT:one
UPDATE core_payments
SET modified            = $2,
    payment_status        = $3
WHERE order_id = $1 
RETURNING created, modified, id, order_id, user_id, paymentgateway_id, currency_code, amount, payment_status, is_active, is_test_payment, payment_datetime, financial_archive_date, is_trial, merchant_params
`

type UpdatePaymentParams struct {
	OrderID       int32         `json:"order_id"`
	Modified      time.Time     `json:"modified"`
	PaymentStatus sql.NullInt32 `json:"payment_status"`
}

func (q *Queries) UpdatePaymentByOrder(ctx context.Context, arg UpdatePaymentParams) (databases.CorePayment, error) {
	row := q.db.QueryRowContext(ctx, updatePayment,
		arg.OrderID,
		arg.Modified,
		arg.PaymentStatus,
	)
	var i databases.CorePayment
	err := row.Scan(
		&i.Created,
		&i.Modified,
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.PaymentgatewayID,
		&i.CurrencyCode,
		&i.Amount,
		&i.PaymentStatus,
		&i.IsActive,
		&i.IsTestPayment,
		&i.PaymentDatetime,
		&i.FinancialArchiveDate,
		&i.IsTrial,
		&i.MerchantParams,
	)
	return i, err
}
