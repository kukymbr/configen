package gentype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPublicName(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"API", "APIProvider"},
		{"APIConfig", "APIConfigProvider"},
		{"api", "API"},
		{"apiConfig", "APIConfig"},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			res := ToPublicName(test.Input)

			assert.Equal(t, test.Expected, res)
		})
	}
}

func TestNameToParts(t *testing.T) {
	tests := []struct {
		Input    string
		Expected []string
	}{
		{Input: "API", Expected: []string{"api"}},
		{Input: "APIConfig", Expected: []string{"api", "config"}},
		{Input: "apiConfig", Expected: []string{"api", "config"}},
		{Input: "apiCONFIG", Expected: []string{"api", "config"}},
		{Input: "name", Expected: []string{"name"}},
		{Input: "test-name", Expected: []string{"test", "name"}},
		{Input: "test_name", Expected: []string{"test", "name"}},
		{Input: "test name", Expected: []string{"test", "name"}},
		{Input: "TestName", Expected: []string{"test", "name"}},
		{Input: "Test_Name", Expected: []string{"test", "name"}},
		{Input: "Test Name", Expected: []string{"test", "name"}},
		{Input: "test.name", Expected: []string{"test", "name"}},
		{Input: "test, name", Expected: []string{"test", "name"}},
		{Input: "test -/- name", Expected: []string{"test", "name"}},
		{Input: "test name.", Expected: []string{"test", "name"}},
		{Input: "testNAME", Expected: []string{"test", "name"}},
		{Input: "__TEST_name.", Expected: []string{"test", "name"}},
		{Input: "TEST1name", Expected: []string{"test", "1", "name"}},
		{Input: "TEST123name", Expected: []string{"test", "123", "name"}},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			parts := nameToWords(test.Input)

			assert.Equal(t, test.Expected, parts)
		})
	}
}
