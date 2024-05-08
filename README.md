# xMapper 🗺️ - Dynamic Struct Transformer

Welcome to `xMapper`, where your Go structs gain superpowers! 🚀 Ever tired of manually copying fields from one struct to another? Forget about the hassle! `xMapper` automates the mapping of structs and dynamically transforms data as it flows from source to destination. It's not just about mapping; it's about transforming data effortlessly with power and style!

## Features

- **Automatic Struct Mapping**: Automate the boring stuff! Map fields between structs without writing boilerplate code.
- **Dynamic Data Transformation**: Apply transformations to your data dynamically as it's being mapped. Upper case, add suffixes, manipulate data on the go!
- **Extensible and Customizable**: Easily extend `xMapper` by adding your own custom transformation functions.
- **Error Handling**: Robust error handling to let you know exactly what went wrong during the mapping process.

## Getting Started

Here’s how to get started with `xMapper`:

### Installation

To start using `xMapper` in your Go project, simply use the `go get` command to retrieve the package:

```bash
go get github.com/yourusername/xmapper
```

Ensure your environment is set up with Go modules (Go 1.11+ required), and this command will manage everything for you, fetching the latest version of `xMapper` and adding it to your project's dependencies.

### Usage

1. **Define Your Transformers**: Create functions that match the `TransformerFunc` signature:

    ```go
    func toUpperCase(input interface{}) interface{} {
        str, ok := input.(string)
        if ok {
            return strings.ToUpper(str)
        }
        return input
    }
    ```

2. **Register Your Transformers**: Before you map your structs, make sure to register your transformers:

    ```go
    xmapper.RegisterTransformer("toUpperCase", toUpperCase)
    ```

3. **Map Your Structs**: Now let the magic happen:

    ```go
    type Source struct {
        Name string `json:"name" transformer:"toUpperCase"`
    }

    type Destination struct {
        Name string `json:"name"`
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
func toUpperCase(input interface{}) interface{} {
    str, ok := input.(string)
    if ok {
        return strings.ToUpper(str)
    }
    return input
}

// Adds an exclamation mark at the end of a string
func addExclamation(input interface{}) interface{} {
    str, ok := input.(string)
    if ok {
        return str + "!"
    }
    return input
}

// Repeats the string twice, separated by a space
func repeatTwice(input interface{}) interface{} {
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
func init() {
    xmapper.RegisterTransformer("toUpperCase", toUpperCase)
    xmapper.RegisterTransformer("addExclamation", addExclamation)
    xmapper.RegisterTransformer("repeatTwice", repeatTwice)
}
```

### Step 3: Set Up Your Structs

Define your source and destination structs. Use the `transformer` tag to specify multiple transformers separated by commas. The transformers will be applied in the order they are listed:

```go
type Source struct {
    Greeting string `json:"greeting" transformer:"toUpperCase,addExclamation,repeatTwice"`
}

type Destination struct {
    Greeting string `json:"greeting"`
}

func main() {
    src := Source{Greeting: "hello"}
    dest := Destination{}

    if err := xmapper.MapStructs(&src, &dest); err != nil {
        fmt.Println("Error mapping structs:", err)
    } else {
        fmt.Println("Mapped Greeting:", dest.Greeting)
        // Output should be: "HELLO! HELLO!"
    }
}
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
