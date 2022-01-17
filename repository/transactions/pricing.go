package transactions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	src2 "scoop-order/cmd/src"
	"scoop-order/internal/configs"
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
	key := "GetPrice_" + string(arg.UserID) + string(arg.OfferID) + arg.DiscountCode
	val, err := transaction.clientRedis.Get(key).Result()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	if val != "" {
		err := json.Unmarshal([]byte(val), &response)
		if err != nil {
			return response, err
		}
		return response, nil
	}

	err = transaction.execDBTx(ctx, func(q *repository.Queries) error {
		offers, err := CheckOffers(ctx, q, arg.OfferID, arg.GeoInfo, arg.CurrencyCode)
		if err != nil {
			return err
		}
		response, err = Pricing(q, ctx, arg, offers)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response, err
	}

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

func (transaction *SQLTransaction) PricingTx2(ctx context.Context, arg schemas.CheckPricingRequest) (schemas.CheckPricingResponse, error) {
	var response schemas.CheckPricingResponse
	key := "GetPrice_" + string(arg.UserID) + string(arg.OfferID) + arg.DiscountCode
	val, err := transaction.clientRedis.Get(key).Result()
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
	}
	if val != "" {
		err := json.Unmarshal([]byte(val), &response)
		if err != nil {
			return response, err
		}
		return response, nil
	}

	err = transaction.execDBTx(ctx, func(q *repository.Queries) error {
		offers, err := CheckOffers(ctx, q, arg.OfferID, arg.GeoInfo, arg.CurrencyCode)
		if err != nil {
			return err
		}
		response, err = Pricing2(q, ctx, arg, offers)
		if err != nil {
			return err
		}
		return fmt.Errorf("OK")
	})
	if err.Error() != "OK" {
		return response, err
	}

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

	return response, nil
}

