package transformers

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// ToUpperCase: Convert string to uppercase
func ToUpperCase(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return strings.ToUpper(str)
	}
	return input
}

// ToLowerCase: Convert string to lowercase
func ToLowerCase(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return strings.ToLower(str)
	}
	return input
}

// Trim: Trim spaces from string
func Trim(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return strings.TrimSpace(str)
	}
	return input
}

// TrimLeft: Trim spaces from left of string
func TrimLeft(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return strings.TrimLeft(str, " ")
	}
	return input
}

// TrimRight: Trim spaces from right of string
func TrimRight(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return strings.TrimRight(str, " ")
	}
	return input
}

// Base64Encode: Encode string to base64
func Base64Encode(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return input
}

// Base64Decode: Decode base64 string
func Base64Decode(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(str)
		if err == nil {
			return string(decoded)
		}
	}
	return input
}

// UrlEncode: Encode string to URL
func UrlEncode(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		return url.QueryEscape(str)
	}
	return input
}

// UrlDecode: Decode URL string
func UrlDecode(input interface{}) interface{} {
	if str, ok := input.(string); ok {
		decoded, err := url.QueryUnescape(str)
		if err == nil {
			return decoded
		}
	}
	return input
}
