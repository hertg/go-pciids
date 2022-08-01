package conv_test

import (
	"testing"

	"github.com/go-playground/assert"
	"github.com/hertg/go-pciids/internal/conv"
)

func TestParseByteNum(t *testing.T) {
	num := conv.ParseByteNum([]byte("2f"))
	assert.Equal(t, uint(47), num)

	num = conv.ParseByteNum([]byte("f4ab"))
	assert.Equal(t, uint(62635), num)

	num = conv.ParseByteNum([]byte("c48118"))
	assert.Equal(t, uint(12878104), num)

	num = conv.ParseByteNum([]byte("df29fe6e"))
	assert.Equal(t, uint(3744071278), num)
}
