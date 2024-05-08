package xmapper_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/dev3mike/go-xmapper"
)

// Define dummy transformers
func toUpperCase(input interface{}) interface{} {
    if str, ok := input.(string); ok {
        return strings.ToUpper(str)
    }
    return input
}

func addExclamation(input interface{}) interface{} {
    if str, ok := input.(string); ok {
        return str + "!"
    }
    return input
}

func repeatTwice(input interface{}) interface{} {
    if str, ok := input.(string); ok {
        return str + " " + str
    }
    return input
}

// Define dummy validatos
func isEmail(input interface{}) bool {
	email, ok := input.(string)
	if !ok {
		return false
	}

	regexPattern := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`
	
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}
	
	return re.MatchString(email)
}

func isGmailAddress(input interface{}) bool {
    email, ok := input.(string)
	if !ok {
		return false
	}
	return strings.HasSuffix(email, "@gmail.com")
}

// TestMapStructsBasic checks the basic functionality of mapping without transformations.
func TestMapStructsBasic(t *testing.T) {
    type Src struct {
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
    }
    type Dest struct {
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
    }

    src := Src{FirstName: "John", LastName: "Doe"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error: %s", err)
    }
    if dest.FirstName != "John" || dest.LastName != "Doe" {
        t.Errorf("Failed to map fields correctly, got: %+v", dest)
    }
}

// TestMapStructsTransformations checks the transformation functionality.
func TestMapStructsTransformations(t *testing.T) {
    xmapper.RegisterTransformer("toUpperCase", toUpperCase)

    type Src struct {
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName" transformer:"toUpperCase"`
    }
    type Dest struct {
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
    }

    src := Src{FirstName: "John", LastName: "Doe"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error: %s", err)
    }
    if dest.FirstName != "John" || dest.LastName != "DOE" {
        t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
    }
}

// TestMapStructsInvalidInput tests the error handling for non-struct pointers.
func TestMapStructsInvalidInput(t *testing.T) {
    var src struct {
        FirstName string `json:"firstName"`
    }
    dest := "not a struct pointer"

    err := xmapper.MapStructs(&src, dest)
    if err == nil {
        t.Errorf("Expected error for non-pointer destination, got nil")
    }
}

// TestMapStructsNoFieldMatch ensures the function does not crash or misbehave when no fields match.
func TestMapStructsNoFieldMatch(t *testing.T) {
    type Src struct {
        FirstName string `json:"firstName"`
    }
    type Dest struct {
        LastName string `json:"lastName"`
    }

    src := Src{FirstName: "John"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error when no fields match: %s", err)
    }
    if dest.LastName != "" {
        t.Errorf("Fields that do not match should not be altered, got: %+v", dest)
    }
}

// TestMapStructsUnregisteredTransformer tests the handling of unregistered transformer names.
func TestMapStructsUnregisteredTransformer(t *testing.T) {
    type Src struct {
        LastName string `json:"lastName" ,transformer:"nonExistent"`
    }
    type Dest struct {
        LastName string `json:"lastName"`
    }

    src := Src{LastName: "Doe"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Expected mapping to continue despite unregistered transformer, error: %s", err)
    }
    if dest.LastName != "Doe" {
        t.Errorf("Expected original value to be set when transformer is unregistered, got: %+v", dest)
    }
}

// TestMapStructsComplexNestedStructs tests mapping between complex nested structures.
func TestMapStructsComplexNestedStructs(t *testing.T) {
    type ContactInfo struct {
        Email   string `json:"email"`
        ZipCode string `json:"zipCode"`
        Tags []string `json:"tags"`
    }
    type ContactInfo2 struct {
        Email   string `json:"email"`
        ZipCode string `json:"zipCode"`
        Tags []string `json:"tags"`
    }
    type Src struct {
        Email   string `json:"email"`
        UserID  string `json:"userId"`
        Contact  ContactInfo `json:"contact"`
    }
    type Dest struct {
        Email   string `json:"email"`
        UserID  string `json:"userId"`
        Contact  ContactInfo2 `json:"contact"`
    }

    src := Src{
        UserID: "12345",
        Email: "john.doe2@example.com",
        Contact: ContactInfo{
            Email:   "john.doe@example.com",
            ZipCode: "90210",
        }}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)

    if err != nil {
        t.Errorf("Unexpected error when mapping nested fields: %s", err)
    }
    if dest.Email == "" || dest.UserID == "" || src.Contact.Email != dest.Contact.Email {
        t.Errorf("Failed to map nested fields, got: %+v", dest)
    }
    if !reflect.DeepEqual(dest.Contact.Tags, src.Contact.Tags) {
        t.Errorf("Failed to map nested array fields correctly, got: %+v, want: %+v", dest.Contact.Tags, src.Contact.Tags)
    }
}

