package examples_test

import (
	"fmt"

	"github.com/sjwl/jwx/v2/jwa"
	"github.com/sjwl/jwx/v2/jwk"
	"github.com/sjwl/jwx/v2/jwt"
)

func ExampleJWT_ParseWithKey() {
	const keysrc = `{"kty":"oct","k":"AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow"}`

	key, err := jwk.ParseKey([]byte(keysrc))
	if err != nil {
		fmt.Printf("jwk.ParseKey failed: %s\n", err)
		return
	}

	tok, err := jwt.Parse([]byte(exampleJWTSignedHMAC), jwt.WithKey(jwa.HS256, key), jwt.WithValidate(false))
	if err != nil {
		fmt.Printf("jwt.Parse failed: %s\n", err)
		return
	}
	_ = tok
	// OUTPUT:
}
