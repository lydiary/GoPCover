package redisclient

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Client struct {
	client *redis.Client
}

func (client *Client) NewClient(opt *redis.Options) {
	client.client = redis.NewClient(opt)
}

func (client *Client) SetValue(key string, value interface{}) error {
	serializedValue, _ := json.Marshal(value)
	err := client.client.Set(key, string(serializedValue), 0).Err()
	return err
}

func (client *Client) GetValue(key string) (string, error) {
	serializedValue, err := client.client.Get(key).Result()
	return serializedValue, err
}

func (client *Client) DeleteValue(key string) {
	client.client.Del(key)
}

func (client *Client) Keys(pattern string) []string {
	keys := client.client.Keys(pattern)
	return keys.Val()
}
