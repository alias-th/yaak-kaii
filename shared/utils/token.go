package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

// RandomByte สร้าง n ไบต์จาก crypto/rand
func RandomByte(nBytes int) ([]byte, error) {
	if nBytes <= 0 {
		return nil, errors.New("nBytes must be > 0")
	}
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// EncodeBase64URL แปลงไบต์เป็น URL-safe base64 (ไม่มี padding)
func EncodeBase64URL(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeBase64URL ถอดรหัส base64 (RawURLEncoding) กลับเป็นไบต์
func DecodeBase64URL(token string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(token)
}

// HashTokenBytes คืนค่า hash (sha256) ในรูป base64-url (ไม่ใส่ padding)
func HashTokenBytes(raw []byte) string {
	h := sha256.Sum256(raw)
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// VerifyToken ตรวจสอบ token (string base64-url) เทียบกับ hash ที่เก็บไว้ (base64-url)
func VerifyToken(token string, storedHash string) (bool, error) {
	decoded, err := DecodeBase64URL(token)
	if err != nil {
		return false, err
	}
	gotHash := HashTokenBytes(decoded)

	// ใช้ constant time compare เพื่อความปลอดภัยจาก timing attack
	if subtle.ConstantTimeCompare([]byte(gotHash), []byte(storedHash)) == 1 {
		return true, nil
	}
	return false, nil
}

// Convenience: สร้าง token พร้อม hash เพื่อเก็บใน DB (คืน token ที่จะส่งให้ client และ hash สำหรับเก็บ)
func GenerateTokenPair(nBytes int) (token string, tokenHash string, err error) {
	raw, err := RandomByte(nBytes)
	if err != nil {
		return "", "", err
	}
	token = EncodeBase64URL(raw)
	tokenHash = HashTokenBytes(raw)

	return token, tokenHash, nil
}
