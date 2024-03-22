package mack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinalValue_Floats_Simple(t *testing.T) {
	// zero
	found := finalValue("0.0")
	assert.Equal(t, "0.0", found)
	// value
	found = finalValue("10.01")
	assert.Equal(t, "10.01", found)
	// trailing zeroes
	found = finalValue("10.0100")
	assert.Equal(t, "10.0100", found)
	// strings are not swept up
	found = finalValue("\"10.0100\"")
	assert.Equal(t, "\"10.0100\"", found)
	found = finalValue("\"J. K. Rowling\"")
	assert.Equal(t, "\"J. K. Rowling\"", found)
	// ints are not floats
	found = finalValue("10")
	assert.Equal(t, "10", found)
	// stupid super exponent numbers
	found = finalValue("10.0100E+3")
	assert.Equal(t, "10.0100E+3", found)
}
