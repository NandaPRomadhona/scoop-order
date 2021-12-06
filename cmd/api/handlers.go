package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"scoop-order/cmd/src"
	"scoop-order/internal/configs"
	"scoop-order/internal/schemas"
	"scoop-order/playgroud"
	"scoop-order/repository"
	"strconv"
	"strings"
	"time"
)

func (server *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		permission, _ := c.Get("permission")
		if permission == "denied" {
			return
		}
		c.Header("Content-Type", "application/json")
		response := map[string]string{
			"status": "success",
			"data":   "order API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (server *Server) GetOrderList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req schemas.GetOrderListRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" GetOrderList : %s", err.Error())))
			return
		}
		var arg = repository.SelectOrderParams{
			Limit:  req.Limit,
			Offset: req.Offset,
		}
		orders, err := server.transaction.SelectOrder(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		result := repository.DestructorsOrders(orders)

		ctx.IndentedJSON(http.StatusOK, result)
	}
}

func (server *Server) GetOrderByID(ctx *gin.Context) {
	var req schemas.GetOrderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" GetOrderList : %s", err.Error())))
		return
	}

	order, err := server.transaction.SelectOrderByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderDetails, err := server.transaction.SelectDetailOrderByID(ctx, order.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := repository.DestructorOrderDetail(order, orderDetails)

	ctx.IndentedJSON(http.StatusOK, result)
}

func (server *Server) GetOrderByOrderNumber(ctx *gin.Context) {
	var req schemas.GetOrderByONRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" GetOrderList : %s", err.Error())))
		return
	}

	order, err := server.transaction.SelectOrderByOrderNumber(ctx, req.ID)
	fmt.Println("order", order.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderDetails, err := server.transaction.SelectDetailOrderByID(ctx, order.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result := repository.DestructorOrderDetail(order, orderDetails)

	ctx.IndentedJSON(http.StatusOK, result)
}

func (server *Server) GetOrderByUser(ctx *gin.Context) {
	server.requireAuthentication()

	user, _ := ctx.Get("user")
	userData, ok := user.(repository.SelectUserRows)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf(" Invlaid Session ")))
	}

	orders, err := server.transaction.SelectOrderByUserID(ctx, userData.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var result []interface{}
	for _, order := range orders {
		orderDetails, err := server.transaction.SelectDetailOrderByID(ctx, order.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		result = append(result, repository.DestructorOrderDetail(order, orderDetails))
	}

	ctx.IndentedJSON(http.StatusOK, result)
}

func (server *Server) GetValidPrice(ctx *gin.Context) {
	var req schemas.GetValidPriceRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" GetOrderList : %s", err.Error())))
		return
	}

	offerIDs := strings.Split(req.OfferID, ",")
	var ids []int32
	for _, id := range offerIDs {
		i, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong offers input")))
		}
		ids = append(ids, int32(i))
	}
	checkPricingParam := schemas.CheckPricingRequest{
		UserID:           req.UserID,
		OfferID:          ids,
		PaymentGatewayID: req.PaymentGatewayID,
		PlatformID:       req.PlatformID,
		DiscountCode:     req.DiscountCode,
		CurrencyCode:     req.CurrencyCode,
		GeoInfo:          schemas.GeoInfo{CountryCode: req.CountryCode},
	}
	pricingOffer, err := server.transaction.PricingTx(ctx, checkPricingParam)
	if err != nil {
		ctx.IndentedJSON(http.StatusConflict, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, successResponse(err.Error(), pricingOffer))
	return
}

