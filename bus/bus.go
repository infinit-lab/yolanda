package bus

import (
	"container/list"
	"log"
)

type Resource struct {
	Status  int
	Id      string
	Data    interface{}
	Context interface{}
}

type Subscriber interface {
	Handle(key int, value *Resource)
}

var c *center

func init() {
	log.Println("Initializing bus...")
	c = new(center)
	c.routes = make(map[int]*list.List)
}

func PublishResource(key int, status int, id string, data interface{}, context interface{}) error {
	r := new(Resource)
	r.Status = status
	r.Id = id
	r.Data = data
	r.Context = context
	return c.publish(key, r)
}

func PublishResourceDelay(key int, status int, id string, data interface{}, context interface{}, delayMs int) error {
	r := new(Resource)
	r.Status = status
	r.Id = id
	r.Data = data
	r.Context = context
	return c.publishDelay(key, r, delayMs)
}

func Subscribe(key int, subscriber Subscriber) {
	c.subscribe(key, subscriber)
}

func Unsubscribe(key int, subscriber Subscriber) {
	c.unsubscribe(key, subscriber)
}
