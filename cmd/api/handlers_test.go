package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"scoop-order/internal/schemas"
	"scoop-order/repository/transactions"
	"testing"
)

func TestServer_GetOrderList(t *testing.T) {
	testCases := []struct {
		name          string
		arg           schemas.GetOrderListRequest
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "No Arguments",
			arg: schemas.GetOrderListRequest{},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "No Offset",
			arg: schemas.GetOrderListRequest{Limit: 5},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "All Arguments",
			arg: schemas.GetOrderListRequest{Limit: 20, Offset: 5},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	trx := transactions.NewTransaction(testDB, testRedis)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, trx)
			orderListPath := "/orders"
			server.router.GET(
				orderListPath,
				server.GetOrderList(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			bodyByte, err := json.Marshal(testCases[i].arg)
			body := bytes.NewReader(bodyByte)
			request, err := http.NewRequest(http.MethodGet, orderListPath, body)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			//var responseObject []repository.OrderPublished
			//bodyBytes, err := ioutil.ReadAll(recorder.Body)
			//fmt.Println(string(bodyBytes))
			//err = json.Unmarshal(bodyBytes, &responseObject)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestServer_GetValidPrice(t *testing.T) {
	testCases := []struct {
		name          string
		//arg           schemas.GetValidPriceRequest
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "No Discount",
			//arg: schemas.GetValidPriceRequest{UserID: 2217630, }
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	trx := transactions.NewTransaction(testDB, testRedis)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, trx)
			finalPricePath := "/finalPrice"
			server.router.GET(
				finalPricePath,
				server.GetValidPrice,
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, finalPricePath+"?user_id=2217630&offer_id=26785&payment_gateway_id=1&platform_id=1&currency_code=IDR&country_code=ID", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
