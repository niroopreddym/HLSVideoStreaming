package services

import (
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis"
	"github.com/niroopreddym/HLSVideoStreaming/helpers"
)

// Redis is struct
type Redis struct {
	RedisClient *redis.Client
}

// NewRedisInstance is the ctor to instantiate the Redis-DB Service
func NewRedisInstance(hostName string, port string) *Redis {
	fmt.Println("Go Redis Tutorial")
	addr := fmt.Sprintf("%v:%v", hostName, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return &Redis{
		RedisClient: client,
	}
}

// AddKeyValuePair adds the key value pair to the Redis DB
func (rdb *Redis) AddKeyValuePair(key string, value interface{}) {
	rdb.RedisClient.Set(key, value, 0)
}

// GetValueByKey fetches the value based on key
func (rdb *Redis) GetValueByKey(key string) string {
	cmd := rdb.RedisClient.Get(key)
	return cmd.Val()
}

// GetValuesInList fetches the list of values based on listkey
func (rdb *Redis) GetValuesInList(listKey string) []string {
	cmd := rdb.RedisClient.LRange(listKey, 0, 50)
	return cmd.Val()
}

// AddKeysToList adds the key value pair to the Redis DB
func (rdb *Redis) AddKeysToList(listName string, value string) {
	if rdb.RedisClient.Exists(listName).Val() == 0 {
		rdb.RedisClient.LPush(listName, value)
	} else {
		rdb.RedisClient.RPush(listName, value)
	}
}

// ContainsKey checks if a key exists or not
func (rdb *Redis) ContainsKey(key string) bool {
	return rdb.RedisClient.Exists(key).Val() == 1
}

// ContainsKey checks if a key exists or not
func (rdb *Redis) DeleteKey(key string) bool {
	return rdb.RedisClient.Del(key).Val() == 1
}

// PlaceFFMPEGDataToRedis converts a file and plaes the data into redis
func (rdb *Redis) PlaceFFMPEGDataToRedis(outDirPath string, inputVideo string) {
	items, _ := ioutil.ReadDir(outDirPath)
	for _, item := range items {
		if item.IsDir() {
			subitems, _ := ioutil.ReadDir(item.Name())
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					fmt.Println(item.Name() + "/" + subitem.Name())
				}
			}
		} else {
			byteArrayItem := helpers.ConvertToByteArray(outDirPath + item.Name())
			rdb.AddKeyValuePair(item.Name(), byteArrayItem)
			rdb.AddKeysToList(inputVideo, item.Name())
		}
	}
}
