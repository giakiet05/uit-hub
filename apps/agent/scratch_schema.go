package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

type Input struct {
	ID string `json:"id"`
}

func main() {
	schema, _ := jsonschema.ForType(reflect.TypeFor[Input](), nil)
	schema.Extra = map[string]any{"concurrencySafe": true}
	b, _ := json.Marshal(schema)
	fmt.Println(string(b))
}
