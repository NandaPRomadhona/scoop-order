package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"os"
	app2 "scoop-order/cmd/api"
	"scoop-order/internal/configs"
	"scoop-order/repository/transactions"
	"time"
)

const(
	serverAddress = "0.0.0.0:9999"
)

var (
	Logger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	NewRelicApp *newrelic.Application
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
	Logger = log.Default()

	NewRelicApp, _ = newrelic.NewApplication(
		newrelic.ConfigAppName("scoop-order"),
		newrelic.ConfigLicense("d9ef899b2a9a369fbd28a7a67193806f433cNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	//InfoLogger.SetPrefix("INFO: ")
	//WarningLogger = log.Default()
	//WarningLogger.SetPrefix("WARNING: ")
	//ErrorLogger = log.Default()
	//ErrorLogger.SetPrefix("ERROR: ")
}

func LogFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage)
}


func main()  {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", configs.DBHost, configs.DBPort, configs.DBUser, configs.DBPassword, configs.DBName)
	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot connect to database: %s", err.Error()), nil)
	}

	InfoLogger.Println("Connect to Databases")
	Logger.SetPrefix("INFO - ")
	Logger.Println("Connect to Databases")

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddr,
		Password: configs.RedisPass,
		DB:       0,
	})
	// create router dependency
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(nrgin.Middleware(NewRelicApp))
	//router.Use(gin.LoggerWithFormatter(LogFormatter))

	transaction := transactions.NewTransaction(conn, clientRedis)
	server := app2.NewServer(router, transaction, Logger, WarningLogger, InfoLogger, ErrorLogger)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot start the server: %s", err.Error()), nil)
	}
	InfoLogger.Println("Server Run in Address: "+serverAddress)

}