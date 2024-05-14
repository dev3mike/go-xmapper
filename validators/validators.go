package validators

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RequiredValidator checks if the input string is not empty
func RequiredValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    if strings.TrimSpace(str) == "" {
        return fmt.Errorf("input is required and cannot be empty")
    }
    return nil
}

// EmailValidator checks if the input string is a valid email address
func EmailValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`)
    if !emailRegex.MatchString(str) {
        return fmt.Errorf("input is not a valid email address")
    }
    return nil
}

// PhoneValidator checks if the input string is a valid international phone number
func PhoneValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
    if !phoneRegex.MatchString(str) {
        return fmt.Errorf("input is not a valid international phone number")
    }
    return nil
}

// StrongPasswordValidator ensures the password meets strength requirements
func StrongPasswordValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    // Validate minimum length
    if len(str) < 8 {
        return fmt.Errorf("password must be at least 8 characters long")
    }

    // Check for at least one uppercase letter
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(str)
    if !hasUpper {
        return fmt.Errorf("password must contain at least one uppercase letter")
    }

    // Check for at least one lowercase letter
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(str)
    if !hasLower {
        return fmt.Errorf("password must contain at least one lowercase letter")
    }

    // Check for at least one digit
    hasDigit := regexp.MustCompile(`\d`).MatchString(str)
    if !hasDigit {
        return fmt.Errorf("password must contain at least one digit")
    }

    // Check for at least one special character
    hasSpecial := regexp.MustCompile(`[\^$*.\[\]{}()?!"@#%&/,><':;|_~` + "`" + `"-]`).MatchString(str)
    if !hasSpecial {
        return fmt.Errorf("password must contain at least one special character")
    }

    return nil
}


// DateValidator checks if the input string is a valid date in YYYY-MM-DD format
func DateValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    _, err := time.Parse("2006-01-02", str)
    if err != nil {
        return fmt.Errorf("input is not a valid date, expected format YYYY-MM-DD")
    }
    return nil
}


// TimeValidator checks if the input string is a valid time in HH:MM:SS format
func TimeValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    _, err := time.Parse("15:04:05", str)
    if err != nil {
        return fmt.Errorf("input is not a valid time, expected format HH:MM:SS")
    }
    return nil
}

// DatetimeValidator checks if the input string is a valid datetime in YYYY-MM-DD HH:MM:SS with timezone format
func DatetimeValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    _, err := time.Parse("2006-01-02 15:04:05 MST", str)
    if err != nil {
        return fmt.Errorf("input is not a valid datetime with timezone, expected format YYYY-MM-DD HH:MM:SS MST")
    }
    return nil
}

// UrlValidator checks if the input string is a valid URL
func UrlValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    u, err := url.ParseRequestURI(str)
    if err != nil || u.Scheme == "" || u.Host == "" {
        return fmt.Errorf("input is not a valid URL")
    }
    return nil
}

// IpValidator checks if the input string is a valid IP address
func IpValidator(input interface{}, _ string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    if net.ParseIP(str) == nil {
        return fmt.Errorf("input is not a valid IP address")
    }
    return nil
}

// MinLengthValidator checks if the input string's length is at least the specified minimum
func MinLengthValidator(input interface{}, length string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    minLength, err := strconv.Atoi(length)
    if err != nil {
        return fmt.Errorf("failed to convert length to integer")
    }
    if len(str) < minLength {
        return fmt.Errorf("input does not meet the minimum length requirement, minimum length is %s", length)
    }
    return nil
}

// MaxLengthValidator checks if the input string's length does not exceed the specified maximum
func MaxLengthValidator(input interface{}, length string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    maxLength, err := strconv.Atoi(length)
    if err != nil {
        return fmt.Errorf("failed to convert length to integer")
    }
    if len(str) > maxLength {
        return fmt.Errorf("input exceeds the maximum length requirement, maximum length is %s", length)
    }
    return nil
}

func GreaterThanValidator(input interface{}, threshold string) error {
    number, thresh, err := convertToFloat64(input, threshold)
    if err != nil {
        return err
    }
    if number <= thresh {
        return fmt.Errorf("input must be greater than %s", threshold)
    }
    return nil
}

func LessThanValidator(input interface{}, threshold string) error {
    number, thresh, err := convertToFloat64(input, threshold)
    if err != nil {
        return err
    }
    if number >= thresh {
        return fmt.Errorf("input must be less than %s", threshold)
    }
    return nil
}

