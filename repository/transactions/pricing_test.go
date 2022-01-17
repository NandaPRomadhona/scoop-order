package transactions

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"scoop-order/internal/databases"
	"scoop-order/internal/schemas"
	"testing"
)

func TestCheckOffer(t *testing.T) {
	offerID := []int32{1}
	geoInfo := schemas.GeoInfo{CountryCode:"ID"}
	offers, _ := CheckOffers(context.Background(), testQueries, offerID, geoInfo)
	require.NotEmpty(t, offers)
	require.Equal(t, 1, len(offers))

	offerID = []int32{1, 2}
	offers, _ = CheckOffers(context.Background(), testQueries, offerID, geoInfo)
	require.Equal(t, 2, len(offers))
	require.Equal(t, "Full Premium", offers[1].Name.String)
}

func TestGetPrice(t *testing.T) {
	offerID := []int32{1}
	geoInfo := schemas.GeoInfo{CountryCode:"ID"}
	offers, _ := CheckOffers(context.Background(), testQueries, offerID, geoInfo)
	fmt.Println("Offer: ", offers)
	basePrice, FinalPrice, offerPrice, offerFinalPrice, _ := getPricing(offers, "IDR")

	require.Equal(t, 100000.00, basePrice)
	require.Equal(t, 100000.00, FinalPrice)
	require.Equal(t, 1, len(offerPrice))
	require.Equal(t, 1, len(offerFinalPrice))

	//offerID = []int32{3}
	//offers, _ = checkOffers(context.Background(), testQueries, offerID, geoInfo)
	//basePrice, FinalPrice, offerPrice, offerFinalPrice, _ = getPricing(offers, "IDR")
	//
	//require.Equal(t, 719000.00, basePrice)
	//require.Equal(t, 359500.00, FinalPrice)
	//require.Equal(t, 1, len(offerPrice))
	//require.Equal(t, 1, len(offerFinalPrice))
	//
	//offerID = []int32{1, 3}
	//offers, _ = checkOffers(context.Background(), testQueries, offerID, geoInfo)
	//basePrice, FinalPrice, offerPrice, offerFinalPrice, _ = getPricing(offers, "IDR")
	//
	//require.Equal(t, 15000.00+719000.00, basePrice)
	//require.Equal(t, 15000.00+359500.00, FinalPrice)
	//require.Equal(t, 2, len(offerPrice))
	//require.Equal(t, 2, len(offerFinalPrice))
}

func TestCheckPlatformOffers(t *testing.T) {
	offerID := []int32{1}
	platformID := int32(1)
	offer, err := checkPlatformOffers(context.Background(), testQueries, offerID, platformID)
	require.NoError(t, err)
	require.Equal(t, 1, len(offer))
	require.Equal(t, 0.99, offer[0].PriceUsd)

	offerID = []int32{1, 2}
	platformID = int32(1)
	offer, err = checkPlatformOffers(context.Background(), testQueries, offerID, platformID)
	require.NoError(t, err)
	require.Equal(t, 2, len(offer))
	require.Equal(t, 0.99, offer[0].PriceUsd)
	require.Equal(t, 89000.00, offer[1].PriceIdr)

	offerID = []int32{1, 2}
	platformID = int32(2)
	offer, err = checkPlatformOffers(context.Background(), testQueries, offerID, platformID)
	require.Error(t, fmt.Errorf("Platforms Offer: sql: no rows in result set "), err)
}

func TestGetFinalPrice(t *testing.T) {
	testCases := []struct {
		name            string
		discountID      int32
		discountCode    string
		price           float64
		isDiscountOffer []bool
		currencyCode    string
		used            bool
		expect          float64
	}{
		{"IDR, Discount 12%", 1, "DISC12%", 100000.00, []bool{false}, "IDR", false, 88000.00},
		{"IDR, Discount 12K", 2, "DISC12K", 100000.00, []bool{false}, "IDR", false, 88000.00},
		{"IDR, Discount 12%", 3, "ALL15K", 100000.00, []bool{false}, "IDR", false, 15000.00},
	}
	for _, testCase := range testCases {
		discount, _ := testQueries.SelectDiscountByID(context.Background(), testCase.discountID)
		discountCode, _ := testQueries.SelectDiscountCodeByCode(context.Background(), testCase.discountCode)
		price := testCase.price
		isDiscountOffer := testCase.isDiscountOffer
		finalPrice := getFinalPrice(discount, discountCode, testCase.currencyCode, price, isDiscountOffer, testQueries, testCase.used)
		if finalPrice != testCase.expect {
			t.Errorf(testCase.name)
			return
		}
	}

}

