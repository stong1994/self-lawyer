package vector

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type CacheOption func(*Cache)

var WithCacheOptionSetCachePath = func(path string) CacheOption {
	return func(c *Cache) {
		c.cachePath = path
	}
}

type Cache struct {
	cachePath string
	file      *os.File
	cache     map[string][]float32
}

func NewCache(opts ...CacheOption) *Cache {
	config := &Cache{
		cachePath: "cache.json",
	}
	for _, opt := range opts {
		opt(config)
	}
	file, err := os.OpenFile(config.cachePath, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("open cache file failed", err)
	}
	config.file = file

	cache := make(map[string][]float32)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), " ")
		embedding := make([]float32, 0, len(splits)-1)
		for i := 1; i < len(splits); i++ {
			f, err := strconv.ParseFloat(splits[i], 32)
			if err != nil {
				log.Fatal("parse float failed", splits, err)
			}
			embedding = append(embedding, float32(f))
		}
		cache[splits[0]] = embedding
	}
	config.cache = cache

	return config
}

func (c *Cache) Get(key string) []float32 {
	if embedding, ok := c.cache[key]; ok {
		return embedding
	}
	return nil
}

func (c *Cache) Set(key string, embedding []float32) {
	c.cache[key] = embedding
	if _, ok := c.cache[key]; ok {
		// update the row in the cache file with key matching
		c.Update(key, embedding)
	} else {
		// just append the new line to the file
		writer := bufio.NewWriter(c.file)
		if _, err := writer.WriteString(cacheLine(key, embedding) + "\n"); err != nil {
			log.Fatal("write cache line failed", err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Cache) Close() {
	if err := c.file.Close(); err != nil {
		log.Fatal("close cache file failed", err)
	}
}

func cacheLine(key string, embedding []float32) string {
	line := key
	for _, f := range embedding {
		line += " " + strconv.FormatFloat(float64(f), 'f', -1, 32)
	}
	return line
}

// Update updates the row in the cache file with key matching.
func (c *Cache) Update(key string, embedding []float32) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(c.file)
	// remove the old line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 0 && parts[0] == key {
			continue
		}
		lines = append(lines, line)
	}
	// append the new line
	lines = append(lines, cacheLine(key, embedding))

	if err := scanner.Err(); err != nil {
		log.Fatal("read cache file failed", err)
	}

	// Truncate the file to 0 size
	if err := c.file.Truncate(0); err != nil {
		log.Fatal("truncate cache file failed", err)
	}
	// Reset the file offset to the beginning
	if _, err := c.file.Seek(0, 0); err != nil {
		log.Fatal("seek cache file failed", err)
	}

	writer := bufio.NewWriter(c.file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			log.Fatal("write cache line failed", err)
		}
	}

	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}
}
