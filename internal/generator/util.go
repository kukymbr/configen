package generator

import (
	"errors"
	"fmt"
	"regexp"
)

var pxIdentifier = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9_]*$`)

func validateIdentifier(name string) error {
	if len(name) == 0 {
		return errors.New("identifier cannot be empty")
	}

	if !pxIdentifier.MatchString(name) {
		return fmt.Errorf("'%s' is not a valid identifier", name)
	}

	return nil
}
