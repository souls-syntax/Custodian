package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
)

func NormalizeQuery(s string) string {

	s = strings.TrimSpace(s)

	s = strings.ToLower(s)

	fields := strings.FieldsFunc(s, unicode.IsSpace)
	return strings.Join(fields, " ")
}

func HashQuery(normalized string) string {

	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])

}
