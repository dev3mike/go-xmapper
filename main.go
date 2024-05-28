package xmapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dev3mike/go-xmapper/transformers"
	"github.com/dev3mike/go-xmapper/validators"
)

// TransformerFunc defines the type for functions that transform data from one form to another.
type TransformerFunc func(interface{}) interface{}

// ValidatorFunc defines the type for functions that validate data.
type ValidatorFunc func(interface{}, string) error

// ErrValidation: Validation methods return this error in case of an error, so you can use it to catch validation errors
var ErrValidation = errors.New("ValidationError")

// transformerRegistry is a map that holds registered transformer functions keyed by their name.
var transformerRegistry = map[string]TransformerFunc{}

// validatorRegistry holds registered validator functions keyed by their name.
var validatorRegistry = map[string]ValidatorFunc{}

func init() {
	// Default validators
	RegisterValidator("required", validators.RequiredValidator) // Should not be empty
	RegisterValidator("email", validators.EmailValidator)
	RegisterValidator("phone", validators.PhoneValidator)                   // International phone number format
	RegisterValidator("strongPassword", validators.StrongPasswordValidator) // Minimum 8 characters, at least one uppercase, one lowercase, one number, and one special character
	RegisterValidator("date", validators.DateValidator)                     // Date in YYYY-MM-DD format
	RegisterValidator("time", validators.TimeValidator)                     // Time in HH:MM:SS format
	RegisterValidator("datetime", validators.DatetimeValidator)             // Date and time in YYYY-MM-DD HH:MM:SS format with timezone
	RegisterValidator("url", validators.UrlValidator)
	RegisterValidator("ip", validators.IpValidator)
	RegisterValidator("minLength", validators.MinLengthValidator)
	RegisterValidator("maxLength", validators.MaxLengthValidator)
	RegisterValidator("gt", validators.GreaterThanValidator)
	RegisterValidator("lt", validators.LessThanValidator)
	RegisterValidator("gte", validators.GreaterThanOrEqualValidator)
	RegisterValidator("lte", validators.LessThanOrEqualValidator)
	RegisterValidator("range", validators.RangeValidator)
	RegisterValidator("enum", validators.EnumValidator)
	RegisterValidator("boolean", validators.BooleanValidator)
	RegisterValidator("contains", validators.ContainsValidator)
	RegisterValidator("notContains", validators.NotContainsValidator)
	RegisterValidator("startsWidth", validators.StartsWidthValidator)
	RegisterValidator("endsWith", validators.EndsWithValidator)

	// Default transformers
	RegisterTransformer("uppercase", transformers.ToUpperCase)
	RegisterTransformer("lowercase", transformers.ToLowerCase)
	RegisterTransformer("trim", transformers.Trim)
	RegisterTransformer("trimLeft", transformers.TrimLeft)
	RegisterTransformer("trimRight", transformers.TrimRight)
	RegisterTransformer("base64Encode", transformers.Base64Encode)
	RegisterTransformer("base64Decode", transformers.Base64Decode)
	RegisterTransformer("urlEncode", transformers.UrlEncode)
	RegisterTransformer("urlDecode", transformers.UrlDecode)
}

// RegisterTransformer adds a transformer function to the registry with a given name.
func RegisterTransformer(name string, f TransformerFunc) {
	transformerRegistry[name] = f
}

// RegisterValidator adds a validator function to the registry.
func RegisterValidator(name string, f ValidatorFunc) {
	validatorRegistry[name] = f
}

// MapStructs validate, transfor and maps data from source struct to destination struct
func MapStructs(src, dest interface{}) error {

	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)
	if !isValidStructPointer(srcValue) || !isValidStructPointer(destValue) {
		return errors.New("both source and destination must be pointer to a struct")
	}

	return mapStructsRecursive(srcValue, destValue)
}