func (server *Server) Checkout(ctx *gin.Context) {
	server.requireAuthentication()

	server.Logger.SetPrefix("DEBUG - ")
	server.Logger.Println("=========== Request =============")
	startTime := time.Now()

	var req schemas.CheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" Checkout : %s", err.Error())))
		return
	}
	server.Logger.Println("Data Request: ", req)
	user, _ := ctx.Get("user")
	userData, ok := user.(repository.SelectUserRows)
	server.Logger.Println("userData: ", userData)
	if !ok {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" Checkout : user session is not found")))
		return
	}

	apiKey := ctx.Request.Header["Signature"][0]
	signatureHash := src.GenerateSignature(userData.UserID, req.OfferID, req.PaymentGatewayID)

	if apiKey != hex.EncodeToString(signatureHash[:]) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" Checkout : invalid signature")))
		return
	}

	arg := schemas.CheckoutTxParams{
		UserID:           userData.UserID,
		OfferID:          req.OfferID,
		DiscountCode:     req.DiscountCode,
		CurrencyCode:     req.CurrencyCode,
		PlatformID:       req.PlatformID,
		PaymentGatewayID: req.PaymentGatewayID,
		GeoInfo:          req.GeoInfo,
	}

	// Create orders (NEW)
	checkout, err := server.transaction.CheckoutTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("checkoutTX : %s", err.Error())))
		return
	}
	var payment map[string]interface{}
	var paymentStatus int32
	if checkout.FinalAmount == 0{
		checkout.OrderStatus = configs.OrderComplete
		paymentStatus = configs.PaymentBilled
		server.Logger.SetPrefix("DEBUG - ")
		server.Logger.Println("=========== Hit Owned Item =============")
		// Hit Owned Items Service
		// ...
	} else {
		// Do Purchasing (WAITING_FOR_PAYMENT = 20001)
		// Hit Payment service
		server.Logger.SetPrefix("DEBUG - ")
		server.Logger.Println("=========== Hit Payment =============")
		paymentGateway, _ := server.transaction.SelectPaymentGateway(ctx, req.PaymentGatewayID)
		payment, err = server.Payment(checkout, paymentGateway)
		paymentStatus = configs.PaymentInProcess

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("hit Payment : %s", err.Error())))
			return
		}
		checkout.OrderStatus = configs.OrderWaitingForPayment
	}

	// Update Order data (WAITING_FOR_PAYMENT = 20001)
	err = server.transaction.PaymentTx(ctx, checkout, paymentStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("PaymentTX : %s", err.Error())))
		return
	}

	response := schemas.OrderResponse{
		Order:   checkout,
		Payment: payment,
	}

	server.Logger.Println("=========== End =============")
	server.Logger.Println("Complete in :", time.Since(startTime))
	server.Logger.Println("=========== Response =============")
	server.Logger.Println(response)
	ctx.JSON(http.StatusOK, successResponse("Success Create Order", response))
	return

}

func (server *Server) Payment(checkout schemas.CheckoutTxResult, paymentGateway repository.SelectPaymentGatewaysRow) (map[string]interface{}, error) {
	type payRes struct {
		Data       map[string]interface{} `json:"data"`
		Message    string      `json:"message"`
		statusCode string      `json:"statusCode"`
	}
	var newRespPayment payRes
	user, _ := server.transaction.SelectUser(context.Background(), checkout.UserID)
	baseURL := configs.URLScoopPayment
	method := "POST"
	var requestByte []byte
	var respPayment []byte
	var err error

	if paymentGateway.PaymentGroup == "midtrans" {
		if strings.Contains(strings.ToLower(paymentGateway.Name), "va") {
			url := baseURL + "/va"
			bankName := strings.TrimPrefix(strings.ToLower(paymentGateway.Name), "va bank ")
			var items []map[string]interface{}
			for _, ol := range checkout.Orderline {
				offer,_ := server.transaction.SelectOfferByID(context.Background(), ol.OfferID)
				olItem := map[string]interface{}{
					"id" : fmt.Sprintf("%d", ol.ID),
					"price": ol.FinalPrice,
					"quantity":1,
					"name":offer.Name.String,
				}
				items = append(items, olItem)
			}
			requestByte, err = json.Marshal(map[string]interface{}{
				"order_id": fmt.Sprintf("%d",checkout.OrderNumber),
				"bank":    bankName,
				"email":   user.Email,
				"username": user.UserName,
				"amount":  checkout.FinalAmount,
				"items": items,
			})
			if err != nil {
				log.Fatal(err)
			}

			client := &http.Client{}
			reqPayment, err := http.NewRequest(method, url, bytes.NewReader(requestByte))
			if err != nil {
				fmt.Println("Error 1:", err)
			}
			reqPayment.Header.Add("Authorization", base64.StdEncoding.EncodeToString([]byte(configs.AuthCoopPayment)))
			reqPayment.Header.Add("Content-Type", "application/json")

			server.Logger.Println("Request Body (payment) :", string(requestByte))

			res, err := client.Do(reqPayment)
			if err != nil {
				fmt.Println("Error 2:", err)
			}

			respPayment, err = ioutil.ReadAll(res.Body)
			server.Logger.Println("Response Body (payment) :", string(respPayment))

			if res.StatusCode != 201{
				return newRespPayment.Data, fmt.Errorf(" Error Payments : %s", string(respPayment))
			}

			defer res.Body.Close()

			if err != nil {
				fmt.Println("Error 3:", err)
			}
		}
	}

	json.Unmarshal(respPayment, &newRespPayment)
	return newRespPayment.Data, nil
}

