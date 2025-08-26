package redisdriver

import (
	"testing"

	"github.com/cidekar/adele-framework/cache"
)

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache when it was never set.")
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo not found in cache when it was set.")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	x, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if x != "bar" {
		t.Error("did not get value when it was set.")
	}
}

func TestRedisCache_Forget(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("Foo found in cache when it was not there.")
	}

}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("foo2", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("Baz", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("Foo found in cache when it was not there.")
	}

	inCache, err = testRedisCache.Has("foo2")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("Foo found in cache when it was not there.")
	}

	inCache, err = testRedisCache.Has("Baz")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("Baz not found in cache when it should be there.")
	}

}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("Foo found in cache when it was not there.")
	}

}

func TestRedisCache_Decode(t *testing.T) {
	entry := cache.Entry{}
	entry["foo"] = "bar"
	bytes, err := cache.Encode(entry)
	if err != nil {
		t.Error(err)
	}

	_, err = cache.Decode(bytes)
	if err != nil {
		t.Error(err)
	}
}
