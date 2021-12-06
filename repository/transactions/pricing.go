package transactions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	src2 "scoop-order/cmd/src"
	"scoop-order/internal/databases"
	"scoop-order/internal/schemas"
	"scoop-order/repository"
	"strconv"
	"time"
)

// Discount Code Type :
// 1 = allowed join with other promo
// 2 = Not allowed join with other promo

// Discount Type:
// 1 = DISCOUNT_OFFER : 'discount offers
// 2 = DISCOUNT_ORDER : 'discount order
// 3 = DISCOUNT_PG_OFFER : 'discount payment gateway offer
// 4 = DISCOUNT_PG_ORDER : 'discount payment gateway order
// 5 = DISCOUNT_CODE : 'discount code'

// Discount Rule:
// 1. BY AMOUNT
// 2. BY PERCENTAGE
// 3. BY TO AMOUNT #flush all price, with discount price

func (transaction *SQLTransaction) PricingTx(ctx context.Context, arg schemas.CheckPricingRequest) (schemas.CheckPricingResponse, error) {
	var response schemas.CheckPricingResponse
	key := "GetPrice_" + string(arg.UserID) + string(arg.OfferID)
	val, err := transaction.clientRedis.Get(key).Result()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	if val != ""{
		err := json.Unmarshal([]byte(val), &response)
		if err != nil {
			return response, err
		}
		return response, nil
	}

	err = transaction.execDBTx(ctx, func(q *repository.Queries) error {
		offers, err := CheckOffers(ctx, q, arg.OfferID, arg.GeoInfo)
		if err != nil {
			return err
		}
		response, err = Pricing(q, ctx, arg, offers)
		if err != nil {
			return err
		}
		return nil
	})

	// save temporary in redis
	redisData, errJson := json.Marshal(response)
	if errJson != nil {
		fmt.Println("errJson : ", errJson)
	}
	errRedis := transaction.clientRedis.Set(key, redisData, time.Minute).Err()

	if errRedis != nil {
		fmt.Println(errRedis)
		fmt.Println("belum tersimpan di redis nih")
	}

	return response, err
}