// TestMapStructsWithArrays tests mapping structures that contain slice or array fields.
func TestMapStructsWithArrays(t *testing.T) {
    type Src struct {
        Tags []string `json:"tags"`
    }
    type Dest struct {
        Tags []string `json:"tags"`
    }

    src := Src{Tags: []string{"go", "programming", "developer"}}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error when mapping array fields: %s", err)
    }
    if !reflect.DeepEqual(dest.Tags, src.Tags) {
        t.Errorf("Failed to map array fields correctly, got: %+v, want: %+v", dest.Tags, src.Tags)
    }
}

// TestMapStructsPointerFields tests mapping structs with pointer fields.
func TestMapStructsPointerFields(t *testing.T) {
    type Src struct {
        Name *string `json:"name"`
    }
    type Dest struct {
        Name *string `json:"name"`
    }

    name := "John Doe"
    src := Src{Name: &name}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error when mapping pointer fields: %s", err)
    }
    if dest.Name == nil || *dest.Name != "John Doe" {
        t.Errorf("Failed to map pointer fields correctly, expected 'John Doe', got: %v", dest.Name)
    }
}

func TestMultipleTransformers(t *testing.T) {
    // Register transformers
    xmapper.RegisterTransformer("toUpperCase", toUpperCase)
    xmapper.RegisterTransformer("addExclamation", addExclamation)
    xmapper.RegisterTransformer("repeatTwice", repeatTwice)

    // Define test struct with transformer tags
    type TestStruct struct {
        Message string `json:"message" transformer:"toUpperCase,addExclamation,repeatTwice"`
    }

    src := TestStruct{Message: "hello"}
    dest := TestStruct{}

    // Perform mapping
    err := xmapper.MapStructs(&src, &dest)
    if err != nil {
        t.Errorf("Unexpected error during mapping: %s", err)
    }

    // Expected result after applying all transformers
    expected := "HELLO! HELLO!"
    if dest.Message != expected {
        t.Errorf("Expected '%s', got '%s'", expected, dest.Message)
    }
}


// TestNonExistentTransformer checks if using a non-existent transformer results in a proper error.
func TestNonExistentTransformer(t *testing.T) {
    // Register only valid transformers
    xmapper.RegisterTransformer("toUpperCase", toUpperCase)
    xmapper.RegisterTransformer("addExclamation", addExclamation)

    // Define a test struct that specifies a non-existent transformer
    type TestStruct struct {
        Message string `json:"message" transformer:"nonExistentTransformer"`
    }

    src := TestStruct{Message: "test"}
    dest := TestStruct{}

    // Perform mapping
    err := xmapper.MapStructs(&src, &dest)
    if err == nil {
        t.Error("Expected error for non-existent transformer, but got none")
    } else {
        expectedErrorMessage := "transformer 'nonExistentTransformer' not found"
        if err.Error() != expectedErrorMessage {
            t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
        }
    }
}


// TestMapStructsValidatorsWithInvalidData checks the validation functionality.
func TestMapStructsValidatorsWithInvalidData(t *testing.T) {
    xmapper.RegisterValidator("isEmail", isEmail)

    type Src struct {
        Email  string `json:"email" validator:"isEmail" transformer:"toUpperCase"`
    }
    type Dest struct {
        EmailAddress  string `json:"email"`
    }

    src := Src{Email: "not_a_valid_email"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)

    if err == nil {
        t.Errorf("Expected error for invalid email, got nil")
    }
}

// TestMapStructsValidatorsAndTransformersWithValidData checks the validation functionality.
func TestMapStructsValidatorsAndTransformersWithValidData(t *testing.T) {
    xmapper.RegisterValidator("isEmail", isEmail)

    type Src struct {
        Email  string `json:"email" validator:"isEmail" transformer:"toUpperCase"`
    }
    type Dest struct {
        EmailAddress  string `json:"email"`
    }

    src := Src{Email: "test@gmail.com"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)

    if err != nil {
        t.Errorf("Unexpected error for valid email: %s", err)
    }

    if dest.EmailAddress != "TEST@GMAIL.COM"{
        t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
    }
}

// TestMapStructsValidatorsAndTransformersWithMultipleValidators checks the validation functionality.
func TestMapStructsValidatorsAndTransformersWithMultipleValidators(t *testing.T) {
    xmapper.RegisterValidator("isEmail", isEmail)
    xmapper.RegisterValidator("isGmailAddress", isGmailAddress)

    type Src struct {
        Email  string `json:"email" validator:"isEmail,isGmailAddress" transformer:"toUpperCase"`
    }
    type Dest struct {
        EmailAddress  string `json:"email"`
    }

    src := Src{Email: "test@gmail.com"}
    dest := Dest{}

    err := xmapper.MapStructs(&src, &dest)

    if err != nil {
        t.Errorf("Unexpected error for valid email: %s", err)
    }

    if dest.EmailAddress != "TEST@GMAIL.COM"{
        t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
    }
}