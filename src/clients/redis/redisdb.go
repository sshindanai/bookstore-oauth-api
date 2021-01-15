package redis

import (
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	// var err error
	// conn, err = redis.DialURL("redis://localhost")
	// defer conn.Close()
	// if err != nil {
	// 	panic(err)
	// }
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

}

func NewRedis() *redis.Client {
	return rdb
}
