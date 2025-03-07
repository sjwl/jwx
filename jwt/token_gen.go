// Code generated by tools/cmd/genjwt/main.go. DO NOT EDIT.

package jwt

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/lestrrat-go/iter/mapiter"
	"github.com/sjwl/jwx/v2/internal/base64"
	"github.com/sjwl/jwx/v2/internal/iter"
	"github.com/sjwl/jwx/v2/internal/json"
	"github.com/sjwl/jwx/v2/internal/pool"
	"github.com/sjwl/jwx/v2/jwt/internal/types"
)

const (
	AudienceKey   = "aud"
	ExpirationKey = "exp"
	IssuedAtKey   = "iat"
	IssuerKey     = "iss"
	JwtIDKey      = "jti"
	NotBeforeKey  = "nbf"
	SubjectKey    = "sub"
)

// Token represents a generic JWT token.
// which are type-aware (to an extent). Other claims may be accessed via the `Get`/`Set`
// methods but their types are not taken into consideration at all. If you have non-standard
// claims that you must frequently access, consider creating accessors functions
// like the following
//
// func SetFoo(tok jwt.Token) error
// func GetFoo(tok jwt.Token) (*Customtyp, error)
//
// Embedding jwt.Token into another struct is not recommended, because
// jwt.Token needs to handle private claims, and this really does not
// work well when it is embedded in other structure
type Token interface {

	// Audience returns the value for "aud" field of the token
	Audience() []string

	// Expiration returns the value for "exp" field of the token
	Expiration() time.Time

	// IssuedAt returns the value for "iat" field of the token
	IssuedAt() time.Time

	// Issuer returns the value for "iss" field of the token
	Issuer() string

	// JwtID returns the value for "jti" field of the token
	JwtID() string

	// NotBefore returns the value for "nbf" field of the token
	NotBefore() time.Time

	// Subject returns the value for "sub" field of the token
	Subject() string

	// PrivateClaims return the entire set of fields (claims) in the token
	// *other* than the pre-defined fields such as `iss`, `nbf`, `iat`, etc.
	PrivateClaims() map[string]interface{}

	// Get returns the value of the corresponding field in the token, such as
	// `nbf`, `exp`, `iat`, and other user-defined fields. If the field does not
	// exist in the token, the second return value will be `false`
	//
	// If you need to access fields like `alg`, `kid`, `jku`, etc, you need
	// to access the corresponding fields in the JWS/JWE message. For this,
	// you will need to access them by directly parsing the payload using
	// `jws.Parse` and `jwe.Parse`
	Get(string) (interface{}, bool)

	// Set assigns a value to the corresponding field in the token. Some
	// pre-defined fields such as `nbf`, `iat`, `iss` need their values to
	// be of a specific type. See the other getter methods in this interface
	// for the types of each of these fields
	Set(string, interface{}) error
	Remove(string) error

	// Options returns the per-token options associated with this token.
	// The options set value will be copied when the token is cloned via `Clone()`
	// but it will not survive when the token goes through marshaling/unmarshaling
	// such as `json.Marshal` and `json.Unmarshal`
	Options() *TokenOptionSet
	Clone() (Token, error)
	Iterate(context.Context) Iterator
	Walk(context.Context, Visitor) error
	AsMap(context.Context) (map[string]interface{}, error)
}
type stdToken struct {
	mu            *sync.RWMutex
	dc            DecodeCtx          // per-object context for decoding
	options       TokenOptionSet     // per-object option
	audience      types.StringList   // https://tools.ietf.org/html/rfc7519#section-4.1.3
	expiration    *types.NumericDate // https://tools.ietf.org/html/rfc7519#section-4.1.4
	issuedAt      *types.NumericDate // https://tools.ietf.org/html/rfc7519#section-4.1.6
	issuer        *string            // https://tools.ietf.org/html/rfc7519#section-4.1.1
	jwtID         *string            // https://tools.ietf.org/html/rfc7519#section-4.1.7
	notBefore     *types.NumericDate // https://tools.ietf.org/html/rfc7519#section-4.1.5
	subject       *string            // https://tools.ietf.org/html/rfc7519#section-4.1.2
	privateClaims map[string]interface{}
}

// New creates a standard token, with minimal knowledge of
// possible claims. Standard claims include"aud", "exp", "iat", "iss", "jti", "nbf" and "sub".
// Convenience accessors are provided for these standard claims
func New() Token {
	return &stdToken{
		mu:            &sync.RWMutex{},
		privateClaims: make(map[string]interface{}),
		options:       DefaultOptionSet(),
	}
}

func (t *stdToken) Options() *TokenOptionSet {
	return &t.options
}

