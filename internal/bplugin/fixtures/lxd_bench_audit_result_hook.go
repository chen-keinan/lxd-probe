package main

import (
	"fmt"
	"github.com/chen-keinan/lxd-probe/pkg/models"
)

//LxdBenchAuditResultHook this plugin method accept lxd audit bench results
//event include test data , description , audit, remediation and result
//nolint
func LxdBenchAuditResultHook(lxdAuditResults models.LxdAuditResults) error {
	fmt.Println("this is LxdBenchAuditResultHook plugin")
	return nil
}
