package cache

import (
	"../config"
	"log"
	"os"
	"testing"
	"time"
)

var c *Cache

const (
	testKey   = "1"
	testValue = 1
)

func TestNewCacheWithConfig(t *testing.T) {
	os.Args = append(os.Args, "cache.lifetime=1")
	config.Exec()
	c = NewCacheWithConfig()
	if c == nil {
		t.Error("Failed to NewCacheWithConfig")
	}
}

func TestCache_Insert(t *testing.T) {
	c.Insert(testKey, testValue)
	v, ok := c.Get(testKey)
	if !ok {
		t.Error("Failed to Insert")
	} else if v.(int) != 1 {
		t.Error("Failed to Get")
	}
	log.Printf("key is %s, value is %v", testKey, v)
	time.Sleep(2 * time.Second)
	v, ok = c.Get(testKey)
	if ok {
		t.Error("Should not Get")
	}
}

func TestCache_Erase(t *testing.T) {
	c.Insert(testKey, testValue)
	c.Erase(testKey)
	_, ok := c.Get(testKey)
	if ok {
		t.Error("Should not Get")
	}
}