func CheckOffers(ctx context.Context, q *repository.Queries, offerID []int32, geoInfo schemas.GeoInfo, currencyCode string) ([]databases.CoreOffer, error) {
	var offers []databases.CoreOffer
	for i, _ := range offerID {
		offer, err := q.SelectOfferByID(ctx, offerID[i])
		if err != nil {
			return nil, err
		}
		// PTS can't use point
		if currencyCode == "PTS" && offer.OfferTypeID == configs.OfferTypeBuffet {
			return nil, fmt.Errorf("package can't buy using point")
		}

		// check offer is available?
		if offer.OfferStatus != configs.OfferStatusReadyToSale {
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
				//GROUP TYPE:
				//1. RESTRICTS
				//2. ALLOWS
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
				find := src2.FindString(restrictCountries, geoInfo.CountryCode)
				if find == false {
					isRestricted = true
				}
			}
		} else {
			isRestricted = src2.FindString(restrictCountries, geoInfo.CountryCode)
		}
		if isRestricted {
			return nil, fmt.Errorf(" The \"" + offer.Name.String + "\" cannot be purchase because it is restricted from your current country")
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
			baseAmount = baseAmount + float64(offer.PricePoint.Int32)
			offerPrice = append(offerPrice, float64(offer.PricePoint.Int32))
			finalAmount = finalAmount + float64(offer.PricePoint.Int32)
			offerFinalPrice = append(offerFinalPrice, float64(offer.PricePoint.Int32))
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

func checkMinMaxPrice(price float64, currency string, d databases.CoreDiscount) (bool, string) {
	switch currency {
	case "IDR":
		if d.MinIdrOrderPrice.Float64 > 0.0 {
			if price < d.MinIdrOrderPrice.Float64 {
				return false, "Less then minimum order"
			}
		}
		if d.MaxIdrOrderPrice.Float64 > 0.0 {
			if price >= d.MaxIdrOrderPrice.Float64 {
				return false, "More then maximum order"
			}
		}
		return true, "OK"
	case "USD":
		if d.MinUsdOrderPrice.Float64 > 0.0 {
			if price < d.MinUsdOrderPrice.Float64 {
				return false, "Less then minimum total order"
			}
		}
		if d.MaxUsdOrderPrice.Float64 > 0.0 {
			if price >= d.MaxUsdOrderPrice.Float64 {
				return false, "More then maximum total order"
			}
		}
		return true, "OK"
	}
	return true, "OK"
}

func getFinalPrice(d databases.CoreDiscount, dc databases.CoreDiscountcode, currency string, price float64, isDiscountOffer []bool, q *repository.Queries, used bool) (float64, string) {
	var finalPrice = price
	var msg = "discount code is used"
	// discount can't override with other discount
	//if dc.DiscountType.Int32 == 2 && contains(isDiscountOffer, true) {
	//	return finalPrice
	//}
	if time.Now().Before(d.ValidTo) && time.Now().After(d.ValidFrom) {
		ok, msgPrice := checkMinMaxPrice(price, currency, d)
		if ok {
			if d.DiscountRule.Int32 == configs.DiscountRuleAmount {
				if dc.CurrentUses.Int32 < dc.MaxUses.Int32 {
					if currency == "IDR" {
						finalPrice = price - d.DiscountIdr
					} else if currency == "USD" {
						finalPrice = price - d.DiscountUsd
					} else {
						finalPrice = price - float64(d.DiscountPoint.Int32)
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
			} else if d.DiscountRule.Int32 == configs.DiscountRulPercentage {
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
			} else if d.DiscountRule.Int32 == configs.DiscountRulePrice {
				if dc.CurrentUses.Int32 < dc.MaxUses.Int32 {
					if currency == "IDR" {
						finalPrice = d.DiscountIdr
					} else if currency == "USD" {
						finalPrice = d.DiscountUsd
					} else {
						finalPrice = float64(d.DiscountPoint.Int32)
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
				return finalPrice, msg
			}
			return finalPrice, msg
		} else {
			msg = msgPrice
		}
	} else {
		msg = "Discount Code is Not Available Now"
	}
	return finalPrice, msg
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
				Name:       offerDetail.Name.String,
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
				Name:       offerDetail.Name.String,
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
			return calDiscount(d, offerDetail, finalPrice, currencyCode)
		}
	}

	return finalPrice, nil
}

func calDiscount(d databases.CoreDiscount, offerDetail databases.CoreOffer, finalPrice float64, currencyCode string) (float64, error) {
	if d.DiscountRule.Int32 == configs.DiscountRuleAmount {
		if currencyCode == "IDR" {
			finalPrice = finalPrice - d.DiscountIdr
		} else if currencyCode == "USD" {
			finalPrice = finalPrice - d.DiscountUsd
		} else {
			finalPrice = finalPrice - float64(d.DiscountPoint.Int32)
		}
	} else if d.DiscountRule.Int32 == configs.DiscountRulPercentage {
		if currencyCode == "IDR" {
			dis := (finalPrice * d.DiscountIdr) / 100
			finalPrice = finalPrice - dis
		} else if currencyCode == "USD" {
			dis := (finalPrice * d.DiscountUsd) / 100
			finalPrice = finalPrice - dis
		} else {
			dis := (d.DiscountPoint.Int32) / 100
			finalPrice = float64(offerDetail.PricePoint.Int32 - dis)
		}
	} else if d.DiscountRule.Int32 == configs.DiscountRulePrice {
		//if currencyCode == "IDR" {
		//	finalPrice = d.DiscountIdr
		//} else if currencyCode == "USD" {
		//	finalPrice = d.DiscountUsd
		//} else {
		//	finalPrice = float64(d.DiscountPoint.Int32)
		//}
		return finalPrice, nil
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
			log.Println(err)
			return result, fmt.Errorf("[ERROR] discount code not found")
		}
		d, err = q.SelectDiscountByID(ctx, dc.DiscountID)
		if err != nil {
			return result, fmt.Errorf("[ERROR] discount is not available")
		}
		isDiscountCode = true
		if dc.DiscountType.Int32 == configs.DiscountOffer {
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
				Name:       offerDetail.Name.String,
				BasePrice:  price,
				FinalPrice: price,
			}
			totalBasePrice = totalBasePrice + price
			offers = append(offers, offer)
		}
		totalFinalPrice, _ := getFinalPrice(d, dc, request.CurrencyCode, totalBasePrice, []bool{false}, q, false)
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
			} else {
				finalPrice = basePrice
			}

			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				Name:       offerDetail.Name.String,
				BasePrice:  basePrice,
				FinalPrice: finalPrice,
				Discount:   basePrice - finalPrice,
			}
			totalBasePrice = totalBasePrice + basePrice
			totalFinalPrice = totalFinalPrice + finalPrice
			offers = append(offers, offer)
		}
		totalFinalPrice, _ = getFinalPrice(d, dc, request.CurrencyCode, totalFinalPrice, []bool{false}, q, false)
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
			} else {
				finalPrice = basePrice
			}

			offer := schemas.Offer{
				OfferID:    offerDetail.ID,
				Name:       offerDetail.Name.String,
				BasePrice:  basePrice,
				FinalPrice: finalPrice,
				Discount:   basePrice - finalPrice,
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

func constructFinalPrice(offerDetails []databases.CoreOffer, currency string) schemas.CheckPricingResponse {
	var dataOffers []schemas.Offer
	var totalBasePrice = 0.0
	var totalDiscount = 0.0
	var totalFinalPrice = 0.0
	for _, offer := range offerDetails {
		var dataOffer schemas.Offer
		switch currency {
		case "IDR":
			dataOffer = schemas.Offer{
				OfferID:    offer.ID,
				Name:       offer.Name.String,
				BasePrice:  offer.PriceIdr,
				FinalPrice: offer.PriceIdr,
				Discount:   0.0,
			}
		case "USD":
			dataOffer = schemas.Offer{
				OfferID:    offer.ID,
				Name:       offer.Name.String,
				BasePrice:  offer.PriceUsd,
				FinalPrice: offer.PriceUsd,
				Discount:   0.0,
			}
		case "PTS":
			dataOffer = schemas.Offer{
				OfferID:    offer.ID,
				Name:       offer.Name.String,
				BasePrice:  float64(offer.PricePoint.Int32),
				FinalPrice: float64(offer.PricePoint.Int32),
				Discount:   0.0,
			}
		}
		dataOffers = append(dataOffers, dataOffer)
		totalBasePrice += dataOffer.BasePrice
		totalFinalPrice += dataOffer.FinalPrice
		totalDiscount += dataOffer.Discount
	}
	return schemas.CheckPricingResponse{
		CurrencyCode:    currency,
		Offers:          dataOffers,
		TotalBasePrice:  totalBasePrice,
		TotalFinalPrice: totalFinalPrice,
		TotalDiscount:   totalDiscount,
	}
}

func Pricing2(q *repository.Queries, ctx context.Context, request schemas.CheckPricingRequest, offerDetails []databases.CoreOffer) (schemas.CheckPricingResponse, error) {
	//  construct the response
	var pricingResponse schemas.CheckPricingResponse
	var dc databases.CoreDiscountcode
	var err error
	var discountOfDC databases.CoreDiscount
	var isDiscountCode = false
	var isDCCanUsed = false
	var msg string

	pricingResponse = constructFinalPrice(offerDetails, request.CurrencyCode)

	/*check discount code can be used*/
	if request.DiscountCode != "" {
		dc, err = q.SelectDiscountCodeByCode(ctx, request.DiscountCode)
		if err != nil {
			log.Println(err)
			return pricingResponse, fmt.Errorf("[ERROR] discount code not found")
		}
		discountOfDC, err = q.SelectDiscountByID(ctx, dc.DiscountID)
		if err != nil {
			return pricingResponse, fmt.Errorf("[ERROR] discount is not available")
		}
		_, err = q.SelectAllowedPG(ctx, discountOfDC.ID, request.PaymentGatewayID)
		if err != nil {
			return pricingResponse, fmt.Errorf("[ERROR] discount is not available for preferred payment gateway")
		}
		_, err = q.SelectAllowedPlatform(ctx, discountOfDC.ID, request.PlatformID)
		if err != nil {
			return pricingResponse, fmt.Errorf("[ERROR] discount is not available for preferred platform")
		}
		isDiscountCode = true
	}

	// 1. Check due date
	if time.Now().Before(discountOfDC.ValidTo) && time.Now().After(discountOfDC.ValidFrom) {
		isDCCanUsed = true
	} else {
		return pricingResponse, fmt.Errorf("[ERROR] discount is not valid")
	}

	// 2. Check Max User
	if dc.MaxUses.Int32 > 0 && dc.CurrentUses.Int32 >= dc.MaxUses.Int32 {
		return pricingResponse, fmt.Errorf("[ERROR] discount code is not valid (out of max uses number)")
	}

	// 3. Check Max Min
	isDCCanUsed, msg = checkMinMaxPrice(pricingResponse.TotalBasePrice, request.CurrencyCode, discountOfDC)
	if isDCCanUsed == false {
		return pricingResponse, fmt.Errorf("[ERROR] discount code is not valid (%s)", msg)
	}

	//	check offer discount
	var isOfferHaveDiscount = make(map[int]bool)
	for idxOffer, offer := range offerDetails {
		discounGoCut := false
		var validDiscountOffer []databases.CoreDiscount
		if offer.IsDiscount.Bool && !request.IsRenewal {
			var discountIDs []int32
			for _, discount := range offer.DiscountID {
				d, _ := strconv.Atoi(discount)
				discountIDs = append(discountIDs, int32(d))
			}
			var discountData []databases.CoreDiscount
			if len(discountIDs) > 0 {
				discountData, err = q.SelectDiscountByIDs(ctx, discountIDs)
			}
			for _, discount := range discountData {
				if time.Now().Before(discount.ValidTo) && time.Now().After(discount.ValidFrom) {
					discounGoCut = true
					validDiscountOffer = append(validDiscountOffer, discount)
				}
			}
		}

		// cut offer price by discount
		if discounGoCut {
			isOfferHaveDiscount[idxOffer] = true
			for _, discount := range validDiscountOffer {
				// Check Trial discount
				/* ... */
				// temp discount
				if offer.IsFree.Bool {
					pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].FinalPrice
					pricingResponse.Offers[idxOffer].FinalPrice = 0.0
					break
				} else {
					finalOfferPrice, _ := calDiscount(discount, offer, pricingResponse.Offers[idxOffer].FinalPrice, pricingResponse.CurrencyCode)
					pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].Discount + (pricingResponse.Offers[idxOffer].FinalPrice - finalOfferPrice)
					pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
				}
			}
		} else {
			// CHECK DISCOUNT DATA =============================== [ PREDEFINED GROUP, SUCKS ]
			var items []databases.CoreItem
			isSingle := true
			if offer.OfferTypeID == configs.OfferTypeSingle && offer.OfferTypeID == configs.OfferTypeBundle {
				items, err = q.SelectItemByOfferID(ctx, offer.ID)
				if err != nil {
					return pricingResponse, fmt.Errorf("[ERROR] has any problem in select item (%s)", err.Error())
				}
			} else if offer.OfferTypeID == configs.OfferTypeSubscription {
				isSingle = false
				items, err = q.SelectItemBrandByOfferID(ctx, offer.ID)
				if err != nil {
					return pricingResponse, fmt.Errorf("[ERROR] has any problem in select item (%s)", err.Error())
				}
			}
			var group_predefined []int32
			var discounts []databases.CoreDiscount
			if len(items) > 0 {
				switch items[0].ItemType {
				case configs.ItemTypeMagazine:
					if isSingle {
						group_predefined = []int32{4, 1, 5, 7}
					} else {
						group_predefined = []int32{4, 1, 6, 8}
					}
				case configs.ItemTypeBook:
					//	# HARPER COLLINS 5% := 467
					if items[0].BrandId == 467 {
						group_predefined = []int32{467}
					} else {
						group_predefined = []int32{4, 2, 5}
					}
				case configs.ItemTypeNewspaper:
					group_predefined = []int32{4, 3, 6}
				}
			}
			discounts, err = q.SelectDiscountByPredefinedGroups(ctx, group_predefined, request.PlatformID, request.PaymentGatewayID, time.Now().String())
			if err != nil{
				discounts = []databases.CoreDiscount{}
			}
			if len(discounts) > 1 {
				isOfferHaveDiscount[idxOffer] = true
				if offer.IsFree.Bool {
					pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].FinalPrice
					pricingResponse.Offers[idxOffer].FinalPrice = 0.0
				} else {
					finalOfferPrice, _ := calDiscount(discounts[0], offer, pricingResponse.Offers[idxOffer].FinalPrice, pricingResponse.CurrencyCode)
					pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].Discount + (pricingResponse.Offers[idxOffer].FinalPrice - finalOfferPrice)
					pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
				}
			} else {
				isOfferHaveDiscount[idxOffer] = false
			}
		}
	}

	//	check discount code
	if isDiscountCode {
		allowedOffer, _ := q.SelectAllowedOffers(ctx, discountOfDC.ID)
		allowed := 0
		for idxOffer, offer := range offerDetails {
			// All Offer allowed
			if discountOfDC.PredefinedGroup.Int32 == 4 { //all offer
				isDCCanUsed = true
			} else {
				// Offer In List
				find := src2.FindInt(offer.ID, allowedOffer)
				if find {
					isDCCanUsed = true
				} else {
					isDCCanUsed = false
				}
			}

			if isDCCanUsed && allowed < 1 {
				if discountOfDC.DiscountRule.Int32 == configs.DiscountRuleAmount {
					allowed += 1
				}
				if isOfferHaveDiscount[idxOffer]{
					/*	check if coupon code can override the others promo?
					 if cant, calculate from base prices
					 if can, (stack promo calculations )*/
					if dc.DiscountType.Int32 == configs.DiscountAllowedJoin{
						finalOfferPrice, _ := calDiscount(discountOfDC, offer, pricingResponse.Offers[idxOffer].FinalPrice, pricingResponse.CurrencyCode)
						pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].Discount + (pricingResponse.Offers[idxOffer].FinalPrice - finalOfferPrice)
						pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
					} else {
						finalOfferPrice, _ := calDiscount(discountOfDC, offer, pricingResponse.Offers[idxOffer].BasePrice, pricingResponse.CurrencyCode)
						pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].BasePrice - finalOfferPrice
						pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
					}
				} else {
					finalOfferPrice, _ := calDiscount(discountOfDC, offer, pricingResponse.Offers[idxOffer].FinalPrice, pricingResponse.CurrencyCode)
					pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].Discount + (pricingResponse.Offers[idxOffer].FinalPrice - finalOfferPrice)
					pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
				}
			}
		}
	}

	// check discount paymentgateway
	discountPG, err := q.SelectDiscountsPaymentGateways(ctx, request.PaymentGatewayID, time.Now().String())

	if err == nil {
		for idxOffer, offer := range offerDetails {
			finalOfferPrice, _ := calDiscount(discountPG, offer, pricingResponse.Offers[idxOffer].BasePrice, pricingResponse.CurrencyCode)
			pricingResponse.Offers[idxOffer].Discount = pricingResponse.Offers[idxOffer].BasePrice - finalOfferPrice
			pricingResponse.Offers[idxOffer].FinalPrice = finalOfferPrice
		}
	}

	//	final price
	var totalBasePrice float64
	var totalDiscount float64
	var totalFinalPrice float64
	for _, offer := range pricingResponse.Offers {
		totalBasePrice = totalBasePrice + offer.BasePrice
		totalFinalPrice = totalFinalPrice + offer.FinalPrice
		totalDiscount = totalDiscount + offer.Discount
	}
	pricingResponse.TotalBasePrice = totalBasePrice
	pricingResponse.TotalFinalPrice = totalFinalPrice
	pricingResponse.TotalDiscount = totalDiscount

	return pricingResponse, nil
}
