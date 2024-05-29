package vector

import (
	"os"
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	// Setup
	cachePath := "test_cache.txt"
	defer os.Remove(cachePath) // clean up after test

	// Test NewCacheConfig
	cache := NewCache(CacheOptionSetCachePath(cachePath))
	if cache.cachePath != cachePath {
		t.Errorf("Expected cachePath to be %s, got %s", cachePath, cache.cachePath)
	}
	defer cache.Close()

	// Test Set and Get
	key := "testKey"
	value := []float32{1.0, 2.0, 3.0}
	cache.Set(key, value)

	result := cache.Get(key)
	if !reflect.DeepEqual(result, value) {
		t.Errorf("Expected value to be %v, got %v", value, result)
	}

	// Test Update
	newValue := []float32{4.0, 5.0, 6.0}
	cache.Set(key, newValue)

	result = cache.Get(key)
	if !reflect.DeepEqual(result, newValue) {
		t.Errorf("Expected value to be %v, got %v", newValue, result)
	}

	// Test Append
	key2 := "testKey2"
	value2 := []float32{7.0, 8.0, 9.0}
	cache.Set(key2, value2)
	result2 := cache.Get(key2)
	if !reflect.DeepEqual(result2, value2) {
		t.Errorf("Expected value to be %v, got %v", value2, result2)
	}
}
