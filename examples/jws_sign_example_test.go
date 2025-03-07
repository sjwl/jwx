package examples_test

import (
	"fmt"

	"github.com/sjwl/jwx/v2/jwa"
	"github.com/sjwl/jwx/v2/jwk"
	"github.com/sjwl/jwx/v2/jws"
)

func ExampleJWS_Sign() {
	key, err := jwk.FromRaw([]byte(`abracadabra`))
	if err != nil {
		fmt.Printf("failed to create key: %s\n", err)
		return
	}

	buf, err := jws.Sign([]byte("Lorem ipsum"), jws.WithKey(jwa.HS256, key))
	if err != nil {
		fmt.Printf("failed to sign payload: %s\n", err)
		return
	}
	fmt.Printf("%s\n", buf)
	// OUTPUT:
	// eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0.EjVtju0uXjSz6QevNgAqN1ESd9aNCP7-tJLifkQ0_C0
}
