package validators_test

import (
	"testing"

	"github.com/dev3mike/go-xmapper/validators"
)

func TestRequiredValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Non-empty string", "Hello", ""},
		{"Empty string", "", "input is required and cannot be empty"},
		{"Whitespace only", "   ", "input is required and cannot be empty"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.RequiredValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestEmailValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Email", "email@example.com", ""},
		{"Invalid Email", "email@.com", "input is not a valid email address"},
		{"Non-string input", 123, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.EmailValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestPhoneValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Phone", "+1234567890123", ""},
		{"Invalid Phone", "12345", "input is not a valid international phone number"},
		{"Non-string input", 123, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.PhoneValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestStrongPasswordValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Strong Password", "Strong1$Password", ""},
		{"Too Short", "Str1$P", "password must be at least 8 characters long"},
		{"No Uppercase", "strong1$password", "password must contain at least one uppercase letter"},
		{"No Lowercase", "STRONG1$PASSWORD", "password must contain at least one lowercase letter"},
		{"No Digit", "Strong$$Password", "password must contain at least one digit"},
		{"No Special", "Strong1Password", "password must contain at least one special character"},
		{"Non-string input", 12345, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.StrongPasswordValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Test '%s' failed: Expected error '%s', got '%v'", tc.name, tc.expect, err)
			}
		})
	}
}

func TestDateValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Date", "2024-01-01", ""},
		{"Invalid Date", "01-01-2024", "input is not a valid date, expected format YYYY-MM-DD"},
		{"Non-string input", 123, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.DateValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestTimeValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Time", "23:59:59", ""},
		{"Invalid Time", "23:59:60", "input is not a valid time, expected format HH:MM:SS"},
		{"Non-string input", 123, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.TimeValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestDatetimeValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Datetime", "2021-12-31 23:59:59 PST", ""},
		{"Invalid Datetime", "2021-12-31 25:00:00 PST", "input is not a valid datetime with timezone, expected format YYYY-MM-DD HH:MM:SS MST"},
		{"Wrong Format", "31-12-2021 23:59:59 PST", "input is not a valid datetime with timezone, expected format YYYY-MM-DD HH:MM:SS MST"},
		{"Non-string input", 12345, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.DatetimeValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestUrlValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid URL", "https://www.example.com", ""},
		{"Valid URL", "http://www.example.com", ""},
		{"Valid URL", "http://example.com", ""},
		{"Invalid URL", "www.example.com", "input is not a valid URL"},
		{"Non-string input", 12345, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.UrlValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestIpValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid IP", "192.168.1.1", ""},
		{"Invalid IP", "999.999.999.999", "input is not a valid IP address"},
		{"Non-string input", 12345, "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.IpValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestMinLengthValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		length string
		expect string
	}{
		{"Valid Length", "HelloWorld", "5", ""},
		{"Invalid Length", "Hello", "10", "input does not meet the minimum length requirement, minimum length is 10"},
		{"Non-string input", 12345, "5", "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.MinLengthValidator(tc.input, tc.length)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestMaxLengthValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		length string
		expect string
	}{
		{"Valid Length", "Hello", "10", ""},
		{"Invalid Length", "HelloWorld", "5", "input exceeds the maximum length requirement, maximum length is 5"},
		{"Non-string input", 12345, "10", "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.MaxLengthValidator(tc.input, tc.length)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestGreaterThanValidator(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		threshold string
		expect    string
	}{
		{"Valid Greater", 10.5, "10", ""},
		{"Valid Greater", 16, "10", ""},
		{"Invalid Greater", 5.0, "10", "input must be greater than 10"},
		{"Non-float Input", "10", "5", "input must be a number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.GreaterThanValidator(tc.input, tc.threshold)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestLessThanValidator(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		threshold string
		expect    string
	}{
		{"Valid Less", 5.0, "10", ""},
		{"Valid Less", 3, "13", ""},
		{"Invalid Less", 10.5, "10", "input must be less than 10"},
		{"Non-float Input", "10", "15", "input must be a number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.LessThanValidator(tc.input, tc.threshold)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestGreaterThanOrEqualValidator(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		threshold string
		expect    string
	}{
		{"Valid Greater Or Equal", 10.0, "10", ""},
		{"Valid Greater Or Equal", 13, "10", ""},
		{"Invalid Greater Or Equal", 5.0, "10", "input must be greater than or equal to 10"},
		{"Non-float Input", "10", "5", "input must be a number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.GreaterThanOrEqualValidator(tc.input, tc.threshold)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestLessThanOrEqualValidator(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		threshold string
		expect    string
	}{
		{"Valid Less Or Equal", 10.0, "10", ""},
		{"Valid Less Or Equal", 8, "10", ""},
		{"Invalid Less Or Equal", 15.0, "10", "input must be less than or equal to 10"},
		{"Non-float Input", "10", "15", "input must be a number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.LessThanOrEqualValidator(tc.input, tc.threshold)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestEnumValidator(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		allowedValues string
		expect        string
	}{
		{"Valid Enum", "apple", "apple-banana-orange", ""},
		{"Invalid Enum", "pear", "apple-banana-orange", "input must be one of the following values: apple-banana-orange"},
		{"Non-string Input", 12345, "apple-banana-orange", "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.EnumValidator(tc.input, tc.allowedValues)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestBooleanValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		expect string
	}{
		{"Valid Boolean True", true, ""},
		{"Valid Boolean False", false, ""},
		{"Invalid Boolean", "true", "input must be a boolean"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.BooleanValidator(tc.input, "")
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestBooleanPointerValidator(t *testing.T) {
	boolean := true

	err := validators.BooleanValidator(&boolean, "")
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func TestContainsValidator(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		allowedValues string
		expect        string
	}{
		{"Contains Allowed", "hello world", "world,universe", ""},
		{"Does Not Contain", "hello world", "test,universe", "input must contain one of the following values: test,universe"},
		{"Non-string Input", 12345, "hello,world", "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.ContainsValidator(tc.input, tc.allowedValues)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestNotContainsValidator(t *testing.T) {
	tests := []struct {
		name             string
		input            interface{}
		disallowedValues string
		expect           string
	}{
		{"Does Not Contain Disallowed", "hello world", "test,universe", ""},
		{"Contains Disallowed", "hello world", "hello,test", "input must not contain the following value: hello"},
		{"Non-string Input", 12345, "hello,world", "failed to map the input to a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.NotContainsValidator(tc.input, tc.disallowedValues)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestRangeValidator(t *testing.T) {

	tests := []struct {
		name     string
		input    interface{}
		rangeStr string
		expect   string
	}{
		{"Within Range", 50.0, "10-100", ""},
		{"Below Range", 5.0, "10-100", "input must be between 10 and 100"},
		{"Above Range", 150.0, "10-100", "input must be between 10 and 100"},
		{"Above Range", 9, "10-100", "input must be between 10 and 100"},
		{"Invalid Range Format", 50.0, "100-10", "minimum value must be less than maximum value"},
		{"Non-float Input", "50", "10-100", "input must be a number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.RangeValidator(tc.input, tc.rangeStr)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestStartsWithValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		prefix string
		expect string
	}{
		{"Starts With", "hello world", "hello", ""},
		{"Does Not Start With", "hello world", "world", "input must start with 'world'"},
		{"Non-string Input", 12345, "hello", "input must be a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.StartsWidthValidator(tc.input, tc.prefix)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}

func TestEndsWithValidator(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		suffix string
		expect string
	}{
		{"Ends With", "hello world", "world", ""},
		{"Does Not End With", "hello world", "hello", "input must end with 'hello'"},
		{"Non-string Input", 12345, "world", "input must be a string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.EndsWithValidator(tc.input, tc.suffix)
			if (err != nil && err.Error() != tc.expect) || (err == nil && tc.expect != "") {
				t.Errorf("Expected error '%s', got '%v'", tc.expect, err)
			}
		})
	}
}
