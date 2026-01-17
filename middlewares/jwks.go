package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type JWKS struct {
	mu   sync.RWMutex
	keys map[string]*rsa.PublicKey
}

func jwkCache(jwksURL string) *JWKS {
	j := &JWKS{keys: map[string]*rsa.PublicKey{}}
	j.refresh(jwksURL)

	go func() {
		for range time.Tick(5 * time.Minute) {
			j.refresh(jwksURL)
		}
	}()

	return j
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (k jwk) PublicKey() (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}

	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

func (j *JWKS) refresh(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data struct {
		Keys []jwk `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	j.keys = map[string]*rsa.PublicKey{} // reset cache

	for _, k := range data.Keys {
		if k.Use != "sig" {
			continue
		}
		if k.Alg != "RS256" {
			continue
		}

		pub, err := k.PublicKey()
		if err != nil {
			continue
		}

		j.keys[k.Kid] = pub
	}
}

func (j *JWKS) GetKey(kid string) (*rsa.PublicKey, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	key, ok := j.keys[kid]
	if !ok {
		return nil, fmt.Errorf("unknown kid")
	}

	return key, nil
}
