package examples_test

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sjwl/jwx/v2/jwt"
)

func ExampleJWT_SerializeJSON() {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/sjwl/jwx`).
		IssuedAt(time.Unix(aLongLongTimeAgo, 0)).
		Build()
	if err != nil {
		fmt.Printf("failed to build token: %s\n", err)
		return
	}

	json.NewEncoder(os.Stdout).Encode(tok)
	// OUTPUT:
	// {"iat":233431200,"iss":"github.com/sjwl/jwx"}
}