// MapSliceOfStructs iterate over the source slice and map each struct to the destination slice
func MapSliceOfStructs(src, dest interface{}) error {

	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	if srcValue.Kind() != reflect.Pointer || srcValue.Elem().Kind() != reflect.Slice {
		return errors.New("source must be a pointer to a slice")
	}
	if destValue.Kind() != reflect.Pointer || destValue.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	srcSlice := srcValue.Elem()
	destElemType := destValue.Elem().Type().Elem().Elem()
	destSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(destElemType)), srcSlice.Len(), srcSlice.Len())

	for i := 0; i < srcSlice.Len(); i++ {
		srcElem := srcSlice.Index(i)
		destElem := reflect.New(destElemType)
		err := MapStructs(srcElem.Interface(), destElem.Interface())
		if err != nil {
			return err
		}
		destSlice.Index(i).Set(destElem)
	}

	destValue.Elem().Set(destSlice)
	return nil
}

// MapJsonStruct decodes a JSON string into the provided struct pointer and applies any necessary validations and transformations
func MapJsonStruct(jsonStr string, target interface{}) error {
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	err := json.Unmarshal([]byte(jsonStr), target)
	if err != nil {
		return err
	}

	return MapStructs(target, target)
}

/**
    * validatorAndTransformerSpec example : "validators:'arg1,arg2:value'transformers:'transformer1,transformer2'"
**/
func ValidateSingleField(value interface{}, validatorAndTransformerSpec string) (interface{}, error) {
	validatorsStr, transformersStr := parseSingleFieldValidatorAndTransformerSpec(validatorAndTransformerSpec)

	if len(validatorsStr) > 0 {
		validators, err := parseFieldValidators(validatorsStr)
		if err != nil {
			return value, err
		}

		for _, validator := range validators {
			if err := validator(value); err != nil {
				return value, fmt.Errorf("validation failed: %w", ErrValidation)
			}
		}
	}

	if len(transformersStr) > 0 {

		transformers, err := parseTransformers(transformersStr)
		if err != nil {
			return value, err
		}

		for _, transformer := range transformers {
			value = transformer(value)
		}

	}

	return value, nil
}

// ValidateStruct validates the struct fields against defined validators.
func ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	if !isValidStructPointer(val) {
		return fmt.Errorf("input must be a pointer to a struct")
	}
	return validateStructRecursive(val)
}

// validateStructRecursive recursively validates each field of a struct.
func validateStructRecursive(val reflect.Value) error {
	structFields := val.Elem()

	structFieldMap := buildDestinationFieldMap(structFields)
	transformers, err := findTransformers(structFields)
	if err != nil {
		return err
	}

	validators, err := findValidators(structFields)
	if err != nil {
		return err
	}

	for i := 0; i < structFields.NumField(); i++ {
		field := structFields.Field(i)
		fieldName := getFieldName(structFields.Type().Field(i), "json")

		if fieldValidators, ok := validators[fieldName]; ok {
			for _, validator := range fieldValidators {
				if err := validator(field.Interface()); err != nil {
					return fmt.Errorf("validation failed for field '%s': %w", fieldName, ErrValidation)
				}
			}
		}

		if structField, ok := structFieldMap[fieldName]; ok && structField.CanSet() {
			if err := setFieldValue(structField, structField, transformers[fieldName]); err != nil {
				return err
			}
		}
	}
	return nil
}

// mapStructsRecursive recursively maps data from source to destination structs.
func mapStructsRecursive(srcVal, destVal reflect.Value) error {
	srcFields := srcVal.Elem()
	destFields := destVal.Elem()

	// Build destination field map and fetch transformers and validators
	destMap := buildDestinationFieldMap(destFields)
	transformers, err := findTransformers(srcFields)
	if err != nil {
		return err
	}

	validators, err := findValidators(srcFields)
	if err != nil {
		return err
	}

	// Iterate through each source field
	for i := 0; i < srcFields.NumField(); i++ {
		srcField := srcFields.Field(i)
		fieldName := getFieldName(srcFields.Type().Field(i), "json")

		if fieldName == "" {
			continue
		}

		// Execute validators for the field if any are defined
		if fieldValidators, ok := validators[fieldName]; ok {
			for _, validator := range fieldValidators {
				if err := validator(srcField.Interface()); err != nil {
					return fmt.Errorf("validation failed for field '%s': %w", fieldName, ErrValidation)
				}
			}
		}

		// If a corresponding destination field exists and can be set, apply transformers and set value
		if destField, ok := destMap[fieldName]; ok && destField.CanSet() {
			if err := setFieldValue(srcField, destField, transformers[fieldName]); err != nil {
				return err
			}
		}
	}
	return nil
}

