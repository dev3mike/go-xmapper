
# xMapper 🔄 - Dynamic Struct Validator and Mapper

  

Welcome to `xMapper`, where your Go structs gain superpowers! 🚀 Ever tired of manually copying fields from one struct to another? Forget about the hassle! `xMapper` automates the mapping of structs and dynamically transforms data as it flows from source to destination. It's not just about mapping; it's about transforming data effortlessly with power and style!

  

## Features

  

-  **Automatic Struct Mapping**: Automate the boring stuff! Map fields between structs without writing boilerplate code.

-  **Dynamic Data Transformation**: Apply transformations to your data dynamically as it's being mapped. Upper case, add suffixes, manipulate data on the go!

-  **Extensible and Customizable**: Easily extend `xMapper` by adding your own custom transformation functions.

-  **Data Validation**: Ensure data integrity with custom validators that check data before transformation.

-  **Default Validation Rules**:   Check input using built-in validators for common requirements like email format, length, or specific starting/ending characters without writing custom code.

-  **Error Handling**: Robust error handling to let you know exactly what went wrong during the mapping process.

  

## Getting Started

  

Here’s how to get started with `xMapper`:

  

### Installation

  

To start using `xMapper` in your Go project, simply use the `go get` command to retrieve the package:

  

```go
go get github.com/dev3mike/go-xmapper
```

  

Ensure your environment is set up with Go modules (Go 1.11+ required), and this command will manage everything for you, fetching the latest version of `xMapper` and adding it to your project's dependencies.

  

### Usage

  

1.  **Define Your Transformers**: Create functions that match the `TransformerFunc` signature:

  

```go
func  toUpperCase(input interface{}) interface{} {
str, ok := input.(string)
if ok {
	return strings.ToUpper(str)
}
	return input
}

```

  

2.  **Register Your Transformers**: Before you map your structs, make sure to register your transformers:

  

```go

xmapper.RegisterTransformer("toUpperCase", toUpperCase)

```

  

3.  **Map Your Structs**: Now let the magic happen:

  

```go

type  Source  struct {

Name string  `json:"name" transformer:"toUpperCase"`

}

  

type  Destination  struct {

Name string  `json:"name"`

}

  

src := Source{Name: "frodo"}

dest := Destination{}

  

err := xmapper.MapStructs(&src, &dest)

if err != nil {

fmt.Println("Oops! Something went wrong:", err)

}

```

  

### Example with Error Handling

  

Want to ensure your transformers are set up correctly? Here’s how to handle errors:

  

```go

err := xmapper.MapStructs(&src, &dest)

if err != nil {

fmt.Println("Failed to map structs:", err)

}

```

  

## Using Multiple Transformers

  

`xMapper` allows you to apply multiple transformations to a single field in sequence, which can be extremely powerful for complex data manipulation. This section guides you through setting up and using multiple transformers on a single struct field.

  

### Step 1: Define Your Transformers

  

First, define each transformer function. Each function should match the `TransformerFunc` signature. Here are examples of three simple transformers:

  

```go

// Converts a string to uppercase

func  toUpperCase(input interface{}) interface{} {

str, ok := input.(string)

if ok {

return strings.ToUpper(str)

}

return input

}

  

// Adds an exclamation mark at the end of a string

func  addExclamation(input interface{}) interface{} {

str, ok := input.(string)

if ok {

return str + "!"

}

return input

}

  

// Repeats the string twice, separated by a space

func  repeatTwice(input interface{}) interface{} {

str, ok := input.(string)

if ok {

return str + " " + str

}

return input

}

```

  

### Step 2: Register Your Transformers

  

Register each transformer with `xMapper` before you attempt to map your structs:

  

```go

func  init() {

xmapper.RegisterTransformer("toUpperCase", toUpperCase)

xmapper.RegisterTransformer("addExclamation", addExclamation)

xmapper.RegisterTransformer("repeatTwice", repeatTwice)

}

```

  

### Step 3: Set Up Your Structs

  

Define your source and destination structs. Use the `transformer` tag to specify multiple transformers separated by commas. The transformers will be applied in the order they are listed:

  

```go

type  Source  struct {

Greeting string  `json:"greeting" transformer:"toUpperCase,addExclamation,repeatTwice"`

}

  

type  Destination  struct {

Greeting string  `json:"greeting"`

}

  

func  main() {

src := Source{Greeting: "hello"}

dest := Destination{}

  

if  err := xmapper.MapStructs(&src, &dest); err != nil {

fmt.Println("Error mapping structs:", err)

} else {

fmt.Println("Mapped Greeting:", dest.Greeting)

// Output should be: "HELLO! HELLO!"

}

}

```

  

## Using Validators

  

Validators in `xMapper` ensure your data meets specific criteria before it's transformed and mapped to the destination struct. Validators can prevent invalid data from being processed and provide descriptive error messages if data validation fails.

