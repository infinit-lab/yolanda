package cache

import (
	"github.com/infinit-lab/yolanda/config"
	l "github.com/infinit-lab/yolanda/logutils"
	"sync"
	"time"
)

type status int

const (
	statusHot  = 1
	statusCold = 2
)

type data struct {
	value interface{}
	s     status
}

type Cache struct {
	mutex     sync.Mutex
	cache     map[string]*data
	lifetimeS int
}

func NewCacheWithConfig() *Cache {
	lifetime := config.GetInt("cache.lifetime")
	l.Trace("Life time is ", lifetime)
	if lifetime == 0 {
		lifetime = 30 * 60
		l.Trace("Life time reset to ", lifetime)
	}
	return NewCache(lifetime)
}

func NewCache(lifetimeS int) *Cache {
	cache := new(Cache)
	cache.cache = make(map[string]*data)
	cache.lifetimeS = lifetimeS
	go func() {
		for {
			time.Sleep(time.Duration(cache.lifetimeS) * time.Second)
			cache.mutex.Lock()
			var keys []string
			for key, value := range cache.cache {
				switch value.s {
				case statusCold:
					keys = append(keys, key)
					break
				case statusHot:
					value.s = statusCold
					break
				}
			}

			for _, key := range keys {
				delete(cache.cache, key)
			}
			cache.mutex.Unlock()
		}
	}()
	return cache
}

func (c *Cache) Insert(key string, value interface{}) {
	c.Erase(key)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	d := new(data)
	d.value = value
	d.s = statusHot
	c.cache[key] = d
}

func (c *Cache) Erase(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	value, ok := c.cache[key]
	if ok {
		value.s = statusHot
		return value.value, ok
	}
	return nil, false
}
