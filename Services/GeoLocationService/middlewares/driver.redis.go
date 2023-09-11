package middlewares

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client = nil

var ctx = context.Background()

func ConnectToRedis() {
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")

	if port == "" {
		port = "6785"
	}
	if host == "" {
		host = "localhost"
	}
	url := fmt.Sprintf("%s:%s", host, port)
	log.Println("Redis url: ", url)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
}

func RedisFindDriver(lon float64, lat float64, geo_key string, min_km float64) ([]redis.GeoLocation, error) {
	res, err := RedisClient.GeoRadius(ctx, geo_key, lon, lat, &redis.GeoRadiusQuery{
		Radius: min_km,
		Unit:   "km",
		Count:  10,
		Sort:   "ASC",
	}).Result()
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func RedisUpdateDriver(lon float64, lat float64, geo_key string, driver_id string) error {
	_, err := RedisClient.GeoAdd(ctx, geo_key, &redis.GeoLocation{
		Name:      driver_id,
		Longitude: lon,
		Latitude:  lat,
	}).Result()
	return err
}