You can use built-in validator as follows:
| Validator Name    | Description                                                          |
|-------------------|----------------------------------------------------------------------|
| `required`        | Checks if the input is not empty.                                    |
| `email`           | Validates that the input is a valid email address.                   |
| `phone`           | Checks if the input is a valid international phone number.           |
| `strongPassword`  | Requires at least 8 characters, including upper, lower, digit, and special character. |
| `date`            | Validates that the input matches the YYYY-MM-DD date format.         |
| `time`            | Validates that the input matches the HH:MM:SS time format.           |
| `datetime`        | Validates date and time with timezone in YYYY-MM-DD HH:MM:SS format. |
| `url`             | Checks if the input is a valid URL.                                  |
| `ip`              | Validates that the input is a valid IP address.                      |
| `minLength`       | Checks if the input has at least a specified minimum length.         |
| `maxLength`       | Ensures the input does not exceed a specified maximum length.        |
| `gt`              | Validates that a number is greater than a specified value.           |
| `lt`              | Validates that a number is less than a specified value.              |
| `gte`             | Checks if a number is greater than or equal to a specified value.    |
| `lte`             | Checks if a number is less than or equal to a specified value.       |
| `range`           | Validates that a number falls within a specified range.              |
| `enum`            | Checks if the input matches one of a list of predefined values.      |
| `boolean`         | Validates that the input is a boolean value.                         |
| `contains`        | Checks if the input contains one of the specified substrings.        |
| `notContains`     | Ensures the input does not contain any of the specified substrings.  |
| `startsWith`      | Validates that the input starts with a specified substring.          |
| `endsWith`        | Checks if the input ends with a specified substring.                 |

**How to use built-in validators:**

    type  User  struct {
	    Name string  `json:"name" validator:"required,minLength:3,maxLength:100"`
	    Email string  `json:"email" validator:"email,maxLength:255"`
	    Username string  `json:"username" validator:"minLength:6,maxLength:27"`
	    Password string  `json:"password" validator:"strongPassword"`
	    Type string  `json:"type" validator:"enum:buyer,seller"`
	}

  

### Step 1: Define Your Validators

  

First, define each validator function. Each function must match the `ValidatorFunc` signature, which typically returns `true` if the validation passes and `false` otherwise:

```go

// Ensures the input string is not empty
// The second argument is always required. 
// It's used to receive specific settings for the validation. 
// For instance, for `minLength:10`, 
// the number `10` is the argument that sets the minimum length.
func  isNotEmpty(input interface{}, _ string) bool {
	str, ok := input.(string)
	return ok && str != ""
}

// Checks that the input string does not contain spaces
func  doesNotContainSpaces(input interface{}, _ string) bool {
	str, ok := input.(string)
	return ok && !strings.Contains(str, " ")
}

// Checks if the input string's length is at least the specified minimum
func  MinLengthValidator(input interface{}, length string) error {
	str, ok  := input.(string)
	if  !ok {
		return fmt.Errorf("failed to map the input to a string")
	}
	minLength, err  := strconv.Atoi(length)
	if err !=  nil {
		return fmt.Errorf("failed to convert length to integer")
	}
	if  len(str) < minLength {
		return fmt.Errorf("input does not meet the minimum length requirement, minimum length is %s", length)
	}
	return  nil
}
```
  

### Step 2: Register Your Validators
Register each validator with `xMapper` just like you register transformers. This registration links the validator name with its corresponding function:


```go

func  init() {
	xmapper.RegisterValidator("isNotEmpty", isNotEmpty)
	xmapper.RegisterValidator("doesNotContainSpaces", doesNotContainSpaces)
	xmapper.RegisterValidator("customMinLength", MinLengthValidator)
}

```

  

### Step 3: Apply Validators to Your Structs

  

When defining your structs, use the `validator` tag to assign validators to struct fields. Multiple validators can be applied to a single field and are separated by commas. Validators will be executed in the order they are listed:

  

```go
type  Person  struct {
	FirstName string  `json:"firstName" validator:"isNotEmpty,doesNotContainSpaces"`
	LastName string  `json:"lastName" validator:"isNotEmpty,customMinLength:4"`
}

type  Destination  struct {
	FirstName string  `json:"firstName"`
	LastName string  `json:"lastName"`
}
```

  

### Step 4: Handle Validation Errors

  

When mapping your structs, handle any validation errors that might arise. If a validation fails, `xMapper` will return an error indicating which field and validator failed:

  

```go
src := Person{FirstName: "John", LastName: "Doe"}
dest := Destination{}

err := xmapper.MapStructs(&src, &dest)
if err != nil {
	fmt.Println("Validation error:", err)
	return
}

fmt.Println("Data successfully validated and mapped:", dest)
```

  

## Contributing

  

Got a cool idea for a new feature? Found a bug? We love contributions!

  

1. Fork the repo.
2. Create a new branch (`git checkout -b cool-new-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin cool-new-feature`).
5. Create a new Pull Request.

  

## License

  

Distributed under the MIT License. See `LICENSE` for more information.

  

---

  

Feel the power of easy struct transformations and focus on what really matters in your Go applications. Try `xMapper` today and say goodbye to boilerplate code! 🎉