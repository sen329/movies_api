package main

import (
	"database/sql"
	controller "movies/Controller"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var testEnvironment string = os.Getenv("TEST_ENVIRONMENT")

func main() {
	var database string
	driver := os.Getenv("DB_DRIVER")
	username := os.Getenv("DB_USER_NAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS")
	if strings.Contains(testEnvironment, "staging") {
		database = os.Getenv("DB_DATABASE_TEST")
	} else {
		database = os.Getenv("DB_DATABASE")
	}
	dbase, err := sql.Open(driver, username+":"+password+"@tcp("+address+")"+"/"+database+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		if err := dbase.Close(); err != nil {
			panic(err)
		}
	}()

	var redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	cache := cache.New(&cache.Options{
		Redis:      redis,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	defer func() {
		if err := redis.Close(); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	movies := &controller.MovieController{
		Db:    dbase,
		Cache: cache,
	}

	movies.MapEndpoints(router)

	router.Run(":8080")

}
