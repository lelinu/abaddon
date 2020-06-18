package hash_utils

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lelinu/api_utils/utils/base64_utils"
	"hash/fnv"
	"io"
	"io/ioutil"
)

var Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func EncryptString(secret string, data string) (string, error) {
	d, err := compress([]byte(data))
	if err != nil {
		return "", err
	}
	d, err = encrypt([]byte(secret), d)
	if err != nil {
		return "", err
	}
	return base64_utils.EncodeFromBytes(d), nil
}

func DecryptString(secret string, data string) (string, error) {
	d, err := base64_utils.DecodeToBytes(data)
	if err != nil {
		fmt.Printf("Error is %v", err)
		return "", err
	}
	d, err = decrypt([]byte(secret), d)
	if err != nil {
		return "", err
	}
	d, err = decompress(d)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func Hash(str string, n int) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hashSize(hasher.Sum(nil), n)
}

func QuickHash(str string, n int) string {
	hasher := fnv.New64()
	hasher.Write([]byte(str))
	return hashSize(hasher.Sum(nil), n)
}

func HashStream(r io.Reader, n int) string {
	hasher := sha256.New()
	io.Copy(hasher, r)
	h := hex.EncodeToString(hasher.Sum(nil))
	if n == 0 {
		return h
	} else if n >= len(h) {
		return h
	}
	return h[0:n]
}

func hashSize(b []byte, n int) string {
	h := ""
	for i := 0; i < len(b); i++ {
		if n > 0 && len(h) >= n {
			break
		}
		h += ReversedBaseChange(Letters, int(b[i]))
	}

	if len(h) > n {
		return h[0 : len(h)-1]
	}
	return h
}

func ReversedBaseChange(alphabet []rune, i int) string {
	str := ""
	for {
		str += string(alphabet[i%len(alphabet)])
		i = i / len(alphabet)
		if i == 0 {
			break
		}
	}
	return str
}

func encrypt(key []byte, plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func compress(something []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(something)
	w.Close()
	return b.Bytes(), nil
}

func decompress(something []byte) ([]byte, error) {
	b := bytes.NewBuffer(something)
	r, err := zlib.NewReader(b)
	if err != nil {
		return []byte(""), nil
	}
	r.Close()
	return ioutil.ReadAll(r)
}

func sign(something []byte) ([]byte, error) {
	return something, nil
}

func verify(something []byte) ([]byte, error) {
	return something, nil
}