func (t *stdToken) Get(name string) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	switch name {
	case AudienceKey:
		if t.audience == nil {
			return nil, false
		}
		v := t.audience.Get()
		return v, true
	case ExpirationKey:
		if t.expiration == nil {
			return nil, false
		}
		v := t.expiration.Get()
		return v, true
	case IssuedAtKey:
		if t.issuedAt == nil {
			return nil, false
		}
		v := t.issuedAt.Get()
		return v, true
	case IssuerKey:
		if t.issuer == nil {
			return nil, false
		}
		v := *(t.issuer)
		return v, true
	case JwtIDKey:
		if t.jwtID == nil {
			return nil, false
		}
		v := *(t.jwtID)
		return v, true
	case NotBeforeKey:
		if t.notBefore == nil {
			return nil, false
		}
		v := t.notBefore.Get()
		return v, true
	case SubjectKey:
		if t.subject == nil {
			return nil, false
		}
		v := *(t.subject)
		return v, true
	default:
		v, ok := t.privateClaims[name]
		return v, ok
	}
}

func (t *stdToken) Remove(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	switch key {
	case AudienceKey:
		t.audience = nil
	case ExpirationKey:
		t.expiration = nil
	case IssuedAtKey:
		t.issuedAt = nil
	case IssuerKey:
		t.issuer = nil
	case JwtIDKey:
		t.jwtID = nil
	case NotBeforeKey:
		t.notBefore = nil
	case SubjectKey:
		t.subject = nil
	default:
		delete(t.privateClaims, key)
	}
	return nil
}

func (t *stdToken) Set(name string, value interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.setNoLock(name, value)
}

func (t *stdToken) DecodeCtx() DecodeCtx {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.dc
}

func (t *stdToken) SetDecodeCtx(v DecodeCtx) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.dc = v
}

func (t *stdToken) setNoLock(name string, value interface{}) error {
	switch name {
	case AudienceKey:
		var acceptor types.StringList
		if err := acceptor.Accept(value); err != nil {
			return fmt.Errorf(`invalid value for %s key: %w`, AudienceKey, err)
		}
		t.audience = acceptor
		return nil
	case ExpirationKey:
		var acceptor types.NumericDate
		if err := acceptor.Accept(value); err != nil {
			return fmt.Errorf(`invalid value for %s key: %w`, ExpirationKey, err)
		}
		t.expiration = &acceptor
		return nil
	case IssuedAtKey:
		var acceptor types.NumericDate
		if err := acceptor.Accept(value); err != nil {
			return fmt.Errorf(`invalid value for %s key: %w`, IssuedAtKey, err)
		}
		t.issuedAt = &acceptor
		return nil
	case IssuerKey:
		if v, ok := value.(string); ok {
			t.issuer = &v
			return nil
		}
		return fmt.Errorf(`invalid value for %s key: %T`, IssuerKey, value)
	case JwtIDKey:
		if v, ok := value.(string); ok {
			t.jwtID = &v
			return nil
		}
		return fmt.Errorf(`invalid value for %s key: %T`, JwtIDKey, value)
	case NotBeforeKey:
		var acceptor types.NumericDate
		if err := acceptor.Accept(value); err != nil {
			return fmt.Errorf(`invalid value for %s key: %w`, NotBeforeKey, err)
		}
		t.notBefore = &acceptor
		return nil
	case SubjectKey:
		if v, ok := value.(string); ok {
			t.subject = &v
			return nil
		}
		return fmt.Errorf(`invalid value for %s key: %T`, SubjectKey, value)
	default:
		if t.privateClaims == nil {
			t.privateClaims = map[string]interface{}{}
		}
		t.privateClaims[name] = value
	}
	return nil
}

func (t *stdToken) Audience() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.audience != nil {
		return t.audience.Get()
	}
	return nil
}

func (t *stdToken) Expiration() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.expiration != nil {
		return t.expiration.Get()
	}
	return time.Time{}
}

func (t *stdToken) IssuedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.issuedAt != nil {
		return t.issuedAt.Get()
	}
	return time.Time{}
}

func (t *stdToken) Issuer() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.issuer != nil {
		return *(t.issuer)
	}
	return ""
}

func (t *stdToken) JwtID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.jwtID != nil {
		return *(t.jwtID)
	}
	return ""
}

func (t *stdToken) NotBefore() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.notBefore != nil {
		return t.notBefore.Get()
	}
	return time.Time{}
}

func (t *stdToken) Subject() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.subject != nil {
		return *(t.subject)
	}
	return ""
}

func (t *stdToken) PrivateClaims() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.privateClaims
}

