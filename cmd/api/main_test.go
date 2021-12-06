package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"scoop-order/repository"
	"scoop-order/repository/transactions"
	"testing"
)

var testQueries *repository.Queries
var testDB *sql.DB
var testRedis *redis.Client
var err error

func TestMain(m *testing.M) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, user, password, dbName )
	//fmt.Println(dsn)
	testDB, err = sql.Open("postgres", dsn)
	if err !=nil{
		panic ("Cannot connect to DB")
	}
	//fmt.Println(conn)
	testQueries = repository.New(testDB)

	testRedis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, trx transactions.Transaction) *Server {
	router := gin.Default()
	router.Use(cors.Default())

	server := NewServer(router, trx)
	require.NoError(t, err)

	return server
}
