package utils

import (
	"fmt"
	"strings"
)

//GetAuditTestsList return processing function by specificTests
func GetAuditTestsList(key, arg string) []string {
	values := strings.ReplaceAll(arg, fmt.Sprintf("%s=", key), "")
	return strings.Split(strings.ToLower(values), ",")
}

//RemoveNewLineSuffix remove new line from suffix
func RemoveNewLineSuffix(str string) string {
	i := len(str)
	if len(str) > 0 && str[i-1:i] == "\n" {
		return str[0 : i-1]
	}
	return str
}

//AddNewLineToNonEmptyStr add new line to non empty string
func AddNewLineToNonEmptyStr(str string) string {
	if !strings.HasSuffix(str, "\n") {
		return fmt.Sprintf("%s\n", str)
	}
	return str
}
