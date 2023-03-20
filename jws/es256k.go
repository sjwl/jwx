//go:build jwx_es256k
// +build jwx_es256k

package jws

import (
	"github.com/sjwl/jwx/v2/jwa"
)

func init() {
	addAlgorithmForKeyType(jwa.EC, jwa.ES256K)
}
