package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadStringToken(t *testing.T) {
	tests := []struct {
		input    string
		token    string
		expected string
	}{
		{"Hello, world!<?", "<?", "Hello, world!"},
		{"Hello, world!?>", "?>", "Hello, world!"},
		{"Hello, world!<?More text", "<?", "Hello, world!"},
		{"Hello, world!?>More text", "?>", "Hello, world!"},
		{"Hello, world!", "<?", "Hello, world!"},
	}

	for _, test := range tests {
		reader := bufio.NewReader(bytes.NewBufferString(test.input))
		result, err := ReadStringToken(reader, test.token)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}
