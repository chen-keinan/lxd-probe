package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_GetSpecificTestsToExecute test
func Test_GetSpecificTestsToExecute(t *testing.T) {
	l := GetAuditTestsList("i", "i=1.2.3,1.4.5")
	assert.Equal(t, l[0], "1.2.3")
	assert.Equal(t, l[1], "1.4.5")
	l = GetAuditTestsList("e", "")
	assert.Equal(t, l[0], "")
}

//Test_RemoveNewLineSuffix test
func Test_RemoveNewLineSuffix(t *testing.T) {
	s := RemoveNewLineSuffix("abc\n")
	assert.Equal(t, s, "abc")
	s = RemoveNewLineSuffix("abc\n134")
	assert.Equal(t, s, "abc\n134")
	s = RemoveNewLineSuffix("abc")
	assert.Equal(t, s, "abc")
}

//Test_AddNewLineToNonEmptyStr test
func Test_AddNewLineToNonEmptyStr(t *testing.T) {
	k := AddNewLineToNonEmptyStr("abc")
	assert.Equal(t, k, "abc\n")
	k = AddNewLineToNonEmptyStr("\n")
	assert.Equal(t, k, "\n")
	k = AddNewLineToNonEmptyStr("abc\n")
	assert.Equal(t, k, "abc\n")
}