func GreaterThanOrEqualValidator(input interface{}, threshold string) error {
    number, thresh, err := convertToFloat64(input, threshold)
    if err != nil {
        return err
    }
    if number < thresh {
        return fmt.Errorf("input must be greater than or equal to %s", threshold)
    }
    return nil
}

func LessThanOrEqualValidator(input interface{}, threshold string) error {
    number, thresh, err := convertToFloat64(input, threshold)
    if err != nil {
        return err
    }
    if number > thresh {
        return fmt.Errorf("input must be less than or equal to %s", threshold)
    }
    return nil
}

// EnumValidator checks if the input string is one of the allowed values
func EnumValidator(input interface{}, allowedValues string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    values := strings.Split(allowedValues, "-")
    for _, value := range values {
        if strings.TrimSpace(str) == strings.TrimSpace(value) {
            return nil
        }
    }
    return fmt.Errorf("input must be one of the following values: %s", allowedValues)
}

// BooleanValidator checks if the input is a boolean
func BooleanValidator(input interface{}, _ string) error {
    _, ok := input.(bool)
    if !ok {
        return fmt.Errorf("input must be a boolean")
    }
    return nil
}

// ContainsValidator checks if the input string contains one of the allowed values
func ContainsValidator(input interface{}, allowedValues string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    values := strings.Split(allowedValues, ",")
    for _, value := range values {
        if strings.Contains(str, strings.TrimSpace(value)) {
            return nil
        }
    }
    return fmt.Errorf("input must contain one of the following values: %s", allowedValues)
}

// NotContainsValidator checks if the input string does not contain any of the disallowed values
func NotContainsValidator(input interface{}, disallowedValues string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("failed to map the input to a string")
    }
    values := strings.Split(disallowedValues, ",")
    for _, value := range values {
        if strings.Contains(str, strings.TrimSpace(value)) {
            return fmt.Errorf("input must not contain the following value: %s", value)
        }
    }
    return nil
}

// RangeValidator checks if a number is between two numbers specified with a dash (e.g., "10-100")
func RangeValidator(input interface{}, rangeStr string) error {
    val := reflect.ValueOf(input)
    if val.Kind() != reflect.Int && val.Kind() != reflect.Float64 && val.Kind() != reflect.Float32 {
        return fmt.Errorf("input must be a number")
    }

    parts := strings.Split(rangeStr, "-")
    if len(parts) != 2 {
        return fmt.Errorf("range format is incorrect, must be 'min-max'")
    }

    min, err := strconv.ParseFloat(parts[0], 64)
    if err != nil {
        return fmt.Errorf("failed to parse minimum value: %s", parts[0])
    }

    max, err := strconv.ParseFloat(parts[1], 64)
    if err != nil {
        return fmt.Errorf("failed to parse maximum value: %s", parts[1])
    }

    if min > max {
        return fmt.Errorf("minimum value must be less than maximum value")
    }

    number, _, err := convertToFloat64(input, parts[0])
    if err != nil {
        return err
    }

    if number < min || number > max {
        return fmt.Errorf("input must be between %s and %s", parts[0], parts[1])
    }
    return nil
}

// StartsWidthValidator validates that a string starts with a specified substring
func StartsWidthValidator(input interface{}, prefix string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("input must be a string")
    }
    if !strings.HasPrefix(str, prefix) {
        return fmt.Errorf("input must start with '%s'", prefix)
    }
    return nil
}

// EndsWithValidator validates that a string ends with a specified substring
func EndsWithValidator(input interface{}, suffix string) error {
    str, ok := input.(string)
    if !ok {
        return fmt.Errorf("input must be a string")
    }
    if !strings.HasSuffix(str, suffix) {
        return fmt.Errorf("input must end with '%s'", suffix)
    }
    return nil
}

func convertToFloat64(input interface{}, threshold string) (float64, float64, error) {
    val := reflect.ValueOf(input)
    if val.Kind() != reflect.Int && val.Kind() != reflect.Int64 && val.Kind() != reflect.Float64 && val.Kind() != reflect.Float32 {
        return 0, 0, fmt.Errorf("input must be a number")
    }

    thresh, err := strconv.ParseFloat(threshold, 64)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to convert threshold to float: %v", err)
    }

    number := val.Convert(reflect.TypeOf(thresh)).Float()
    return number, thresh, nil
}