func (t *stdToken) makePairs() []*ClaimPair {
	t.mu.RLock()
	defer t.mu.RUnlock()

	pairs := make([]*ClaimPair, 0, 7)
	if t.audience != nil {
		v := t.audience.Get()
		pairs = append(pairs, &ClaimPair{Key: AudienceKey, Value: v})
	}
	if t.expiration != nil {
		v := t.expiration.Get()
		pairs = append(pairs, &ClaimPair{Key: ExpirationKey, Value: v})
	}
	if t.issuedAt != nil {
		v := t.issuedAt.Get()
		pairs = append(pairs, &ClaimPair{Key: IssuedAtKey, Value: v})
	}
	if t.issuer != nil {
		v := *(t.issuer)
		pairs = append(pairs, &ClaimPair{Key: IssuerKey, Value: v})
	}
	if t.jwtID != nil {
		v := *(t.jwtID)
		pairs = append(pairs, &ClaimPair{Key: JwtIDKey, Value: v})
	}
	if t.notBefore != nil {
		v := t.notBefore.Get()
		pairs = append(pairs, &ClaimPair{Key: NotBeforeKey, Value: v})
	}
	if t.subject != nil {
		v := *(t.subject)
		pairs = append(pairs, &ClaimPair{Key: SubjectKey, Value: v})
	}
	for k, v := range t.privateClaims {
		pairs = append(pairs, &ClaimPair{Key: k, Value: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key.(string) < pairs[j].Key.(string)
	})
	return pairs
}

func (t *stdToken) UnmarshalJSON(buf []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.audience = nil
	t.expiration = nil
	t.issuedAt = nil
	t.issuer = nil
	t.jwtID = nil
	t.notBefore = nil
	t.subject = nil
	dec := json.NewDecoder(bytes.NewReader(buf))
LOOP:
	for {
		tok, err := dec.Token()
		if err != nil {
			return fmt.Errorf(`error reading token: %w`, err)
		}
		switch tok := tok.(type) {
		case json.Delim:
			// Assuming we're doing everything correctly, we should ONLY
			// get either '{' or '}' here.
			if tok == '}' { // End of object
				break LOOP
			} else if tok != '{' {
				return fmt.Errorf(`expected '{', but got '%c'`, tok)
			}
		case string: // Objects can only have string keys
			switch tok {
			case AudienceKey:
				var decoded types.StringList
				if err := dec.Decode(&decoded); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, AudienceKey, err)
				}
				t.audience = decoded
			case ExpirationKey:
				var decoded types.NumericDate
				if err := dec.Decode(&decoded); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, ExpirationKey, err)
				}
				t.expiration = &decoded
			case IssuedAtKey:
				var decoded types.NumericDate
				if err := dec.Decode(&decoded); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, IssuedAtKey, err)
				}
				t.issuedAt = &decoded
			case IssuerKey:
				if err := json.AssignNextStringToken(&t.issuer, dec); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, IssuerKey, err)
				}
			case JwtIDKey:
				if err := json.AssignNextStringToken(&t.jwtID, dec); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, JwtIDKey, err)
				}
			case NotBeforeKey:
				var decoded types.NumericDate
				if err := dec.Decode(&decoded); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, NotBeforeKey, err)
				}
				t.notBefore = &decoded
			case SubjectKey:
				if err := json.AssignNextStringToken(&t.subject, dec); err != nil {
					return fmt.Errorf(`failed to decode value for key %s: %w`, SubjectKey, err)
				}
			default:
				if dc := t.dc; dc != nil {
					if localReg := dc.Registry(); localReg != nil {
						decoded, err := localReg.Decode(dec, tok)
						if err == nil {
							t.setNoLock(tok, decoded)
							continue
						}
					}
				}
				decoded, err := registry.Decode(dec, tok)
				if err == nil {
					t.setNoLock(tok, decoded)
					continue
				}
				return fmt.Errorf(`could not decode field %s: %w`, tok, err)
			}
		default:
			return fmt.Errorf(`invalid token %T`, tok)
		}
	}
	return nil
}

func (t stdToken) MarshalJSON() ([]byte, error) {
	buf := pool.GetBytesBuffer()
	defer pool.ReleaseBytesBuffer(buf)
	buf.WriteByte('{')
	enc := json.NewEncoder(buf)
	for i, pair := range t.makePairs() {
		f := pair.Key.(string)
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteRune('"')
		buf.WriteString(f)
		buf.WriteString(`":`)
		switch f {
		case AudienceKey:
			if err := json.EncodeAudience(enc, pair.Value.([]string), t.options.IsEnabled(FlattenAudience)); err != nil {
				return nil, fmt.Errorf(`failed to encode "aud": %w`, err)
			}
			continue
		case ExpirationKey, IssuedAtKey, NotBeforeKey:
			enc.Encode(pair.Value.(time.Time).Unix())
			continue
		}
		switch v := pair.Value.(type) {
		case []byte:
			buf.WriteRune('"')
			buf.WriteString(base64.EncodeToString(v))
			buf.WriteRune('"')
		default:
			if err := enc.Encode(v); err != nil {
				return nil, fmt.Errorf(`failed to marshal field %s: %w`, f, err)
			}
			buf.Truncate(buf.Len() - 1)
		}
	}
	buf.WriteByte('}')
	ret := make([]byte, buf.Len())
	copy(ret, buf.Bytes())
	return ret, nil
}

func (t *stdToken) Iterate(ctx context.Context) Iterator {
	pairs := t.makePairs()
	ch := make(chan *ClaimPair, len(pairs))
	go func(ctx context.Context, ch chan *ClaimPair, pairs []*ClaimPair) {
		defer close(ch)
		for _, pair := range pairs {
			select {
			case <-ctx.Done():
				return
			case ch <- pair:
			}
		}
	}(ctx, ch, pairs)
	return mapiter.New(ch)
}

func (t *stdToken) Walk(ctx context.Context, visitor Visitor) error {
	return iter.WalkMap(ctx, t, visitor)
}

func (t *stdToken) AsMap(ctx context.Context) (map[string]interface{}, error) {
	return iter.AsMap(ctx, t)
}