func CheckOffers(ctx context.Context, q *repository.Queries, offerID []int32, geoInfo schemas.GeoInfo) ([]databases.CoreOffer, error) {
	var offers []databases.CoreOffer
	for i, _ := range offerID {
		offer, err := q.SelectOfferByID(ctx, offerID[i])
		if err != nil {
			return nil, err
		}
		// check offer is available?
		if offer.OfferStatus != 7 {
			return nil, fmt.Errorf("offer is not available")
		}

		// check brand / items of offers
		var restrictCountries []string
		var allowedCountries []string
		if offer.OfferTypeID == 1 { //single offer
			item, err := q.SelectItemOfSingleOffer(ctx, offer.ID)
			if err != nil {
				return nil, fmt.Errorf("check Offer: %s", err)
			}
			if item.ItemDistributionCountryGroupID.Int32 == 0 {
				restrictCountries = []string{}
			} else {
				restrict, _ := q.SelectRestrictCountriesByOffer(ctx, offer.ID)
				if restrict.RestrictType == 1 {
					restrictCountries = restrict.Countries
				} else {
					allowedCountries = restrict.Countries
				}
			}
		} else {
			restrictCountries = []string{}
		}
		// check restricted
		var isRestricted bool
		if len(restrictCountries) == 0 {
			isRestricted = false
			if len(allowedCountries) != 0 {
				find := src2.FindCountry(restrictCountries, geoInfo.CountryCode)
				if find == false {
					isRestricted = true
				}
			}
		} else {
			isRestricted = src2.FindCountry(restrictCountries, geoInfo.CountryCode)
		}
		if isRestricted {
			return nil, fmt.Errorf(offer.Name.String + "cannot be purchase because it is restricted from your current country")
		}

		// check payment able to purchasing free item <<<BELUM>>
		if offer.IsFree.Bool {
			isAbleToGetFree := true
			if isAbleToGetFree != true {
				return nil, fmt.Errorf("Payment Invalid for " + offer.Name.String)
			}
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func checkPlatformOffers(ctx context.Context, q *repository.Queries, offerID []int32, PlatformID int32) ([]databases.CoreOffer, error) {
	var offers []databases.CoreOffer
	for i, _ := range offerID {
		offer, err := q.SelectPlatformOffer(ctx, offerID[i], PlatformID)
		if err != nil {
			return nil, fmt.Errorf("Platforms Offer: %s ", err)
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func getPricing(offers []databases.CoreOffer, currency string) (float64, float64, []float64, []float64, []bool) {
	var baseAmount float64
	var finalAmount float64
	var isDiscount []bool
	var offerPrice []float64
	var offerFinalPrice []float64
	if currency == "IDR" {
		for _, offer := range offers {
			//check the offer has discount offer
			baseAmount = baseAmount + offer.PriceIdr
			offerPrice = append(offerPrice, offer.PriceIdr)
			finalAmount = finalAmount + offer.PriceIdr
			offerFinalPrice = append(offerFinalPrice, offer.PriceIdr)
			//if offer.IsDiscount.Bool {
			//	isDiscount = append(isDiscount, true)
			//	finalAmount = finalAmount + offer.DiscountPriceIdr.Float64
			//	offerFinalPrice = append(offerFinalPrice, offer.DiscountPriceIdr.Float64)
			//} else {
			//	finalAmount = finalAmount + offer.PriceIdr
			//	offerFinalPrice = append(offerFinalPrice, offer.PriceIdr)
			//}
		}
	} else if currency == "USD" {
		for _, offer := range offers {
			//check the offer has discount offer
			baseAmount = baseAmount + offer.PriceUsd
			offerPrice = append(offerPrice, offer.PriceUsd)
			finalAmount = finalAmount + offer.PriceUsd
			offerFinalPrice = append(offerFinalPrice, offer.PriceUsd)
			//if offer.IsDiscount.Bool {
			//	isDiscount = append(isDiscount, true)
			//	finalAmount = finalAmount + offer.DiscountPriceUsd.Float64
			//	offerFinalPrice = append(offerFinalPrice, offer.DiscountPriceUsd.Float64)
			//} else {
			//	finalAmount = finalAmount + offer.PriceUsd
			//	offerFinalPrice = append(offerFinalPrice, offer.PriceUsd)
			//}
		}
	} else {
		for _, offer := range offers {
			//check the offer has discount offer
			baseAmount = baseAmount + offer.PriceUsd
			offerPrice = append(offerPrice, offer.PriceUsd)
			finalAmount = finalAmount + offer.PriceUsd
			offerFinalPrice = append(offerFinalPrice, offer.PriceUsd)
			//if offer.IsDiscount.Bool {
			//	isDiscount = append(isDiscount, true)
			//	finalAmount = finalAmount + offer.DiscountPriceUsd.Float64
			//	offerFinalPrice = append(offerFinalPrice, offer.DiscountPriceUsd.Float64)
			//} else {
			//	finalAmount = finalAmount + offer.PriceUsd
			//	offerFinalPrice = append(offerFinalPrice, offer.PriceUsd)
			//}
		}
	}

	return baseAmount, finalAmount, offerPrice, offerFinalPrice, isDiscount
}

func getFinalPrice(d databases.CoreDiscount, dc databases.CoreDiscountcode, currency string, price float64, isDiscountOffer []bool, q *repository.Queries, used bool) float64 {
	var finalPrice = price
	// discount can't override with other discount
	//if dc.DiscountType.Int32 == 2 && contains(isDiscountOffer, true) {
	//	return finalPrice
	//}
	if time.Now().Before(d.ValidTo) && time.Now().After(d.ValidFrom) && price >= d.MinIdrOrderPrice.Float64 && price <= d.MaxIdrOrderPrice.Float64 {
		if d.DiscountRule.Int32 == 1 {
			if dc.CurrentUses.Int32 < dc.MaxUses.Int32 {
				if currency == "IDR" {
					finalPrice = price - d.DiscountIdr
				} else {
					finalPrice = price - d.DiscountUsd
				}
				if used {
					dc.CurrentUses = sql.NullInt32{Int32: dc.CurrentUses.Int32 + 1, Valid: true}
					q.UpdateDiscountCode(context.Background(), repository.UpdateDiscountCodeParams{
						ID:          dc.ID,
						Modified:    sql.NullTime{Time: time.Now(), Valid: true},
						CurrentUses: dc.CurrentUses,
					})
				}
			}
		} else if d.DiscountRule.Int32 == 2 {
			if dc.CurrentUses.Int32 < dc.MaxUses.Int32 {
				dis := (price * d.DiscountIdr) / 100
				finalPrice = price - dis
				if used {
					dc.CurrentUses = sql.NullInt32{Int32: dc.CurrentUses.Int32 + 1, Valid: true}
					q.UpdateDiscountCode(context.Background(), repository.UpdateDiscountCodeParams{
						ID:          dc.ID,
						Modified:    sql.NullTime{Time: time.Now(), Valid: true},
						CurrentUses: dc.CurrentUses,
					})
				}
			}
		} else if d.DiscountRule.Int32 == 3 {
			if dc.CurrentUses.Int32 < dc.MaxUses.Int32 {
				if currency == "IDR" {
					finalPrice = d.DiscountIdr
				} else {
					finalPrice = d.DiscountUsd
				}
				if used {
					dc.CurrentUses = sql.NullInt32{Int32: dc.CurrentUses.Int32 + 1, Valid: true}
					q.UpdateDiscountCode(context.Background(), repository.UpdateDiscountCodeParams{
						ID:          dc.ID,
						Modified:    sql.NullTime{Time: time.Now(), Valid: true},
						CurrentUses: dc.CurrentUses,
					})
				}
			}

		} else {
			return finalPrice
		}
		return finalPrice
	}
	return finalPrice
}

func checkPaymentGateways(q *repository.Queries, paymentGatewayID int32) (repository.SelectPaymentGatewaysRow, error) {
	paymentGateway, err := q.SelectPaymentGateway(context.Background(), paymentGatewayID)
	if err != nil {
		return paymentGateway, err
	}
	if paymentGateway.IsActive == false {
		return paymentGateway, fmt.Errorf("Payment Gateways is not active ")
	}

	return paymentGateway, err
}

func calculateOfferDiscount(offerDetails []databases.CoreOffer, currencyCode string) []schemas.Offer {
	var offers []schemas.Offer
	for _, offerDetail := range offerDetails {
		var price = 0.0
		var totalBasePrice = 0.0
		if offerDetail.IsDiscount.Bool == true {
			if currencyCode == "IDR" {
				price = offerDetail.PriceIdr
			} else if currencyCode == "USD" {
				price = offerDetail.PriceUsd
			} else {
				price = float64(offerDetail.PricePoint.Int32)
			}

			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				BasePrice:  price,
				FinalPrice: price,
			}
			totalBasePrice = totalBasePrice + price
			offers = append(offers, offer)
		} else {
			if currencyCode == "IDR" {
				price = offerDetail.PriceIdr
			} else if currencyCode == "USD" {
				price = offerDetail.PriceUsd
			} else {
				price = float64(offerDetail.PricePoint.Int32)
			}
			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				BasePrice:  price,
				FinalPrice: price,
			}
			totalBasePrice = totalBasePrice + price
			offers = append(offers, offer)
		}
	}
	return nil
}

func getDiscountOfferPrice(q *repository.Queries, ctx context.Context, offerDetail databases.CoreOffer, currencyCode string) (float64, error) {
	discounts := offerDetail.DiscountID
	finalPrice := offerDetail.PriceIdr
	if currencyCode == "USD" {
		finalPrice = offerDetail.PriceUsd
	}
	var discountIDs []int32
	for _, discount := range discounts {
		d, _ := strconv.Atoi(discount)
		discountIDs = append(discountIDs, int32(d))
	}
	if len(discountIDs) > 0 {
		detailDiscounts, err := q.SelectDiscountByIDs(ctx, discountIDs)
		if err != nil {
			return finalPrice, fmt.Errorf("[ERROR] Get Discount Offer, %s", err.Error())
		}
		d := detailDiscounts[0]

		// check valid discount
		if time.Now().Before(d.ValidTo) && time.Now().After(d.ValidFrom) {
			fmt.Println("check valid discount true")
			if d.DiscountRule.Int32 == 1 {
				if currencyCode == "IDR" {
					finalPrice = offerDetail.PriceIdr - d.DiscountIdr
				} else {
					finalPrice = offerDetail.PriceUsd - d.DiscountUsd
				}
			} else if d.DiscountRule.Int32 == 2 {
				if currencyCode == "IDR" {
					dis := (offerDetail.PriceIdr * d.DiscountIdr) / 100
					finalPrice = offerDetail.PriceIdr - dis
				} else {
					dis := (offerDetail.PriceUsd * d.DiscountIdr) / 100
					finalPrice = offerDetail.PriceUsd - dis
				}
			} else if d.DiscountRule.Int32 == 3 {
				if currencyCode == "IDR" {
					finalPrice = d.DiscountIdr
				} else {
					finalPrice = d.DiscountUsd
				}

			} else {
				return finalPrice, nil
			}
			return finalPrice, nil
		}
	}

	return finalPrice, nil
}

func Pricing(q *repository.Queries, ctx context.Context, request schemas.CheckPricingRequest, offerDetails []databases.CoreOffer) (schemas.CheckPricingResponse, error) {
	var result schemas.CheckPricingResponse
	var dc databases.CoreDiscountcode
	var err error
	var d databases.CoreDiscount
	isDiscountCode := false
	isDCwDO := false
	if request.DiscountCode != "" {
		dc, err = q.SelectDiscountCodeByCode(ctx, request.DiscountCode)
		if err != nil {
			return result, fmt.Errorf("[ERROR] discount code: %s", err.Error())
		}
		d, err = q.SelectDiscountByID(ctx, dc.ID)
		if err != nil {
			return result, fmt.Errorf("[ERROR] discount: %s", err.Error())
		}
		isDiscountCode = true
		if dc.DiscountType.Int32 == 1 {
			isDCwDO = true
		}
	}

	if isDiscountCode && !isDCwDO {
		var offers []schemas.Offer
		totalBasePrice := 0.0
		// discount code not allowed join with other discount (only discount code)
		for _, offerDetail := range offerDetails {
			var price = 0.0
			if request.CurrencyCode == "IDR" {
				price = offerDetail.PriceIdr
			} else if request.CurrencyCode == "USD" {
				price = offerDetail.PriceUsd
			} else {
				price = float64(offerDetail.PricePoint.Int32)
			}
			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				BasePrice:  price,
				FinalPrice: price,
			}
			totalBasePrice = totalBasePrice + price
			offers = append(offers, offer)
		}
		totalFinalPrice := getFinalPrice(d, dc, request.CurrencyCode, totalBasePrice, []bool{false}, q, false)
		result = schemas.CheckPricingResponse{
			CurrencyCode:    request.CurrencyCode,
			Offers:          offers,
			TotalBasePrice:  totalBasePrice,
			TotalDiscount:   totalBasePrice - totalFinalPrice,
			TotalFinalPrice: totalFinalPrice,
		}
		return result, nil
	} else if isDiscountCode && isDCwDO {
		// discount code, allowed other discount
		var offers []schemas.Offer
		totalBasePrice := 0.0
		totalFinalPrice := 0.0
		// discount code not allowed join with other discount (only discount code)
		for _, offerDetail := range offerDetails {
			if err != nil {
				return result, fmt.Errorf("[ERROR] offer: %s", err.Error())
			}
			var basePrice = 0.0
			var finalPrice = 0.0
			if request.CurrencyCode == "IDR" {
				basePrice = offerDetail.PriceIdr
			} else if request.CurrencyCode == "USD" {
				basePrice = offerDetail.PriceUsd
			} else {
				basePrice = float64(offerDetail.PricePoint.Int32)
			}
			// check if offer have discount
			if offerDetail.IsDiscount.Bool {
				finalPrice, err = getDiscountOfferPrice(q, ctx, offerDetail, request.CurrencyCode)
			}else {
				finalPrice = basePrice
			}

			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				BasePrice:  basePrice,
				FinalPrice: finalPrice,
				Discount: basePrice-finalPrice,
			}
			totalBasePrice = totalBasePrice + basePrice
			totalFinalPrice = totalFinalPrice + finalPrice
			offers = append(offers, offer)
		}
		totalFinalPrice = getFinalPrice(d, dc, request.CurrencyCode, totalFinalPrice, []bool{false}, q, false)
		result = schemas.CheckPricingResponse{
			CurrencyCode:    request.CurrencyCode,
			Offers:          offers,
			TotalBasePrice:  totalBasePrice,
			TotalDiscount:   totalBasePrice - totalFinalPrice,
			TotalFinalPrice: totalFinalPrice,
		}
		return result, nil
	} else {
		// check offer price only
		var offers []schemas.Offer
		totalBasePrice := 0.0
		totalFinalPrice := 0.0
		// discount code not allowed join with other discount (only discount code)
		for _, offerDetail := range offerDetails {
			if err != nil {
				return result, fmt.Errorf("[ERROR] offer: %s", err.Error())
			}
			var basePrice = 0.0
			var finalPrice = 0.0
			if request.CurrencyCode == "IDR" {
				basePrice = offerDetail.PriceIdr
			} else if request.CurrencyCode == "USD" {
				basePrice = offerDetail.PriceUsd
			} else {
				basePrice = float64(offerDetail.PricePoint.Int32)
			}
			// check if offer have discount
			if offerDetail.IsDiscount.Bool {
				finalPrice, err = getDiscountOfferPrice(q, ctx, offerDetail, request.CurrencyCode)
			}else {
				finalPrice = basePrice
			}

			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				BasePrice:  basePrice,
				FinalPrice: finalPrice,
				Discount: basePrice-finalPrice,
			}

			totalBasePrice = totalBasePrice + basePrice
			totalFinalPrice = totalFinalPrice + finalPrice
			offers = append(offers, offer)
		}
		result = schemas.CheckPricingResponse{
			CurrencyCode:    request.CurrencyCode,
			Offers:          offers,
			TotalBasePrice:  totalBasePrice,
			TotalDiscount:   totalBasePrice - totalFinalPrice,
			TotalFinalPrice: totalFinalPrice,
		}
		return result, nil

	}
}