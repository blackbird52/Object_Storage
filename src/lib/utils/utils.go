package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
)

func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 || digest[:8] != "SHA-256=" {
		return ""
	} else {
		return digest[8:]
	}
}

func GetSizeFromHeader(h http.Header) int64 {
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	return size
}

func CalculateHash(r io.Reader) string {
	hash := sha256.New()
	_, _ = io.Copy(hash, r)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
