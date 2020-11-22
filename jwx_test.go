package jwx_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"strings"
	"testing"

	"github.com/lestrrat-go/jwx"
	"github.com/lestrrat-go/jwx/internal/jose"
	"github.com/lestrrat-go/jwx/internal/json"
	"github.com/lestrrat-go/jwx/internal/jwxtest"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"
)

type jsonUnmarshalWrapper struct {
	buf []byte
}

func (w jsonUnmarshalWrapper) Decode(v interface{}) error {
	return json.Unmarshal(w.buf, v)
}

func TestDecoderSetting(t *testing.T) {
	const src = `{"foo": 1}`

	for _, useNumber := range []bool{true, false} {
		useNumber := useNumber
		t.Run(fmt.Sprintf("jwx.WithUseNumber(%t)", useNumber), func(t *testing.T) {
			if useNumber {
				jwx.DecoderSettings(jwx.WithUseNumber(useNumber))
				t.Cleanup(func() {
					jwx.DecoderSettings(jwx.WithUseNumber(false))
				})
			}

			// json.NewDecoder must be called AFTER the above jwx.DecoderSettings call
			decoders := []struct {
				Name    string
				Decoder interface{ Decode(interface{}) error }
			}{
				{Name: "Decoder", Decoder: json.NewDecoder(strings.NewReader(src))},
				{Name: "Unmarshal", Decoder: jsonUnmarshalWrapper{buf: []byte(src)}},
			}

			for _, tc := range decoders {
				tc := tc
				t.Run(tc.Name, func(t *testing.T) {
					var m map[string]interface{}
					if !assert.NoError(t, tc.Decoder.Decode(&m), `Decode should succeed`) {
						return
					}

					v, ok := m["foo"]
					if !assert.True(t, ok, `m["foo"] should exist`) {
						return
					}

					if useNumber {
						if !assert.Equal(t, json.Number("1"), v, `v should be a json.Number object`) {
							return
						}
					} else {
						if !assert.Equal(t, float64(1), v, `v should be a float64`) {
							return
						}
					}
				})
			}
		})
	}
}

