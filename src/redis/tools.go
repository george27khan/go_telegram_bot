package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func getOpts() (opts *redis.Options) {
	// loads DB settings from .env into the system
	//if err := godotenv.Load("./db.env"); err != nil {
	//	slog.Logger.Error("File ./db.env not found")
	//}
	rdHost := os.Getenv("REDIS_HOST")
	rdPort := os.Getenv("REDIS_PORT")
	rdPwd := os.Getenv("REDIS_PASSWORD")

	opts = &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rdHost, rdPort),
		Password: rdPwd,
		DB:       0, // use default DB
	}
	return
}
func Connect() *redis.Client {
	rdb := redis.NewClient(getOpts())
	return rdb
}
