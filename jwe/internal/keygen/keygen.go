package keygen

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"

	"github.com/sjwl/jwx/v2/internal/ecutil"
	"github.com/sjwl/jwx/v2/jwa"
	"github.com/sjwl/jwx/v2/jwe/internal/concatkdf"
	"github.com/sjwl/jwx/v2/jwk"
	"github.com/sjwl/jwx/v2/x25519"
)

// Bytes returns the byte from this ByteKey
func (k ByteKey) Bytes() []byte {
	return []byte(k)
}

// Size returns the size of the key
func (g Static) Size() int {
	return len(g)
}

// Generate returns the key
func (g Static) Generate() (ByteSource, error) {
	buf := make([]byte, g.Size())
	copy(buf, g)
	return ByteKey(buf), nil
}

// NewRandom creates a new Generator that returns
// random bytes
func NewRandom(n int) Random {
	return Random{keysize: n}
}

// Size returns the key size
func (g Random) Size() int {
	return g.keysize
}

// Generate generates a random new key
func (g Random) Generate() (ByteSource, error) {
	buf := make([]byte, g.keysize)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return nil, fmt.Errorf(`failed to read from rand.Reader: %w`, err)
	}
	return ByteKey(buf), nil
}

// NewEcdhes creates a new key generator using ECDH-ES
func NewEcdhes(alg jwa.KeyEncryptionAlgorithm, enc jwa.ContentEncryptionAlgorithm, keysize int, pubkey *ecdsa.PublicKey, apu, apv []byte) (*Ecdhes, error) {
	return &Ecdhes{
		algorithm: alg,
		enc:       enc,
		keysize:   keysize,
		pubkey:    pubkey,
		apu:       apu,
		apv:       apv,
	}, nil
}

// Size returns the key size associated with this generator
func (g Ecdhes) Size() int {
	return g.keysize
}

// Generate generates new keys using ECDH-ES
func (g Ecdhes) Generate() (ByteSource, error) {
	priv, err := ecdsa.GenerateKey(g.pubkey.Curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate key for ECDH-ES: %w`, err)
	}

	var algorithm string
	if g.algorithm == jwa.ECDH_ES {
		algorithm = g.enc.String()
	} else {
		algorithm = g.algorithm.String()
	}

	pubinfo := make([]byte, 4)
	binary.BigEndian.PutUint32(pubinfo, uint32(g.keysize)*8)

	if !priv.PublicKey.Curve.IsOnCurve(g.pubkey.X, g.pubkey.Y) {
		return nil, fmt.Errorf(`public key used does not contain a point (X,Y) on the curve`)
	}
	z, _ := priv.PublicKey.Curve.ScalarMult(g.pubkey.X, g.pubkey.Y, priv.D.Bytes())
	zBytes := ecutil.AllocECPointBuffer(z, priv.PublicKey.Curve)
	defer ecutil.ReleaseECPointBuffer(zBytes)
	kdf := concatkdf.New(crypto.SHA256, []byte(algorithm), zBytes, g.apu, g.apv, pubinfo, []byte{})
	kek := make([]byte, g.keysize)
	if _, err := kdf.Read(kek); err != nil {
		return nil, fmt.Errorf(`failed to read kdf: %w`, err)
	}

	return ByteWithECPublicKey{
		PublicKey: &priv.PublicKey,
		ByteKey:   ByteKey(kek),
	}, nil
}

// NewX25519 creates a new key generator using ECDH-ES
func NewX25519(alg jwa.KeyEncryptionAlgorithm, enc jwa.ContentEncryptionAlgorithm, keysize int, pubkey x25519.PublicKey) (*X25519, error) {
	return &X25519{
		algorithm: alg,
		enc:       enc,
		keysize:   keysize,
		pubkey:    pubkey,
	}, nil
}

// Size returns the key size associated with this generator
func (g X25519) Size() int {
	return g.keysize
}

// Generate generates new keys using ECDH-ES
func (g X25519) Generate() (ByteSource, error) {
	pub, priv, err := x25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate key for X25519: %w`, err)
	}

	var algorithm string
	if g.algorithm == jwa.ECDH_ES {
		algorithm = g.enc.String()
	} else {
		algorithm = g.algorithm.String()
	}

	pubinfo := make([]byte, 4)
	binary.BigEndian.PutUint32(pubinfo, uint32(g.keysize)*8)

	zBytes, err := curve25519.X25519(priv.Seed(), g.pubkey)
	if err != nil {
		return nil, fmt.Errorf(`failed to compute Z: %w`, err)
	}
	kdf := concatkdf.New(crypto.SHA256, []byte(algorithm), zBytes, []byte{}, []byte{}, pubinfo, []byte{})
	kek := make([]byte, g.keysize)
	if _, err := kdf.Read(kek); err != nil {
		return nil, fmt.Errorf(`failed to read kdf: %w`, err)
	}

	return ByteWithECPublicKey{
		PublicKey: pub,
		ByteKey:   ByteKey(kek),
	}, nil
}

// HeaderPopulate populates the header with the required EC-DSA public key
// information ('epk' key)
func (k ByteWithECPublicKey) Populate(h Setter) error {
	key, err := jwk.FromRaw(k.PublicKey)
	if err != nil {
		return fmt.Errorf(`failed to create JWK: %w`, err)
	}

	if err := h.Set("epk", key); err != nil {
		return fmt.Errorf(`failed to write header: %w`, err)
	}
	return nil
}

// HeaderPopulate populates the header with the required AES GCM
// parameters ('iv' and 'tag')
func (k ByteWithIVAndTag) Populate(h Setter) error {
	if err := h.Set("iv", k.IV); err != nil {
		return fmt.Errorf(`failed to write header: %w`, err)
	}

	if err := h.Set("tag", k.Tag); err != nil {
		return fmt.Errorf(`failed to write header: %w`, err)
	}

	return nil
}

// HeaderPopulate populates the header with the required PBES2
// parameters ('p2s' and 'p2c')
func (k ByteWithSaltAndCount) Populate(h Setter) error {
	if err := h.Set("p2c", k.Count); err != nil {
		return fmt.Errorf(`failed to write header: %w`, err)
	}

	if err := h.Set("p2s", k.Salt); err != nil {
		return fmt.Errorf(`failed to write header: %w`, err)
	}

	return nil
}
