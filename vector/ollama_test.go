package vector

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOllama_Embed(t *testing.T) {
	content, err := NewOllama().Embed(context.Background(), "hello")
	assert.NoError(t, err)
	fmt.Println(len(content))
}
