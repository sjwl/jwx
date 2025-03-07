// Code generated by tools/cmd/genjwa/main.go. DO NOT EDIT

package jwa_test

import (
	"testing"

	"github.com/sjwl/jwx/v2/jwa"
	"github.com/stretchr/testify/assert"
)

func TestCompressionAlgorithm(t *testing.T) {
	t.Parallel()
	t.Run(`accept jwa constant Deflate`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept(jwa.Deflate), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.Deflate, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`accept the string DEF`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept("DEF"), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.Deflate, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`accept fmt.Stringer for DEF`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept(stringer{src: "DEF"}), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.Deflate, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`stringification for DEF`, func(t *testing.T) {
		t.Parallel()
		if !assert.Equal(t, "DEF", jwa.Deflate.String(), `stringified value matches`) {
			return
		}
	})
	t.Run(`accept jwa constant NoCompress`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept(jwa.NoCompress), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.NoCompress, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`accept the string `, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept(""), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.NoCompress, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`accept fmt.Stringer for `, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.NoError(t, dst.Accept(stringer{src: ""}), `accept is successful`) {
			return
		}
		if !assert.Equal(t, jwa.NoCompress, dst, `accepted value should be equal to constant`) {
			return
		}
	})
	t.Run(`stringification for `, func(t *testing.T) {
		t.Parallel()
		if !assert.Equal(t, "", jwa.NoCompress.String(), `stringified value matches`) {
			return
		}
	})
	t.Run(`bail out on random integer value`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.Error(t, dst.Accept(1), `accept should fail`) {
			return
		}
	})
	t.Run(`do not accept invalid (totally made up) string value`, func(t *testing.T) {
		t.Parallel()
		var dst jwa.CompressionAlgorithm
		if !assert.Error(t, dst.Accept(`totallyInvfalidValue`), `accept should fail`) {
			return
		}
	})
	t.Run(`check list of elements`, func(t *testing.T) {
		t.Parallel()
		var expected = map[jwa.CompressionAlgorithm]struct{}{
			jwa.Deflate:    {},
			jwa.NoCompress: {},
		}
		for _, v := range jwa.CompressionAlgorithms() {
			if _, ok := expected[v]; !assert.True(t, ok, `%s should be in the expected list`, v) {
				return
			}
			delete(expected, v)
		}
		if !assert.Len(t, expected, 0) {
			return
		}
	})
}
