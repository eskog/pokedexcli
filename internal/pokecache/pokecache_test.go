package pokecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	testdata := []byte("test")
	cache := NewCache(time.Second * 3)
	cache.Add("test", testdata)
	result, exist := cache.Get("test")
	if assert.NotNil(t, exist) {
		assert.Equal(t, testdata, result)
		t.Logf("Got expected data %v", result)
	} else {
		t.Errorf("failed to get cache data")
	}
	time.Sleep(time.Second * 40)
	result, exist = cache.Get("test")
	if exist {
		t.Errorf("Cache has not been cleaned, got data %v", result)
	} else {
		t.Log("Cache has been cleaned. Test entry no longer active")
	}

}
