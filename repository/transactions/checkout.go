package transactions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"scoop-order/cmd/src"
	"scoop-order/internal/configs"
	"scoop-order/internal/schemas"
	"scoop-order/repository"
	"time"
)

// CheckoutTx is process creating the order, store in DB and count the total amount
/* ORDER STATUS

 */
func (transaction *SQLTransaction) CheckoutTx(ctx context.Context, arg schemas.CheckoutTxParams) (schemas.CheckoutTxResult, error) {
	// empty result
	var result schemas.CheckoutTxResult
	var currency string
	logger := log.Default()
	logger.SetPrefix("DEBUG - ")

	err := transaction.execDBTx(ctx, func(q *repository.Queries) error {
		// check offer, is it available?
		offers, err := CheckOffers(ctx, q, arg.OfferID, arg.GeoInfo)
		if err != nil {
			return fmt.Errorf("checkOffers: %s", err.Error())
		}

		// check platform not web
		if arg.PlatformID != 4 {
			offers, err = checkPlatformOffers(ctx, q, arg.OfferID, arg.PlatformID)
			if err != nil {
				return fmt.Errorf("checkPlatformOffers: %s", err.Error())
			}
		}

		// check payment gateways
		paymentGatewayID := arg.PaymentGatewayID
		paymentGateway, err := checkPaymentGateways(q, paymentGatewayID)
		if err != nil {
			return fmt.Errorf("checkPaymentGateways: %s", err.Error())
		}

		// mapping currency
		currencyID := paymentGateway.BaseCurrencyID
		if currencyID == 1 {
			currency = "USD"
		} else if currencyID == 2 {
			currency = "IDR"
		} else {
			currency = "PTS"
		}

		// Cara Pertama
		var pricing schemas.CheckPricingResponse
		key := "GetPrice_" + string(arg.UserID) + string(arg.OfferID)
		val, err := transaction.clientRedis.Get(key).Result()
		if err != nil {
			logger.Println("Checkout Trx - ", err.Error())
		}
		if val != "" {
			json.Unmarshal([]byte(val), &pricing)
		} else {
			// Cara Kedua
			getPriceArg := schemas.CheckPricingRequest{
				UserID:           arg.UserID,
				OfferID:          arg.OfferID,
				PaymentGatewayID: paymentGatewayID,
				PlatformID:       arg.PlatformID,
				DiscountCode:     arg.DiscountCode,
				CurrencyCode:     currency,
				GeoInfo:          arg.GeoInfo,
			}
			pricing, _ = Pricing(q, ctx, getPriceArg, offers)
		}

		IsDiscount := false
		if arg.DiscountCode != "" {
			key := "DC_" + arg.DiscountCode
			val, _ := transaction.clientRedis.Get(key).Result()
			if val == "" {
				IsDiscount = true
				err = transaction.clientRedis.Set(key, true, time.Minute*60).Err()
				dc, err := q.SelectDiscountCodeByCode(ctx, arg.DiscountCode)
				if err != nil {
					return err
				}
				dc.CurrentUses = sql.NullInt32{Int32: dc.CurrentUses.Int32 + 1, Valid: true}
				q.UpdateDiscountCode(context.Background(), repository.UpdateDiscountCodeParams{
					ID:          dc.ID,
					Modified:    sql.NullTime{Time: time.Now(), Valid: true},
					CurrentUses: dc.CurrentUses,
				})
			}
		}

		// Create Order Data
		order, err := q.CreateOrder(ctx, repository.CreateOrderParams{
			Created:          time.Now(),
			Modified:         time.Now(),
			OrderNumber:      src.GenerateInvoiceNumber(arg.UserID),
			UserID:           arg.UserID,
			OrderStatus:      sql.NullInt32{Int32: configs.OrderNew, Valid: true},
			IsActive:         sql.NullBool{Bool: true, Valid: true},
			TotalAmount:      pricing.TotalBasePrice,
			FinalAmount:      pricing.TotalFinalPrice,
			CurrencyCode:     currency,
			PaymentgatewayID: paymentGatewayID,
		})
		if err != nil {
			return fmt.Errorf(" orders: %s", err.Error())
		}

		// Create Order Line Data
		var orderlines []repository.OrderLinePublished
		for o, offer := range pricing.Offers {
			var discount = false
			if offer.Discount > 0.0 {
				discount = true
			}
			ol, err := q.CreateOrderlines(ctx, repository.CreateOrderlinesParams{
				Created:         time.Now(),
				Modified:        time.Now(),
				OrderID:         sql.NullInt32{Int32: order.ID, Valid: true},
				OfferID:         offer.OfferID,
				IsDiscount:      sql.NullBool{Bool: discount, Valid: true},
				IsActive:        sql.NullBool{Bool: true, Valid: true},
				CurrencyCode:    sql.NullString{String: currency, Valid: true},
				Price:           sql.NullFloat64{Float64: offer.BasePrice, Valid: true},
				FinalPrice:      sql.NullFloat64{Float64: offer.FinalPrice, Valid: true},
				Name:            offers[o].Name,
				UserID:          sql.NullInt32{Int32: arg.UserID, Valid: true},
				OrderlineStatus: sql.NullInt32{Int32: configs.OrderNew, Valid: true},
				IsTrial:         sql.NullBool{Bool: false, Valid: true},
			})
			if err != nil {
				return fmt.Errorf(" orders line: %s", err.Error())
			}
			orderlines = append(orderlines, repository.OrderLinePublished{
				ID:         ol.ID,
				IsDiscount: ol.IsDiscount.Bool,
				OfferID:    ol.OfferID,
				TotalPrice: ol.Price.Float64,
				FinalPrice: ol.FinalPrice.Float64,
			})
		}

		// Create Order Detail
		_, err = q.CreateOrderDetail(ctx, repository.CreateOrderDetailParams{
			Created:     time.Now(),
			Modified:    time.Now(),
			OrderID:     order.ID,
			UserID:      sql.NullInt32{Int32: arg.UserID, Valid: true},
			UserCountry: sql.NullString{String: arg.GeoInfo.CountryCode, Valid: true},
		})

		if err != nil {
			return err
		}

		if order.TotalAmount != order.FinalAmount {
			IsDiscount = true
		}

		result.UserID = arg.UserID
		result.PaymentGatewayID = paymentGatewayID
		result.Currency = currency
		result.TotalAmount = pricing.TotalBasePrice
		result.FinalAmount = pricing.TotalFinalPrice
		result.OrderID = order.ID
		result.OrderNumber = order.OrderNumber
		result.OrderStatus = order.OrderStatus.Int32
		result.Orderline = orderlines
		result.IsDiscount = IsDiscount
		result.IsActive = true
		return nil
	})

	return result, err
}
