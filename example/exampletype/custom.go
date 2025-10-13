package exampletype

import (
	"fmt"
	"strings"
)

type KeyVal struct {
	Key   string
	Value string
}

func (kv *KeyVal) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), "=")

	if len(parts) != 2 {
		return fmt.Errorf("invalid key value pair: %q", text)
	}

	*kv = KeyVal{
		Key:   parts[0],
		Value: parts[1],
	}

	return nil
}
