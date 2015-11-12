package jwe

// NewEncrypt creates a new Encrypt struct. The caller is responsible
// for instantiating valid inputs for ContentEncrypter, KeyGenerator,
// and KeyEncrypters.
func NewEncrypt(cc ContentEncrypter, kg KeyGenerator, ke ...KeyEncrypter) *Encrypt {
	e := &Encrypt{
		ContentEncrypter: cc,
		KeyGenerator:     kg,
		KeyEncrypters:    ke,
	}
	return e
}

// Encrypt takes the plaintext and encrypts into a JWE message.
func (e Encrypt) Encrypt(plaintext []byte) (*Message, error) {
	cek, err := e.KeyGenerator.KeyGenerate()
	if err != nil {
		return nil, err
	}
	debug("Encrypt: generated cek len = %d", len(cek))

	protected := NewEncodedHeader()
	protected.Set("enc", e.ContentEncrypter.Algorithm())

	// In JWE, multiple recipients may exist -- they receive an
	// encrypted version of the CEK, using their key encryption
	// algorithm of choice.
	recipients := make([]Recipient, len(e.KeyEncrypters))
	for i, enc := range e.KeyEncrypters {
		r := NewRecipient()
		r.Header.Set("alg", enc.Algorithm())
		if v := enc.Kid(); v != "" {
			r.Header.Set("kid", v)
		}
		enckey, err := enc.KeyEncrypt(cek)
		if err != nil {
			return nil, err
		}
		r.EncryptedKey = enckey
		debug("Encrypt: encrypted_key = %x", enckey)
		recipients[i] = *r
	}

	// If there's only one recipient, you want to include that in the
	// protected header
	if len(recipients) == 1 {
		protected.Set("alg", recipients[0].Header.Algorithm)
	}

	aad, err := protected.Base64Encode()
	if err != nil {
		return nil, err
	}

	// ...on the other hand, there's only one content cipher.
	iv, ciphertext, tag, err := e.ContentEncrypter.Encrypt(cek, plaintext, aad)

	debug("Encrypt.Encrypt: cek        = %x (%d)", cek, len(cek))
	debug("Encrypt.Encrypt: aad        = %x", aad)
	debug("Encrypt.Encrypt: ciphertext = %x", ciphertext)
	debug("Encrypt.Encrypt: iv         = %x", iv)
	debug("Encrypt.Encrypt: tag        = %x", tag)

	msg := NewMessage()
	msg.AuthenticatedData.Base64Decode(aad)
	msg.CipherText = ciphertext
	msg.InitializationVector = iv
	msg.ProtectedHeader = protected
	msg.Recipients = recipients
	msg.Tag = tag

	return msg, nil
}
