package bus

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/infinit-lab/yolanda/config"
	l "github.com/infinit-lab/yolanda/logutils"
	"sync"
	"time"
)

type center struct {
	mutex  sync.Mutex
	routes map[int]*list.List
}

type message struct {
	key   int
	value *Resource
}

type node struct {
	subscriber Subscriber
	channel    chan interface{}
}

func (c *center) subscribe(key int, subscriber Subscriber) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	channels, ok := c.routes[key]
	if ok {
		for temp := channels.Front(); temp != nil; temp = temp.Next() {
			n, ok := temp.Value.(*node)
			if !ok {
				continue
			}
			if n.subscriber == subscriber {
				return
			}
		}
	} else {
		channels = list.New()
		c.routes[key] = channels
	}
	n := new(node)
	n.subscriber = subscriber
	cacheNum := config.GetInt("bus.cacheNum")
	l.Trace("Bus cache number is ", cacheNum)
	if cacheNum == 0 {
		n.channel = make(chan interface{})
	} else {
		n.channel = make(chan interface{}, cacheNum)
	}
	channels.PushBack(n)
	go func() {
		for {
			data, ok := <-n.channel
			if !ok {
				break
			}
			msg, ok := data.(*message)
			if !ok {
				continue
			}
			n.subscriber.Handle(msg.key, msg.value)
		}
	}()
}

func (c *center) unsubscribe(key int, subscriber Subscriber) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	channels, ok := c.routes[key]
	if ok {
		for temp := channels.Front(); temp != nil; temp = temp.Next() {
			n, ok := temp.Value.(*node)
			if !ok {
				return
			}
			if n.subscriber == subscriber {
				close(n.channel)
				channels.Remove(temp)
				if channels.Len() == 0 {
					delete(c.routes, key)
				}
				break
			}
		}
	}
}

func (c *center) channels(key int) (*list.List, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tempList, ok := c.routes[key]
	allList, allOk := c.routes[0]
	if ok || allOk {
		channels := list.New()
		if ok {
			channels.PushBackList(tempList)
		}
		if allOk {
			channels.PushBackList(allList)
		}
		return channels, nil
	}
	return nil, errors.New(fmt.Sprintf("No subscriber"))
}

func (c *center) publish(key int, value *Resource) error {
	channels, err := c.channels(key)
	if err != nil {
		l.WarningF("Failed to get channels of %d, error: %s", key, err.Error())
		return err
	}
	for chn := channels.Front(); chn != nil; chn = chn.Next() {
		n, ok := chn.Value.(*node)
		if ok {
			msg := new(message)
			msg.key = key
			msg.value = value
			n.channel <- msg
		} else {
			l.Error("Failed to convert to node")
		}
	}
	return nil
}

func (c *center) publishDelay(key int, value *Resource, delayMs int) error {
	go func() {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		_ = c.publish(key, value)
	}()
	return nil
}
