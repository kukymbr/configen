package gentype_test

import (
	"testing"

	"github.com/kukymbr/configen/internal/generator/gentype"
	"github.com/stretchr/testify/assert"
)

func TestPrefixCommentLines(t *testing.T) {
	tests := []struct {
		Name     string
		Comment  string
		Prefix   string
		Expected string
	}{
		{Name: "empty", Comment: " \n ", Prefix: "# ", Expected: ""},
		{Name: "single line", Comment: "Port to listen.\n", Prefix: "# ", Expected: "# Port to listen."},
		{Name: "multiline", Comment: "First line.\nSecond line.\n", Prefix: "// ", Expected: "// First line.\n// Second line."},
		{Name: "empty line inside", Comment: "First.\n\nThird.", Prefix: "// ", Expected: "// First.\n//\n// Third."},
		{Name: "trailing spaces", Comment: "First.  \nSecond.\t", Prefix: "# ", Expected: "# First.\n# Second."},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, gentype.PrefixCommentLines(test.Comment, test.Prefix))
		})
	}
}
