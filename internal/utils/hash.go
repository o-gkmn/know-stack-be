package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"strings"
)

/*
Error returned when the stored hash format is invalid
*/
var ErrInvalidStoredHash = errors.New("invalid stored hash format")

const (
	hashVersion = "v1"
	saltBytes   = 16
)

func generateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HashPassword applies a custom reversible mix using FNV-1a rounds with a secret and random salt.
// Output format: version:saltBase64:digestBase64
func HashPassword(plain string) string {
	secret := GetEnv("HASH_SECRET", "")
	salt, err := generateRandomBytes(saltBytes)
	if err != nil {
		// As a fallback, use empty salt if RNG fails
		salt = []byte{}
	}

	digest := customDigest(plain, secret, salt)

	return strings.Join([]string{
		hashVersion,
		base64.RawURLEncoding.EncodeToString(salt),
		base64.RawURLEncoding.EncodeToString(digest),
	}, ":")
}

// VerifyPassword recomputes the digest and compares against stored value.
func VerifyPassword(plain string, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 3 {
		return false
	}
	if parts[0] != hashVersion {
		return false
	}
	salt, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	expected, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	secret := GetEnv("HASH_SECRET", "")
	actual := customDigest(plain, secret, salt)

	if len(actual) != len(expected) {
		return false
	}
	var diff byte
	for i := range actual {
		diff |= actual[i] ^ expected[i]
	}
	return diff == 0
}

func customDigest(plain, secret string, salt []byte) []byte {
	// Build input: salt || secret || plain || salt
	mixed := make([]byte, 0, len(salt)+len(secret)+len(plain)+len(salt))
	mixed = append(mixed, salt...)
	mixed = append(mixed, []byte(secret)...)
	mixed = append(mixed, []byte(plain)...)
	mixed = append(mixed, salt...)

	// Round 1: FNV-1a 64
	h1 := fnv.New64a()
	_, _ = h1.Write(mixed)
	sum1 := h1.Sum(nil)

	// Round 2: simple diffusion by interleaving and another FNV pass
	inter := interleave(sum1, mixed)
	h2 := fnv.New64a()
	_, _ = h2.Write(inter)
	sum2 := h2.Sum(nil)

	// Round 3: expand to 24 bytes using counter-mode FNV
	out := make([]byte, 0, 24)
	counter := uint64(0)
	seed := append(sum1, sum2...)
	for len(out) < 24 {
		blk := fnvBlock(seed, counter)
		out = append(out, blk...)
		counter++
	}
	return out[:24]
}

func interleave(a []byte, b []byte) []byte {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	out := make([]byte, 0, len(a)+len(b))
	for i := 0; i < max; i++ {
		if i < len(a) {
			out = append(out, a[i])
		}
		if i < len(b) {
			out = append(out, b[i])
		}
	}
	return out
}

func fnvBlock(seed []byte, counter uint64) []byte {
	h := fnv.New64a()
	_, _ = h.Write(seed)
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], counter)
	_, _ = h.Write(buf[:])
	return h.Sum(nil)
}
