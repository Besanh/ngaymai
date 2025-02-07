package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	GetClient() *redis.Client
	Connect() error
}

var Redis IRedis

type RedisClient struct {
	Client *redis.Client
	config Config
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedis(config Config) (IRedis, error) {
	r := &RedisClient{
		config: config,
	}
	if err := r.Connect(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}

func (r *RedisClient) Connect() error {
	Client := redis.NewClient(&redis.Options{
		Addr:     r.config.Addr,
		Password: r.config.Password,
		DB:       r.config.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	str, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(str)
	r.Client = Client
	return nil
}
