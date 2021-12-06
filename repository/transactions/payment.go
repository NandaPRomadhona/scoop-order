package transactions

import (
	"context"
	"database/sql"
	"fmt"
	"scoop-order/internal/schemas"
	"scoop-order/repository"
	"time"
)

//func (transaction *SQLTransaction) PaymentTx(ctx context.Context, checkout schemas.CheckoutTxResult) (interface{}, error) {
//	panic("implement me")
//}

func (transaction *SQLTransaction) PaymentTx(ctx context.Context, checkout schemas.CheckoutTxResult, paymentStatus int32) error {

	err := transaction.execDBTx(ctx, func(q *repository.Queries) error {
		var err error
		// Update Order data (WAITING_FOR_PAYMENT = 20001)
		paymentArg := repository.CreatePaymentParams{
			OrderID:              checkout.OrderID,
			UserID:               checkout.UserID,
			PaymentgatewayID:     checkout.PaymentGatewayID,
			CurrencyCode:         checkout.Currency,
			Amount:               checkout.FinalAmount,
			PaymentStatus:        paymentStatus,
			IsActive:             true,
			IsTestPayment:        false,
			IsTrial:              false,
			PaymentDatetime:      time.Now(),
			FinancialArchiveDate: time.Now(),
			//MerchantParams: json.RawMessage{},
		}
		_, err = transaction.CreatePayment(ctx, paymentArg)
		if err != nil {
			errorMsg := fmt.Errorf("CreatePayment : %s", err.Error())
			return errorMsg
		}

		err = updateOrderData(checkout, transaction)
		if err != nil {
			errorMsg := fmt.Errorf("updateOrderData : %s", err.Error())
			return errorMsg
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (transaction *SQLTransaction) CompletePaymentTx(ctx context.Context, arg schemas.PaymentTxParams) (string, error) {
	var status string
	err := transaction.execDBTx(ctx, func(q *repository.Queries) error {
		//update order
		updateOrderParams := repository.UpdateOrderParams{
			ID:          arg.OrderID,
			Modified:    time.Now(),
			OrderStatus: sql.NullInt32{Int32: arg.OrderStatus, Valid: true},
		}
		order, err := transaction.UpdateOrder(context.Background(), updateOrderParams)
		if err != nil {
			return fmt.Errorf(" failed upadate complete Order #%d: , %s", order.ID, err.Error())
		}

		//update orderlines
		updateOrderLineParams := repository.UpdateOrderlinesParams{
			OrderID:         arg.OrderID,
			Modified:        time.Now(),
			OrderlineStatus: sql.NullInt32{Int32: arg.OrderStatus, Valid: true},
		}
		ol, err := transaction.UpdateOrderlines(context.Background(), updateOrderLineParams)
		if err != nil {
			return fmt.Errorf(" failed upadate complete Order line #%d: , %s", ol.ID, err.Error())
		}

		//update payments
		updatePaymentParams := repository.UpdatePaymentParams{
			OrderID:       arg.OrderID,
			Modified:      time.Now(),
			PaymentStatus: sql.NullInt32{Int32: arg.PaymentStatus, Valid: true},
		}
		payment, err := transaction.UpdatePaymentByOrder(context.Background(), updatePaymentParams)
		if err != nil {
			return fmt.Errorf(" failed upadate complete payment %d: , %s", payment.ID, err.Error())
		}
		status = "complete"

		//produce complete order data in kafka or hit owned user
		return nil
	})
	if err != nil {
		fmt.Println("ERRORS: ", err.Error())
		return "failed", err
	}
	return status, nil
}

func updateOrderData(checkout schemas.CheckoutTxResult, transaction *SQLTransaction) error {
	arg := repository.UpdateOrderParams{
		ID:          checkout.OrderID,
		Modified:    time.Now(),
		OrderStatus: sql.NullInt32{Int32: checkout.OrderStatus, Valid: true},
	}
	_, err := transaction.UpdateOrder(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("upadate Order: , %s", err.Error())
	}
	argOl := repository.UpdateOrderlinesParams{
		OrderID:         checkout.OrderID,
		Modified:        time.Now(),
		OrderlineStatus: sql.NullInt32{Int32: checkout.OrderStatus, Valid: true},
	}
	_, err = transaction.UpdateOrderlines(context.Background(), argOl)
	if err != nil {
		return fmt.Errorf("upadate Order line: , %s", err.Error())
	}
	return nil
}
