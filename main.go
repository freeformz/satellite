package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
)

func newPool(server string) *redis.Pool {
    return &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}

func main() {
	router := httprouter.New()

	redisUrlKey := os.Getenv("REDIS_URL_KEY")
	redisURL := os.Getenv(redisUrlKey)
	if redisURL == "" {
		redisURL = "localhost:6379"	
	}
    	
	pool := newPool(redisURL)

	broadcaster := NewBroadcastHandler(pool)
	router.HandlerFunc("GET", "/broadcast/:channel", broadcaster)
	router.HandlerFunc("POST", "/broadcast/:channel", broadcaster)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3001"
	}
	http.ListenAndServe(":"+port, router)
}
