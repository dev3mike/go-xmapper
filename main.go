package xmapper

import (
	"fmt"
	"reflect"
	"strings"
)

// TransformerFunc defines the type for functions that transform data from one form to another.
type TransformerFunc func(interface{}) interface{}

// ValidatorFunc defines the type for functions that validate data.
type ValidatorFunc func(interface{}) bool

// transformerRegistry is a map that holds registered transformer functions keyed by their name.
var transformerRegistry = map[string]TransformerFunc{}

// validatorRegistry holds registered validator functions keyed by their name.
var validatorRegistry = map[string]ValidatorFunc{}

// RegisterTransformer adds a transformer function to the registry with a given name.
func RegisterTransformer(name string, f TransformerFunc) {
    transformerRegistry[name] = f
}

// RegisterValidator adds a validator function to the registry.
func RegisterValidator(name string, f ValidatorFunc) {
    validatorRegistry[name] = f
}

// MapStructs maps data from source struct to destination struct using reflection.
func MapStructs(src, dest interface{}) error {
    srcValue := reflect.ValueOf(src)
    destValue := reflect.ValueOf(dest)
    if !isValidStructPointer(srcValue) || !isValidStructPointer(destValue) {
        return fmt.Errorf("both source and destination must be pointer to a struct")
    }
    return mapStructsRecursive(srcValue, destValue)
}

// mapStructsRecursive recursively maps data from source to destination structs.
func mapStructsRecursive(srcVal, destVal reflect.Value) error {
    srcFields := srcVal.Elem()
    destFields := destVal.Elem()

    destMap := buildDestinationFieldMap(destFields)
    transformers, err := findTransformers(srcFields)
    if err != nil {
        return err
    }

    validators, err := findValidators(srcFields)
    if err != nil {
        return err
    }

    for i := 0; i < srcFields.NumField(); i++ {
        srcField := srcFields.Field(i)
        fieldName := getFieldName(srcFields.Type().Field(i), "json")
        if fieldName == "" {
            continue
        }

        if validatorsList, ok := validators[fieldName]; ok {
            for _, validator := range validatorsList {
                if !validator(srcField.Interface()) {
                    return fmt.Errorf("validation failed for field '%s'", fieldName)
                }
            }
        }

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
        transformerNames := field.Tag.Get("transformer")
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

// setFieldValue sets the destination field value from the source field, potentially using multiple transformers.
func setFieldValue(srcField, destField reflect.Value, transformers []TransformerFunc) error {
    if srcField.Kind() == reflect.Struct && destField.Kind() == reflect.Struct {
        return mapStructsRecursive(srcField.Addr(), destField.Addr())
    }

    valueToSet := srcField.Interface()
    for _, transformer := range transformers {
        valueToSet = transformer(valueToSet)
    }
    destField.Set(reflect.ValueOf(valueToSet))
    return nil
}

func findValidators(fields reflect.Value) (map[string][]ValidatorFunc, error) {
    validators := make(map[string][]ValidatorFunc)
    for i := 0; i < fields.NumField(); i++ {
        field := fields.Type().Field(i)
        validatorNames := field.Tag.Get("validator")
        if validatorNames != "" {
            jsonName := getFieldName(field, "json")
            validatorList, err := parseValidators(validatorNames)
            if err != nil {
                return nil, err
            }
            validators[jsonName] = validatorList
        }
    }
    return validators, nil
}

func parseValidators(names string) ([]ValidatorFunc, error) {
    nameList := strings.Split(names, ",")
    validatorList := make([]ValidatorFunc, 0, len(nameList))
    for _, name := range nameList {
        name = strings.TrimSpace(name)
        if validator, exists := validatorRegistry[name]; exists {
            validatorList = append(validatorList, validator)
        } else {
            return nil, fmt.Errorf("validator '%s' not found", name)
        }
    }
    return validatorList, nil
}