// isValidStructPointer checks if the provided value is a pointer to a struct.
func isValidStructPointer(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}

// buildDestinationFieldMap creates a map of destination fields keyed by their JSON tag names.
func buildDestinationFieldMap(destFields reflect.Value) map[string]reflect.Value {
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < destFields.NumField(); i++ {
		field := destFields.Type().Field(i)
		fieldName := getFieldName(field, "json")
		if fieldName != "" {
			fieldMap[fieldName] = destFields.Field(i)
		}
	}
	return fieldMap
}

// getFieldName returns the first part of a struct field's tag associated with the provided key or an empty string if not set.
func getFieldName(field reflect.StructField, key string) string {
	tag := field.Tag.Get(key)
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

// findTransformers collects lists of transformers for fields that have a transformer tag specified.
// It returns an error if any specified transformer does not exist.
func findTransformers(fields reflect.Value) (map[string][]TransformerFunc, error) {
	transformers := make(map[string][]TransformerFunc)
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Type().Field(i)
		transformerNames := field.Tag.Get("transformers")
		if transformerNames != "" {
			jsonName := getFieldName(field, "json")
			transformerList, err := parseTransformers(transformerNames)
			if err != nil {
				return nil, err
			}
			transformers[jsonName] = transformerList
		}
	}
	return transformers, nil
}

// parseTransformers parses a comma-separated list of transformer names and returns a slice of TransformerFunc.
// It returns an error if any transformer cannot be found in the registry.
func parseTransformers(names string) ([]TransformerFunc, error) {
	nameList := strings.Split(names, ",")
	transformerList := make([]TransformerFunc, 0, len(nameList))
	for _, name := range nameList {
		name = strings.TrimSpace(name)
		if transformer, exists := transformerRegistry[name]; exists {
			transformerList = append(transformerList, transformer)
		} else {
			return nil, fmt.Errorf("transformer '%s' not found", name)
		}
	}
	return transformerList, nil
}

