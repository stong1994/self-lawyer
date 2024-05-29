package vector

import "context"

type CacheOllama struct {
	o     *Ollama
	cache *Cache
}

func NewCacheOllama(o *Ollama, cache *Cache) *CacheOllama {
	c := &CacheOllama{
		o:     o,
		cache: cache,
	}
	return c
}

func (c *CacheOllama) Embed(ctx context.Context, content string) ([]float32, error) {
	if embedding := c.cache.Get(content); embedding != nil {
		return embedding, nil
	}
	embedding, err := c.o.Embed(ctx, content)
	if err != nil {
		return nil, err
	}

	return embedding, nil
}

func (c *CacheOllama) GetDim() int {
	return c.o.GetDim()
}
