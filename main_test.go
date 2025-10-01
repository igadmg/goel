package main

import (
	"bufio"
	"bytes"
	"strings"
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

func TestProcessBuffer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "StateStream with plain text",
			input: `Hello, World!
<?`,
			expected: "wr.Write([]byte(`Hello, World!\n`))\n",
		},
		{
			name:     "StateCode with Go code",
			input:    `<?fmt.Println("Hello, Go!")?>`,
			expected: `fmt.Println("Hello, Go!")` + "\n",
		},
		{
			name: "StateCode with Go code with *",
			input: `Hello
	*World <?fmt.Println("Hello, Go!")?> `,
			expected: `fmt.Println("Hello, Go!")` + "\n",
		},
		{
			name: "Mixed states",
			input: `Hello, World!
<?fmt.Println("Hello, Go!")?>
Goodbye, World!`,
			expected: `wr.Write([]byte(` + "`Hello, World!\n`" + `))` + "\n" +
				`fmt.Println("Hello, Go!")` + "\n" +
				`wr.Write([]byte(` + "`Goodbye, World!`" + `))` + "\n",
		},
		{
			name:     "StateCode with expression",
			input:    `<?=42?>`,
			expected: `wr.Write([]byte(fmt.Sprintf("%v", 42)))` + "\n",
		},
		{
			name:     "StateCode with title case expression",
			input:    `<?^=hello world?>`,
			expected: `wr.Write([]byte(cases.Title(language.English, cases.NoLower).String(fmt.Sprintf("%v", hello world))))` + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			reader := bufio.NewReader(strings.NewReader(tc.input))

			// Act
			output := process_buffer(reader)

			// Assert
			assert.Equal(t, tc.expected, output.String())
		})
	}
}
