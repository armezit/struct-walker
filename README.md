# Golang Struct Walker

A simple Go library for traversing nested structs.
It largely relies on [reflect](https://golang.org/pkg/reflect/) methods, so be careful about performance.

## Installing

```sh
go get -u github.com/armezit/struct-walker
```

## Usage

```go
package main

import (
	"fmt"
	walker "github.com/armezit/struct-walker"
	"reflect"
	"strings"
)

func main() {
	type Bar struct {
		AA string
		BB int
	}

	type Foo struct {
		A string
		B struct {
			B1 float64
			B2 float64
		}
		C struct {
			C1 uint
			C2 uint
		}
		D map[string]Bar
		E []Bar
	}

	// foo initialization
	foo := Foo{
		// ...
	}

	visitor := func(value reflect.Value, branch []interface{}, path []string, field *reflect.StructField) {
		k := strings.Join(path, ".")
		v := value.Interface()
		fmt.Printf("%s: %v", k, v)
	}
	walker.Walk(foo, visitor)
}
```