// Test compatibility against `jose` tool
func TestJoseCompatibility(t *testing.T) {
	if testing.Short() {
		t.Logf("Skipped during short tests")
		return
	}

	if !jose.Available() {
		t.Logf("`jose` binary not availale, skipping tests")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("jwk", func(t *testing.T) {
		testcases := []struct {
			Name      string
			Raw       interface{}
			Template  string
			VerifyKey func(context.Context, *testing.T, jwk.Key) bool
		}{
			{
				Name:     "RSA Private Key (256)",
				Raw:      rsa.PrivateKey{},
				Template: `{"alg": "RS256"}`,
			},
			{
				Name:     "RSA Private Key (384)",
				Raw:      rsa.PrivateKey{},
				Template: `{"alg": "RS384"}`,
			},
			{
				Name:     "RSA Private Key (512)",
				Raw:      rsa.PrivateKey{},
				Template: `{"alg": "RS512"}`,
			},
			{
				Name:     "RSA Private Key with Private Parameters",
				Raw:      rsa.PrivateKey{},
				Template: `{"alg": "RS256", "x-jwx": 1234}`,
				VerifyKey: func(ctx context.Context, t *testing.T, key jwk.Key) bool {
					m, err := key.AsMap(ctx)
					if !assert.NoError(t, err, `key.AsMap() should succeed`) {
						return false
					}

					if !assert.Equal(t, float64(1234), m["x-jwx"], `private parameters should match`) {
						return false
					}

					return true
				},
			},
		}

		for _, tc := range testcases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				keyfile, cleanup, err := jose.GenerateJwk(ctx, t, tc.Template)
				if !assert.NoError(t, err, `jose.GenerateJwk should succeed`) {
					return
				}
				defer cleanup()

				webkey, err := jwxtest.ParseJwkFile(ctx, keyfile)
				if !assert.NoError(t, err, `ParseJwkFile should succeed`) {
					return
				}

				if vk := tc.VerifyKey; vk != nil {
					if !vk(ctx, t, webkey) {
						return
					}
				}

				if !assert.NoError(t, webkey.Raw(&tc.Raw), `jwk.Raw should succeed`) {
					return
				}
			})
		}
	})
	t.Run("jwe", func(t *testing.T) {
		expected := []byte("Lorem ipsum")
		t.Run("ECDH", func(t *testing.T) {
			// XXX for ECDH-ES (remove the following t.SkipNow() to enable this test)
			t.SkipNow()

			// let jose generate a key file
			joseJwkFile, joseJwkCleanup, err := jose.GenerateJwk(ctx, t, `{"alg": "ECDH-ES"}`)
			if !assert.NoError(t, err, `jose.GenerateJwk should succeed`) {
				return

			}
			defer joseJwkCleanup()

			// Load the JWK generated by jose
			jwxJwk, err := jwxtest.ParseJwkFile(ctx, joseJwkFile)
			if !assert.NoError(t, err, `jwxtest.ParseJwkFile should succeed`) {
				return
			}

			t.Run("Parse JWK via jwx", func(t *testing.T) {
				// Better be a ECDSA private key
				var rawkey ecdsa.PrivateKey
				if !assert.NoError(t, jwxJwk.Raw(&rawkey), `jwk.Raw should succeed`) {
					return
				}
				_ = rawkey
			})
			t.Run("Encrypt with jose, Decrypt with jwx", func(t *testing.T) {
				// let jose encrypt payload using the key file
				joseCryptFile, joseCryptCleanup, err := jose.EncryptJwe(ctx, t, expected, joseJwkFile)
				if !assert.NoError(t, err, `jose.EncryptJwe should succeed`) {
					return
				}
				defer joseCryptCleanup()

				jwxtest.DumpFile(t, joseCryptFile)

				// let jwx decrypt the jose crypted file
				payload, err := jwxtest.DecryptJweFile(ctx, joseCryptFile, jwa.ECDH_ES, joseJwkFile)
				if !assert.NoError(t, err, `decryptFile.DecryptJwe should succeed`) {
					jwxtest.DumpFile(t, joseCryptFile)
					return
				}

				if !assert.Equal(t, expected, payload, `decrypted payloads should match`) {
					return
				}
			})
			t.Run("Encrypt with jwx, Decrypt with jose", func(t *testing.T) {
				jwxCryptFile, jwxCryptCleanup, err := jwxtest.EncryptJweFile(ctx, expected, jwa.ECDH_ES, joseJwkFile, jwa.A128GCM, jwa.NoCompress)
				if !assert.NoError(t, err, `jwxtest.EncryptJweFile should succeed`) {
					return
				}
				defer jwxCryptCleanup()

				payload, err := jose.DecryptJwe(ctx, t, jwxCryptFile, joseJwkFile)
				if !assert.NoError(t, err, `jose.DecryptJwe should succeed`) {
					return
				}

				if !assert.Equal(t, expected, payload, `decrypted payloads should match`) {
					return
				}
			})
			t.Run("Encrypt with jwx, Decrypt with jose", func(t *testing.T) {
				jwxCryptFile, jwxCryptCleanup, err := jwxtest.EncryptJweFile(ctx, expected, jwa.ECDH_ES, joseJwkFile, jwa.A128GCM, jwa.NoCompress)
				if !assert.NoError(t, err, `jwxtest.EncryptJweFile should succeed`) {
					return
				}
				defer jwxCryptCleanup()

				payload, err := jose.DecryptJwe(ctx, t, jwxCryptFile, joseJwkFile)
				if !assert.NoError(t, err, `jose.DecryptJwe should succeed`) {
					return
				}

				if !assert.Equal(t, expected, payload, `decrypted payloads should match`) {
					return
				}
			})
		})
	})
}
