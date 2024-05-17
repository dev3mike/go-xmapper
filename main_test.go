package xmapper_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

// Define dummy validators
func isEmail(input interface{}, _ string) error {

	email, ok := input.(string)
	if !ok {
		return fmt.Errorf("input is not a valid string")
	}

	regexPattern := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`
	re, err := regexp.Compile(regexPattern)

	if err != nil {
		return fmt.Errorf("failed to compile regex pattern: %v", err)
	}

	if !re.MatchString(email) {
		return fmt.Errorf("input is not a valid email")
	}

	return nil
}

func isGmailAddress(input interface{}, _ string) error {
	email, ok := input.(string)

	if !ok {
		return fmt.Errorf("failed to map the input to a string")
	}
	if !strings.HasSuffix(email, "@gmail.com") {
		return fmt.Errorf("input is not a valid gmail address")
	}

	return nil
}

func minLength(input interface{}, length string) error {
	textInput, ok := input.(string)

	if !ok {
		return fmt.Errorf("failed to map the input to a string")
	}

	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		return fmt.Errorf("failed to convert length to integer")
	}
	if len(textInput) < lengthInt {
		return fmt.Errorf("input does not meet the minimum length requirement, minimum length is %s", length)
	}
	return nil
}

// TestMapStructsBasic checks the basic functionality of mapping without transformations.
func TestMapStructsBasic(t *testing.T) {
	type Src struct {
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
	}
	type Dest struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	src := Src{FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890"}
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
		LastName  string `json:"lastName" transformers:"toUpperCase"`
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

// TestValidateSingleField checks the validation functionality for a single field.
func TestValidateSingleFieldWithAValidEmail(t *testing.T) {
	xmapper.RegisterValidator("isEmail", isEmail)
	xmapper.RegisterTransformer("toUpperCase", toUpperCase)

	value := "test@example.com"

	transformedValue, err := xmapper.ValidateSingleField(value, "validators:'isEmail' transformers:'toUpperCase'")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if transformedValue != "TEST@EXAMPLE.COM" {
		t.Errorf("Failed to apply transformations correctly, got: %+v", value)
	}
}

// TestValidateSingleFieldWithAnInvalidEmail checks the validation functionality for a single field with an invalid email.
func TestValidateSingleFieldWithAnInvalidEmail(t *testing.T) {
	xmapper.RegisterValidator("isEmail", isEmail)
	xmapper.RegisterTransformer("toUpperCase", toUpperCase)

	value := "invalid-email"

	_, err := xmapper.ValidateSingleField(value, "validators:'isEmail' transformers:'toUpperCase'")

	if err == nil {
		t.Error("Expected an error for invalid email, but got none")
	}
}

// TestValidateSingleFieldWithAnInvalidEmail checks the validation functionality for a single field with an invalid email.
func TestValidateSingleFieldWithMultipleValidators(t *testing.T) {
	value := "abc"

	_, err := xmapper.ValidateSingleField(value, "validators:'isEmail,minLength:5,maxLength:10'")

	value2 := "abc@test.com"

	_, err2 := xmapper.ValidateSingleField(value2, "validators:'isEmail,minLength:50,maxLength:10'")

	if err == nil {
		t.Error("Expected an error for invalid email and length constraints, but got none")
	}

	if err2 == nil {
		t.Error("Expected an error for invalid email and length constraints, but got none")
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
		LastName string `json:"lastName" transformers:"nonExistent"`
	}
	type Dest struct {
		LastName string `json:"lastName"`
	}

	src := Src{LastName: "Doe"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err == nil {
		t.Errorf("Expected error when the transformer is not registered but got %s", err)
	}
}

// TestMapStructsComplexNestedStructs tests mapping between complex nested structures.
func TestMapStructsComplexNestedStructs(t *testing.T) {
	type ContactInfo struct {
		Email   string   `json:"email"`
		ZipCode string   `json:"zipCode"`
		Tags    []string `json:"tags"`
	}
	type ContactInfo2 struct {
		Email   string   `json:"email"`
		ZipCode string   `json:"zipCode"`
		Tags    []string `json:"tags"`
	}
	type Src struct {
		Email   string      `json:"email"`
		UserID  string      `json:"userId"`
		Contact ContactInfo `json:"contact"`
	}
	type Dest struct {
		Email   string       `json:"email"`
		UserID  string       `json:"userId"`
		Contact ContactInfo2 `json:"contact"`
	}

	src := Src{
		UserID: "12345",
		Email:  "john.doe2@example.com",
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
		Message string `json:"message" transformers:"toUpperCase,addExclamation,repeatTwice"`
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
		Message string `json:"message" transformers:"nonExistentTransformer"`
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
		Email string `json:"email" validators:"isEmail" transformers:"toUpperCase"`
	}
	type Dest struct {
		EmailAddress string `json:"email"`
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
	xmapper.RegisterTransformer("toUpperCase", toUpperCase)

	type Src struct {
		Email  string `json:"email" validators:"isEmail" transformers:"toUpperCase"`
		Status string `json:"status" validators:"enum:active-inactive"`
	}
	type Dest struct {
		EmailAddress string `json:"email"`
	}

	src := Src{Email: "test@gmail.com", Status: "active"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if dest.EmailAddress != "TEST@GMAIL.COM" {
		t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
	}
}

// TestMapStructsValidatorsAndTransformersWithMultipleValidators checks the validation functionality.
func TestMapStructsValidatorsAndTransformersWithMultipleValidators(t *testing.T) {
	xmapper.RegisterValidator("isEmail", isEmail)
	xmapper.RegisterValidator("isGmailAddress", isGmailAddress)
	xmapper.RegisterValidator("minLength", minLength)

	type Src struct {
		Email string `json:"email" validators:"isEmail,isGmailAddress,minLength:4" transformers:"toUpperCase"`
	}
	type Dest struct {
		EmailAddress string `json:"email"`
	}

	src := Src{Email: "test@gmail.com"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if dest.EmailAddress != "TEST@GMAIL.COM" {
		t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
	}
}

// TestMapStructsValidatorsAndTransformersWithMultipleValidators checks the validation functionality.
func TestDynamicVariablesWithDefaultValidatorsWithValidEmail(t *testing.T) {
	type Src struct {
		Email string `json:"email" validators:"email,minLength:4" transformers:"toUpperCase"`
	}
	type Dest struct {
		EmailAddress string `json:"email"`
	}

	src := Src{Email: "test@gmail.com"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if dest.EmailAddress != "TEST@GMAIL.COM" {
		t.Errorf("Failed to apply transformations correctly, got: %+v", dest)
	}
}

// TestMapStructsValidatorsAndTransformersWithMultipleValidators checks the validation functionality.
func TestDynamicVariablesWithDefaultValidatorsWithInvalidEmail(t *testing.T) {
	type Src struct {
		Email string `json:"email" validators:"email,minLength:5" transformers:"toUpperCase"`
	}
	type Dest struct {
		EmailAddress string `json:"email"`
	}

	src := Src{Email: "test"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)

	if err == nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

// TestDynamicVariablesWithDefaultValidatorsWithMaxLength checks the validation functionality.
func TestDynamicVariablesWithDefaultValidatorsWithMaxLength(t *testing.T) {
	type Src struct {
		Text string `json:"text" validators:"maxLength:5" transformers:"toUpperCase"`
	}
	type Dest struct {
		Text string `json:"email"`
	}

	src := Src{Text: "more_than_5"}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)

	if err == nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

// TestMapStructsStructSliceToJson tests mapping structures that contain slice or array fields.
func TestMapStructsStructSliceToJson(t *testing.T) {
	type Tag struct {
		Name string `json:"name"`
	}

	type Src struct {
		Tags []Tag `json:"tags"`
	}

	type Dest struct {
		Tags string `json:"tags"`
	}

	src := Src{Tags: []Tag{
		{Name: "go"},
		{Name: "programming"},
		{Name: "developer"},
	}}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error when mapping array fields: %s", err)
	}

	// Convert the source slice to JSON for comparison
	expectedTags, err := json.Marshal(src.Tags)
	if err != nil {
		t.Fatalf("Failed to marshal source tags to JSON: %s", err)
	}

	if dest.Tags != string(expectedTags) {
		t.Errorf("Failed to map array fields correctly, got: %s, want: %s", dest.Tags, string(expectedTags))
	}
}

// TestMapStructsJsonToStructSlice tests mapping structures where a JSON string is converted to a slice.
func TestMapStructsJsonToStructSlice(t *testing.T) {
	type Tag struct {
		Name string `json:"name"`
	}

	type Src struct {
		Tags string `json:"tags"`
	}

	type Dest struct {
		Tags []Tag `json:"tags"`
	}

	// JSON string representing an array of Tag structs
	jsonTags := `[{"name":"go"},{"name":"programming"},{"name":"developer"}]`
	src := Src{Tags: jsonTags}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error when mapping JSON string to slice: %s", err)
	}

	// Convert the JSON string to a slice of Tag structs for comparison
	var expectedTags []Tag
	if err := json.Unmarshal([]byte(jsonTags), &expectedTags); err != nil {
		t.Fatalf("Failed to unmarshal JSON string to expected tags: %s", err)
	}

	if !reflect.DeepEqual(dest.Tags, expectedTags) {
		t.Errorf("Failed to map JSON string to slice correctly, got: %+v, want: %+v", dest.Tags, expectedTags)
	}
}

// TestMapStructToJson tests mapping a struct to a JSON string.
func TestMapStructToJson(t *testing.T) {
	type Tag struct {
		Name string `json:"name"`
	}

	type Src struct {
		Tag Tag `json:"tag"`
	}

	type Dest struct {
		Tag string `json:"tag"`
	}

	src := Src{Tag: Tag{Name: "developer"}}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error when mapping struct to JSON string: %s", err)
	}

	// Convert the source struct to JSON for comparison
	expectedTag, err := json.Marshal(src.Tag)
	if err != nil {
		t.Fatalf("Failed to marshal source tag to JSON: %s", err)
	}

	if dest.Tag != string(expectedTag) {
		t.Errorf("Failed to map struct to JSON string correctly, got: %s, want: %s", dest.Tag, string(expectedTag))
	}
}

// TestMapJsonToStruct tests mapping a JSON string to a struct.
func TestMapJsonToStruct(t *testing.T) {
	type Tag struct {
		Name string `json:"name"`
	}

	type Src struct {
		Tag string `json:"tag"`
	}

	type Dest struct {
		Tag Tag `json:"tag"`
	}

	// JSON string representing a Tag struct
	jsonTag := `{"name":"developer"}`
	src := Src{Tag: jsonTag}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error when mapping JSON string to struct: %s", err)
	}

	// Convert the JSON string to a Tag struct for comparison
	var expectedTag Tag
	if err := json.Unmarshal([]byte(jsonTag), &expectedTag); err != nil {
		t.Fatalf("Failed to unmarshal JSON string to expected tag: %s", err)
	}

	if !reflect.DeepEqual(dest.Tag, expectedTag) {
		t.Errorf("Failed to map JSON string to struct correctly, got: %+v, want: %+v", dest.Tag, expectedTag)
	}
}

// TestMappingStringPointers checks the mapping functionality for string pointers.
func TestMappingStringPointers(t *testing.T) {
	type Src struct {
		Text string `json:"text"`
	}
	type Dest struct {
		Text *string `json:"text"`
	}

	// Test with non-empty source text
	srcText := "example_text"
	src := Src{Text: srcText}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error during mapping: %s", err)
	}

	if dest.Text == nil || *dest.Text != srcText {
		t.Errorf("Expected dest.Text to be %s, but got %v", srcText, dest.Text)
	}

	// Test with empty source text
	srcEmptyText := ""
	src = Src{Text: srcEmptyText}
	dest = Dest{}

	err = xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Errorf("Unexpected error during mapping with empty text: %s", err)
	}

	if dest.Text == nil || *dest.Text != srcEmptyText {
		t.Errorf("Expected dest.Text to be an empty string, but got %v", dest.Text)
	}
}

// TestMappingFromPointerString checks the mapping functionality from a string pointer in the source struct to a regular string in the destination struct.
func TestMappingFromPointerString(t *testing.T) {
	type Src struct {
		Text *string `json:"text"`
	}
	type Dest struct {
		Text string `json:"text"`
	}

	// Test with non-empty source text
	srcText := "example_text"
	src := Src{Text: &srcText}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping: %s", err)
	}

	if dest.Text != srcText {
		t.Errorf("Expected dest.Text to be %s, but got %v", srcText, dest.Text)
	}

	// Test with empty source text
	emptyText := ""
	src = Src{Text: &emptyText}
	dest = Dest{}

	err = xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping with empty text: %s", err)
	}

	if dest.Text != emptyText {
		t.Errorf("Expected dest.Text to be an empty string, but got %v", dest.Text)
	}

	// Test with nil source text pointer
	src = Src{Text: nil}
	dest = Dest{}

	err = xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping with nil text: %s", err)
	}

	if dest.Text != "" {
		t.Errorf("Expected dest.Text to be an empty string, but got %v", dest.Text)
	}
}

// TestMappingBetweenPointerStrings checks the mapping functionality from a string pointer in the source struct to a string pointer in the destination struct.
func TestMappingBetweenPointerStrings(t *testing.T) {
	type Src struct {
		Text *string `json:"text"`
	}
	type Dest struct {
		Text *string `json:"text"`
	}

	// Test with non-empty source text
	srcText := "example_text"
	src := Src{Text: &srcText}
	dest := Dest{}

	err := xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping: %s", err)
	}

	if dest.Text == nil || *dest.Text != srcText {
		t.Errorf("Expected dest.Text to be %s, but got %v", srcText, dest.Text)
	}

	// Test with empty source text
	emptyText := ""
	src = Src{Text: &emptyText}
	dest = Dest{}

	err = xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping with empty text: %s", err)
	}

	if dest.Text == nil || *dest.Text != emptyText {
		t.Errorf("Expected dest.Text to be an empty string, but got %v", dest.Text)
	}

	// Test with nil source text pointer
	src = Src{Text: nil}
	dest = Dest{}

	err = xmapper.MapStructs(&src, &dest)
	if err != nil {
		t.Fatalf("Unexpected error during mapping with nil text: %s", err)
	}

	if dest.Text != nil {
		t.Errorf("Expected dest.Text to be nil, but got %v", dest.Text)
	}
}

// Define a simple struct for testing the MapJsonStruct function
type TestProfile struct {
	Name string `json:"name" transformers:"toUpperCase"`
	Age  int    `json:"age"`
}

// TestMapJsonStructSuccess tests the successful unmarshalling of a JSON string into the struct
func TestMapJsonStructSuccess(t *testing.T) {
	jsonStr := `{"name":"John Doe","age":30}`
	var profile TestProfile

	err := xmapper.MapJsonStruct(jsonStr, &profile)
	if err != nil {
		t.Errorf("MapJsonStruct failed: %s", err)
	}

	if profile.Name != "JOHN DOE" || profile.Age != 30 {
		t.Errorf("MapJsonStruct did not correctly unmarshal the JSON. Got %v", profile)
	}
}

// TestMapJsonStructInvalidJSON tests the function with invalid JSON input
func TestMapJsonStructInvalidJSON(t *testing.T) {
	jsonStr := `{"name":"John Doe", "age": "thirty"}`
	var profile TestProfile

	err := xmapper.MapJsonStruct(jsonStr, &profile)
	if err == nil {
		t.Errorf("MapJsonStruct should have failed on invalid JSON but did not")
	}
}

// TestMapJsonStructNonPointer tests the function with a non-pointer target
func TestMapJsonStructNonPointer(t *testing.T) {
	jsonStr := `{"name":"John Doe","age":30}`
	profile := TestProfile{}

	err := xmapper.MapJsonStruct(jsonStr, profile)
	if err == nil || err.Error() != "target must be a pointer to a struct" {
		t.Errorf("MapJsonStruct should have failed with a non-pointer target but did not")
	}
}

// TestValidateStructWithInvalidData checks the validation functionality for an invalid email format.
func TestValidateStructWithInvalidData(t *testing.T) {
	// Define the source structure with validation tags.
	type Src struct {
		Email string `json:"email" validators:"email"`
	}

	// Create an instance of Src with invalid email.
	src := Src{Email: "not_a_valid_email"}

	// Call ValidateStruct to check validation.
	err := xmapper.ValidateStruct(&src)

	// Check if error is returned for invalid email.
	if err == nil {
		t.Errorf("Expected error for invalid email, got nil")
	} else {
		t.Logf("Received expected error: %v", err)
	}
}

// TestValidateStructWithValidDataAndTransformer checks the validation functionality for a single struct with transformer.
func TestValidateStructWithValidDataAndTransformer(t *testing.T) {
	// Define the source structure with validation tags.
	type Src struct {
		Name string `json:"name" validators:"required" transformers:"uppercase"`
	}

	// Create an instance of Src with invalid email.
	src := Src{Name: "test"}

	// Call ValidateStruct to check validation.
	err := xmapper.ValidateStruct(&src)

	// Check if error is returned for invalid email.
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if src.Name != "TEST" {
		t.Errorf("Failed to apply transformations correctly, got: %+v", src)
	}
}
