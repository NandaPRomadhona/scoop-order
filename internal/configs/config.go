package configs

import (
	"os"
	"time"
)

var DBUser = os.Getenv("DB_USER")
var DBPassword = os.Getenv("DB_PASS")
var DBHost = os.Getenv("DB_HOST")
var DBPort = os.Getenv("DB_PORT")
var DBName = os.Getenv("DB_NAME")
var RedisAddr = os.Getenv("REDIS_ADDR")
var RedisPass = os.Getenv("REDIS_PASS")
var JWTSecrete = os.Getenv("JWT_SECRET")
var AuthCoopPayment = os.Getenv("AUTH_SCOOP_PAYMENT")
var URLScoopPayment = os.Getenv("URL_SCOOP_PAYMENT")
var RedisDuration = time.Duration(2)

