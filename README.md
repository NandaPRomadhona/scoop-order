# scoop-order

Order Service


## Deployment

To deploy this project need some environment

```bash
DB_HOST
DB_USER
DB_PORT
DB_PASS
DB_NAME
REDIS_ADDR
REDIS_PASS
JWT_SECRET
AUTH_SCOOP_PAYMENT
URL_SCOOP_PAYMENT

```


## API Reference

#### Get Final Price

```http
  GET /v1/finalPrice
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `user_id` | `Integer` | **Required**.user's ID |
| `offer_id` | `Integer` | **Required**.Offer's ID|
| `offer_id` | `Integer` | **Required**.Offer's ID|
| `payment_gateway_id` | `Integer` | **Required**|
| `platform_id` | `Integer` | |
| `currency_code` | `String` ||
| `discount_code` | `String` ||
| `country_code` | `String` ||

example: 
```http
  GET /v1/finalPricefinalPrice?user_id=2217630&offer_id=3,1,2&payment_gateway_id=1
  &platform_id=1&currency_code=IDR&discount_code=ALLOWALL5K&country_code=ID
```

#### Checkout

```http
  GET /v1/checkout
```

| Header | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Client-ID`      | `integer` | **Required**. Id of Client |
| `Authorization`      | `integer` | **Required**. JWT Token|
| `Signature`      | `integer` | **Required**. Secret Key from each transaction|

Example JSON Body:
```bigquery
{
	"offer_id" : [43103],
	"currency_code": "IDR",
	"discount_code" : "",
	"platform_id":4,
	"payment_gateway_id": 29,
	"geo_info":{
		"latitude":12.20,
		"longitude":23.23,
		"country_code" :"IN",
		"country_name":"indonesia"
	}
}
```

#### Complete

```http
  GET /v1/complete
```

| Header | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `Client-ID`      | `integer` | **Required**. Id of Client |
| `Authorization`      | `integer` | **Required**. JWT Token|

Example JSON Body:
```bigquery
{
	"order_id" : 32021612217630
}
```
