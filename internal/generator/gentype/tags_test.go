package gentype_test

import (
	"testing"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/stretchr/testify/assert"
)

func TestValueSeparators(t *testing.T) {
	tests := []struct {
		Name           string
		Tag            string
		ExpectedItem   string
		ExpectedKeyVal string
	}{
		{Name: "defaults", Tag: `env:"HOSTS"`, ExpectedItem: ",", ExpectedKeyVal: ":"},
		{Name: "empty tag", Tag: "", ExpectedItem: ",", ExpectedKeyVal: ":"},
		{Name: "custom", Tag: `env:"DBS" envSeparator:";" envKeyValSeparator:"="`, ExpectedItem: ";", ExpectedKeyVal: "="},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			item, keyVal := gentype.ValueSeparators(test.Tag)

			assert.Equal(t, test.ExpectedItem, item)
			assert.Equal(t, test.ExpectedKeyVal, keyVal)
		})
	}
}
