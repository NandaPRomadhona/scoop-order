package playgroud

import (
	"encoding/json"
	"time"
)

var VAResponse = `{
  "data": {
    "status_code": "201",
    "status_message": "Success, Bank Transfer transaction is created",
    "transaction_id": "d6199bfe-10fd-433e-b38a-8933f5bd5a47",
    "transaction_status": "pending",
    "transaction_time": "2021-11-01 15:04:00",
    "payment_type": "bank_transfer",
    "fraud_status": "accept",
    "order_id": "284239sdsds4jk",
    "va_number": {
      "bank": "bca",
      "va_number": "06204606600"
    }
  },
  "message": "Success, Bank Transfer transaction is created",
  "statusCode": "201"
}`

func GetVAResponse(arg string) string {
	var newRespPayment map[string]interface{}
	var newArg map[string]interface{}

	json.Unmarshal([]byte(VAResponse),&newRespPayment)
	json.Unmarshal([]byte(arg),&newArg)
	newRespPayment["transaction_time"] = time.Now()
	newRespPayment["order_id"] = newArg["order_id"]

	resp, _ := json.Marshal(newRespPayment)
	return string(resp)
}

var GOPayResponse = `{
  "data": {
    "status_code": "201",
    "status_message": "GoPay transaction is created",
    "transaction_id": "04d375b8-cc9b-4df4-ad40-3cd6a17b7620",
    "order_id": "sdsda8927349",
    "redirect_url": "",
    "gross_amount": "7000.00",
    "currency": "IDR",
    "payment_type": "gopay",
    "transaction_time": "",
    "transaction_status": "pending",
    "fraud_status": "",
    "Actions": [
      {
        "name": "generate-qr-code",
        "method": "GET",
        "url": "https://api.sandbox.veritrans.co.id/v2/gopay/04d375b8-cc9b-4df4-ad40-3cd6a17b7620/qr-code",
        "fields": null
      },
      {
        "name": "deeplink-redirect",
        "method": "GET",
        "url": "https://simulator.sandbox.midtrans.com/gopay/partner/app/payment-pin?id=341290c9-3f81-4c67-82ab-58d81c3d8b1b",
        "fields": null
      },
      {
        "name": "get-status",
        "method": "GET",
        "url": "https://api.sandbox.veritrans.co.id/v2/04d375b8-cc9b-4df4-ad40-3cd6a17b7620/status",
        "fields": null
      },
      {
        "name": "cancel",
        "method": "POST",
        "url": "https://api.sandbox.veritrans.co.id/v2/04d375b8-cc9b-4df4-ad40-3cd6a17b7620/cancel",
        "fields": null
      }
    ]
  },
  "message": "GoPay transaction is created",
  "statusCode": "201"
}`

var VASuccessResponse = `{
  "data": {
    "id": 0,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "status_code": "200",
    "status_message": "Success, transaction is found",
    "signature_key": "",
    "bank": "",
    "fraud_status": "accept",
    "payment_type": "bank_transfer",
    "order_id": "",
    "transaction_id": "d6199bfe-10fd-433e-b38a-8933f5bd5a47",
    "transaction_status": "settlement",
    "gross_amount": "7000.00",
    "masked_card": "",
    "currency": "",
    "card_type": "",
    "channel_response_code": "",
    "channel_response_message": "",
    "approval_code": ""
  },
  "message": "Success, transaction is found",
  "statusCode": "200"
}`