func TestGetFinalPrice2(t *testing.T) {
	var d databases.CoreDiscount
	var dc databases.CoreDiscountcode
	var currency string
	var price float64

	d, err = testQueries.SelectDiscountByID(context.Background(), 1)
	dc, err = testQueries.SelectDiscountCodeByCode(context.Background(), "DISC12%")
	currency = "IDR"
	price = 100000.00
	//isDiscountOffer :=[]bool{true}
	//finalPrice := getFinalPrice(d, dc, currency, price, isDiscountOffer, testQueries, true)
	//require.Equal(t, 100000.00, finalPrice)

	isDiscountOffer := []bool{false}
	finalPrice := getFinalPrice(d, dc, currency, price, isDiscountOffer, testQueries, true)
	require.Equal(t, 88000.00, finalPrice)

	dc.CurrentUses = dc.MaxUses
	finalPrice = getFinalPrice(d, dc, currency, price, isDiscountOffer, testQueries, true)
	require.Equal(t, 100000.00, finalPrice)

	// Test Discount Rule 2
	d, err = testQueries.SelectDiscountByID(context.Background(), 2)
	dc, err = testQueries.SelectDiscountCodeByCode(context.Background(), "DISC12K")
	currency = "IDR"
	price = 112000.00
	//isDiscountOffer =[]bool{true}
	//finalPrice = getFinalPrice(d, dc, currency, price, isDiscountOffer,testQueries, true)
	//require.Equal(t, 112000.00, finalPrice)

	isDiscountOffer = []bool{false}
	finalPrice = getFinalPrice(d, dc, currency, price, isDiscountOffer, testQueries, true)
	require.Equal(t, 100000.00, finalPrice)

	dc.CurrentUses = dc.MaxUses
	finalPrice = getFinalPrice(d, dc, currency, price, isDiscountOffer, testQueries, true)
	require.Equal(t, 112000.00, finalPrice)
}

func TestGetDiscountOfferPrice(t *testing.T) {
	testCases := []struct {
		name         string
		offerID      int32
		currencyCode string
		expect       float64
	}{
		{"Offer without discount", 4, "IDR", 15000.00},
		{"Offer with discount rule 1", 5, "IDR", 35200.00},
		{"Offer with discount rule 2", 3, "IDR", 15000.00},
	}

	for _, testCase := range testCases {
		offerDetail, _ := testQueries.SelectOfferByID(context.Background(), testCase.offerID)
		price, err := getDiscountOfferPrice(testQueries, context.Background(), offerDetail, testCase.currencyCode)
		if price != testCase.expect {
			t.Errorf(testCase.name, ": ", err)
			require.Equal(t, testCase.expect, price)
			return
		}
	}
}

func TestPricing(t *testing.T) {
	testCases := []struct {
		name        string
		request     schemas.CheckPricingRequest
		offerIDs    []int32
		expectFinal float64
	}{
		{name: "No Discount, No Offer Discount", request: schemas.CheckPricingRequest{OfferID: []int32{2}, PaymentGatewayID: 1, PlatformID: 1, CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{2}, expectFinal: 89000.00},
		{name: "No Discount Code, Offer Discount Rule 1", request: schemas.CheckPricingRequest{OfferID: []int32{5}, PaymentGatewayID: 1, PlatformID: 1, CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{5}, expectFinal: 35200.00},
		{name: "No Discount Code, Offer Discount Rule 2", request: schemas.CheckPricingRequest{OfferID: []int32{6}, PaymentGatewayID: 1, PlatformID: 1, CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{6}, expectFinal: 28000.00},
		{name: "No Discount Code, Offer Discount Rule 3", request: schemas.CheckPricingRequest{OfferID: []int32{3}, PaymentGatewayID: 1, PlatformID: 1, CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{3}, expectFinal: 15000.00},
		{name: "Discount Code Not Allowed other discount", request: schemas.CheckPricingRequest{OfferID: []int32{1}, PaymentGatewayID: 1, PlatformID: 1, DiscountCode: "DISC12%", CurrencyCode: "", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{1}, expectFinal: 88000.00},
		{name: "Discount Code Allowed other discount, No Offer Discount", request: schemas.CheckPricingRequest{OfferID: []int32{4}, PaymentGatewayID: 1, PlatformID: 1, DiscountCode: "ALLOWALL5K", CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{4}, expectFinal: 10000.00},
		{name: "Discount Code Allowed other discount, Offer Discount Rule 1", request: schemas.CheckPricingRequest{OfferID: []int32{1}, PaymentGatewayID: 1, PlatformID: 1, DiscountCode: "ALLOWALL5K", CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{1}, expectFinal: 83000.00},
		{name: "Discount Code Allowed other discount, Offer Discount Rule 2", request: schemas.CheckPricingRequest{OfferID: []int32{6}, PaymentGatewayID: 1, PlatformID: 1, DiscountCode: "ALLOWALL5K", CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{6}, expectFinal: 23000.00},
		{name: "Discount Code Allowed other discount, Offer Discount Rule 3", request: schemas.CheckPricingRequest{OfferID: []int32{3}, PaymentGatewayID: 1, PlatformID: 1, DiscountCode: "ALLOWALL5K", CurrencyCode: "IDR", GeoInfo: schemas.GeoInfo{CountryCode: "in"}}, offerIDs: []int32{3}, expectFinal: 10000.00},
	}
	for _, testCase := range testCases {
		var offerDetails []databases.CoreOffer
		for _, offerID := range testCase.offerIDs {
			o, _ := testQueries.SelectOfferByID(context.Background(), offerID)
			offerDetails = append(offerDetails, o)
		}
		result, err := Pricing(testQueries, context.Background(), testCase.request, offerDetails)
		fmt.Println(result)
		if result.TotalFinalPrice != testCase.expectFinal{
			t.Errorf(testCase.name, ": ", err)
			require.Equal(t, testCase.expectFinal, result.TotalFinalPrice)
		}
	}

}