func (server *Server) Complete(ctx *gin.Context) {
	server.Logger.SetPrefix("DEBUG - ")
	server.Logger.Println("=========== Request =============")
	startTime := time.Now()

	var req schemas.CompleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(" GetOrderList : %s", err.Error())))
		return
	}
	server.Logger.Println("Data Request: ", req)

	order, err := server.transaction.SelectOrderByOrderNumber(context.Background(), req.OrderID)
	if err != nil {
		ctx.IndentedJSON(http.StatusConflict, err)
	}
	fmt.Println("Complete Order ID: ", order.ID)

	CompletePaymentTxParams := schemas.PaymentTxParams{
		OrderID:       order.ID,
		OrderStatus:   int32(configs.OrderComplete),
		PaymentStatus: int32(configs.PaymentBilled),
	}
	status, err := server.transaction.CompletePaymentTx(context.Background(), CompletePaymentTxParams)
	if err != nil {
		log.Fatal(fmt.Errorf("[Failed] Update Complete Data: %s", err.Error()))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	server.Logger.Println("=========== End Request =============")
	server.Logger.Println("Complete in :", time.Since(startTime))
	ctx.JSON(http.StatusOK, successResponse(" Complete Order Success", status))
}

func (server *Server) CompletePayment() {
	// Get Pending Order
	//arg := repository.SelectOrderParams{
	//	Limit:  10,
	//	Offset: 10,
	//}
	//orders, _ := server.transaction.SelectOrder(context.Background(), arg)
	order, _ := server.transaction.SelectOrderByOrderNumber(context.Background(), 32021612217630)

	//for _, order := range orders {
		var currentOrderStatus int32
		// Check Payment Status
		// url : http://localhost:5000/check-status
		type reqParam struct {
			orderID string `json:"order_id"`
		}
		request := reqParam{
			orderID: fmt.Sprintf("%d", order.OrderNumber),
		}
		requestJson, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(requestJson))
		//resp, err := http.Post("http://localhost:5000/va", "application/json",
		//	bytes.NewBuffer(requestJson))

		var newRespPayment map[string]interface{}

		//json.NewDecoder(respPayment.Body).Decode(&newRespPayment)
		json.Unmarshal([]byte(playgroud.VASuccessResponse), &newRespPayment)
		if newRespPayment["transaction_status"] == "settlement" {
			currentOrderStatus = 90000
			CompletePaymentTxParams := schemas.PaymentTxParams{
				OrderID:       order.ID,
				OrderStatus:   currentOrderStatus,
				PaymentStatus: currentOrderStatus,
			}
			_, err := server.transaction.CompletePaymentTx(context.Background(), CompletePaymentTxParams)
			if err != nil {
				log.Fatal(fmt.Errorf("[Failed] Update Complete Data: %s", err.Error()))
			}
		}
	}
//}

func (server *Server) GetPaymentGateway() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paymentGateways, err := server.transaction.SelectPaymentGateways(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.IndentedJSON(http.StatusOK, paymentGateways)
	}
}

