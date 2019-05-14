package redis

import (
	"bytes"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"option.bzza.com/system"
)

var RedisClient *redis.Client
var CacheCodec *cache.Codec

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     system.Conf.Redis.Server,
		Password: system.Conf.Redis.Password, // no password set
		DB:       0,                          // use default DB
	})

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

	CacheCodec = &cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func Store(key string, obj interface{}, ex time.Duration) (err error) {
	err = CacheCodec.Set(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: ex,
	})
	return
}

func Get(key string, obj interface{}) (err error) {
	err = CacheCodec.Get(key, obj)
	return
}

type Cache struct {
	Key    string
	Field  string
	It     interface{} // for cache,address only
	Expire time.Duration
	Buf    bytes.Buffer
	IsArr  bool
}

func (c *Cache) Struct2RedisHSet() (err error) {
	enc := msgpack.NewEncoder(&(*c).Buf).StructAsArray((*c).IsArr)
	err = enc.Encode((*c).It)
	if err != nil {
		return
	}

	rbc := RedisClient.HSet((*c).Key, (*c).Field, (*c).Buf.Bytes())
	if rbc.Err() != nil {
		err = rbc.Err()
		return
	}

	rbc = RedisClient.Expire((*c).Key, (*c).Expire)
	if rbc.Err() != nil {
		err = rbc.Err()
		return
	}
	return
}

func (c *Cache) RedisHGet2Struct() (err error) {
	scmd := RedisClient.HGet((*c).Key, (*c).Field)
	if err = scmd.Err(); err != nil {
		return
	}
	b, err := scmd.Bytes()
	if err != nil {
		return
	}
	_, err = c.Buf.Write(b)
	if err != nil {
		return
	}

	dec := msgpack.NewDecoder(&(*c).Buf)
	err = dec.Decode((*c).It)
	if err != nil {
		return
	}
	//fmt.Println((*c).It)
	return
}