func setFieldValue(srcField, destField reflect.Value, transformers []TransformerFunc) error {
	// Handle pointers
	if srcField.Kind() == reflect.Ptr {
		if srcField.IsNil() {
			// Set destination field to nil if source is nil
			destField.Set(reflect.Zero(destField.Type()))
			return nil
		}
		srcField = srcField.Elem()
	}
	if destField.Kind() == reflect.Ptr {
		if destField.IsNil() {
			// Initialize destination pointer if it's nil
			destField.Set(reflect.New(destField.Type().Elem()))
		}
		destField = destField.Elem()
	}

	// Handle time.Time to time.Time conversion
	if srcField.Type() == reflect.TypeOf(time.Time{}) && destField.Type() == reflect.TypeOf(time.Time{}) {
		srcTime := srcField.Interface().(time.Time)
		destField.Set(reflect.ValueOf(srcTime))
		return nil
	}

	// Handle struct to JSON string conversion
	if srcField.Kind() == reflect.Struct && destField.Kind() == reflect.String {
		jsonBytes, err := json.Marshal(srcField.Interface())
		if err != nil {
			return err
		}
		destField.SetString(string(jsonBytes))
		return nil
	}

	if srcField.Kind() == reflect.Struct && destField.Kind() == reflect.Struct {
		return mapStructsRecursive(srcField.Addr(), destField.Addr())
	}

	if srcField.Kind() == reflect.Slice && destField.Kind() == reflect.Slice {
		destElemType := destField.Type().Elem()
		convertedSlice := reflect.MakeSlice(destField.Type(), srcField.Len(), srcField.Cap())

		for i := 0; i < srcField.Len(); i++ {
			srcElem := srcField.Index(i)
			convertedElem := reflect.New(destElemType).Elem()

			// Convert the element recursively or use transformers if needed
			if err := setFieldValue(srcElem, convertedElem, transformers); err != nil {
				return err
			}

			convertedSlice.Index(i).Set(convertedElem)
		}

		destField.Set(convertedSlice)
		return nil
	}

	// Handle JSON string to struct conversion
	if srcField.Kind() == reflect.String && destField.Kind() == reflect.Struct {
		jsonStr := srcField.String()

		if len(jsonStr) == 0 {
			jsonStr = "{}"
		}

		structValue := reflect.New(destField.Type()).Interface()
		if err := json.Unmarshal([]byte(jsonStr), structValue); err != nil {
			return err
		}
		destField.Set(reflect.ValueOf(structValue).Elem())
		return nil
	}

	// Handle JSON string to slice conversion
	if srcField.Kind() == reflect.String && destField.Kind() == reflect.Slice {
		jsonStr := srcField.String()

		if len(jsonStr) == 0 {
			jsonStr = "[]"
		}

		sliceValue := reflect.New(destField.Type()).Interface()
		if err := json.Unmarshal([]byte(jsonStr), sliceValue); err != nil {
			return err
		}
		destField.Set(reflect.ValueOf(sliceValue).Elem())
		return nil
	}

	// Handle slice to JSON string conversion
	if srcField.Kind() == reflect.Slice && destField.Kind() == reflect.String {
		jsonBytes, err := json.Marshal(srcField.Interface())
		if err != nil {
			return err
		}
		destField.SetString(string(jsonBytes))
		return nil
	}

	// Handle string to slice of strings conversion
	if srcField.Kind() == reflect.String && destField.Kind() == reflect.Slice && destField.Type().Elem().Kind() == reflect.String {
		str := srcField.String()
		slice := strings.Split(str, ",")
		destSlice := reflect.MakeSlice(destField.Type(), len(slice), len(slice))
		for i, v := range slice {
			destSlice.Index(i).Set(reflect.ValueOf(strings.TrimSpace(v)))
		}
		destField.Set(destSlice)
		return nil
	}

	// Handle slice of strings to single string conversion (joining with commas)
	if srcField.Kind() == reflect.Slice && srcField.Type().Elem().Kind() == reflect.String && destField.Kind() == reflect.String {
		slice := srcField.Interface().([]string)
		joinedStr := strings.Join(slice, ",")
		destField.SetString(joinedStr)
		return nil
	}

	// Apply transformers if any and set the value
	valueToSet := srcField.Interface()
	for _, transformer := range transformers {
		valueToSet = transformer(valueToSet)
	}
	destField.Set(reflect.ValueOf(valueToSet))
	return nil
}

func findValidators(fields reflect.Value) (map[string][]func(interface{}) error, error) {
	validators := make(map[string][]func(interface{}) error)
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Type().Field(i)
		validatorSpec := field.Tag.Get("validators")
		if validatorSpec == "" {
			continue
		}

		jsonName := getFieldName(field, "json")
		fieldValidators, err := parseFieldValidators(validatorSpec)
		if err != nil {
			return nil, fmt.Errorf("error parsing validators for field '%s': %v", jsonName, err)
		}
		validators[jsonName] = fieldValidators
	}
	return validators, nil
}

func parseFieldValidators(validatorSpec string) ([]func(interface{}) error, error) {
	var validators []func(interface{}) error
	validatorEntries := strings.Split(validatorSpec, ",")
	for _, entry := range validatorEntries {
		parts := strings.SplitN(entry, ":", 2)
		validatorName := strings.TrimSpace(parts[0])
		arg := ""
		if len(parts) > 1 {
			arg = strings.TrimSpace(parts[1])
		}

		validatorFunc, exists := validatorRegistry[validatorName]
		if !exists {
			return nil, fmt.Errorf("validator '%s' not found", validatorName)
		}

		// Wrap the validator function to include its argument
		validators = append(validators, func(value interface{}) error {
			return validatorFunc(value, arg)
		})
	}
	return validators, nil
}

// Helper function to extract values from input based on a given prefix.
func extractValueFromInput(input, prefix string) string {
	start := strings.Index(input, prefix)
	if start == -1 {
		return ""
	}

	end := strings.Index(input[start+len(prefix):], "'") + start + len(prefix)
	if end == -1 {
		return ""
	}

	return input[start+len(prefix) : end]
}

func parseSingleFieldValidatorAndTransformerSpec(input string) (string, string) {
	validators := extractValueFromInput(input, "validators:'")
	transformers := extractValueFromInput(input, "transformers:'")
	return validators, transformers
}
