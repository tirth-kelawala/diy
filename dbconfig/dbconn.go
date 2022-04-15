package dbconfig

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

func CreatePostgresConn() (sqlDb *sql.DB) {
	host := "localhost"
	port := "5432"
	username := "tirth"
	dbName := "diy"

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, username, dbName,
	)

	//testing2

	conn, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}

func CreateRedisConn() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}
