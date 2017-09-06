package util

import (
	"github.com/bysir-zl/bygo/cache"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
)

type MapStorage struct {
	id    string
	data  map[string]string
	redis *cache.BRedis
}

func (p *MapStorage) Set(key, value string) {
	p.data[key] = value
	p.redis.HMSET("", key, value, 0)
}

func (p *MapStorage) Get(key string) (value string, ok bool) {
	value, ok = p.data[key]
	return
}
func (p *MapStorage) Del(key string) {
	delete(p.data, key)
	p.redis.HDEL("", key)
}

func (p *MapStorage) Save() {
	// todo save to db
}

func (p *MapStorage) Load() {
	// todo load from db
	data, err := p.redis.HGETALL("")
	if err != nil {
		log.ErrorT("MapStorage", p.id+" error:", err)
		return
	}
	p.data = make(map[string]string, len(data))
	for key, value := range data {
		if v, ok := util.Interface2String(value, true); ok {
			p.data[key] = v
		}
	}
}

func NewMapStorage(id string) *MapStorage {
	redis := cache.NewRedis("127.0.0.1:6379")
	redis.SetPrefix(id)
	return &MapStorage{
		id:    id,
		data:  map[string]string{},
		redis: redis,
	}
}
