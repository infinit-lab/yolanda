package bus

import (
	l "github.com/infinit-lab/yolanda/logutils"
	"testing"
	"time"
)

type handler struct {
}

const (
	testKey    = 1
	testStatus = 2
	testData   = 3
	testId     = "123"
)

func (h *handler) Handle(key int, resource *Resource) {
	l.TraceF("key is %d, status is %d, id is %s, data is %d", key, resource.Status, resource.Id, resource.Data.(int))
	t, ok := resource.Context.(*testing.T)
	if !ok {
		l.Error("Failed to convert to testing.T")
		return
	}
	if key != testKey {
		t.Errorf("Key should be %d, actual %d", testKey, key)
	}
	if resource.Id != testId {
		t.Errorf("Id should be %s, actual %s", testId, resource.Id)
	}
	if resource.Status != testStatus {
		t.Errorf("Status should be %d, actual %d", testStatus, resource.Status)
	}
	d := resource.Data.(int)
	if d != testData {
		t.Errorf("Data should be %d, actual %d", testData, d)
	}
}

var h handler

func TestSubscribe(t *testing.T) {
	Subscribe(testKey, &h)
}

func TestPublishResource(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := PublishResource(testKey, testStatus, testId, testData, t)
		if err != nil {
			t.Error("Failed to PublishResource")
		}
	}
	time.Sleep(100 * time.Millisecond)
}

func TestUnsubscribe(t *testing.T) {
	Unsubscribe(testKey, &h)
}
