package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	msg := []byte("the quick brown fox jumps over the lazy dog")

	sha256Digest(msg)
	hmacSHA256(msg)
	aesGCMRoundTrip(msg)
	rsaPSSSignVerify(msg)
	ecdsaP256SignVerify(msg)

	fmt.Println("all FIPS-approved primitives exercised OK")
}

func sha256Digest(msg []byte) {
	sum := sha256.Sum256(msg)
	fmt.Printf("SHA-256:        %s\n", hex.EncodeToString(sum[:]))
}

func hmacSHA256(msg []byte) {
	key := mustRand(32)
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	fmt.Printf("HMAC-SHA256:    %s\n", hex.EncodeToString(mac.Sum(nil)))
}

func aesGCMRoundTrip(msg []byte) {
	key := mustRand(32) // AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes.NewCipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher.NewGCM: %v", err)
	}
	nonce := mustRand(gcm.NonceSize())
	ct := gcm.Seal(nil, nonce, msg, nil)
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		log.Fatalf("gcm.Open: %v", err)
	}
	if string(pt) != string(msg) {
		log.Fatalf("AES-GCM round-trip mismatch")
	}
	fmt.Printf("AES-256-GCM:    ct=%s (round-trip OK)\n", hex.EncodeToString(ct))
}

func rsaPSSSignVerify(msg []byte) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	digest := sha256.Sum256(msg)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA256}
	sig, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, digest[:], opts)
	if err != nil {
		log.Fatalf("rsa.SignPSS: %v", err)
	}
	if err := rsa.VerifyPSS(&key.PublicKey, crypto.SHA256, digest[:], sig, opts); err != nil {
		log.Fatalf("rsa.VerifyPSS: %v", err)
	}
	fmt.Printf("RSA-PSS-2048:   sig=%d bytes (verified)\n", len(sig))
}

func ecdsaP256SignVerify(msg []byte) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("ecdsa.GenerateKey: %v", err)
	}
	digest := sha256.Sum256(msg)
	sig, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		log.Fatalf("ecdsa.SignASN1: %v", err)
	}
	if !ecdsa.VerifyASN1(&key.PublicKey, digest[:], sig) {
		log.Fatalf("ecdsa.VerifyASN1: signature did not verify")
	}
	fmt.Printf("ECDSA-P256:     sig=%d bytes (verified)\n", len(sig))
}

func mustRand(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("rand.Read: %v", err)
	}
	return b
}
