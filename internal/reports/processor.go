package reports

import (
	"github.com/chen-keinan/lxd-probe/internal/models"
	"github.com/gosuri/uitable"
)

//GenerateAuditReport generate failed audit report
func GenerateAuditReport(adtsReport []*models.AuditBench) *uitable.Table {
	table := uitable.New()
	for _, failedAudit := range adtsReport {
		table.MaxColWidth = 100
		status := "Failed"
		if failedAudit.NonApplicable {
			status = "Warn"
		}
		table.Wrap = true // wrap columns
		table.AddRow("--------------", "-------------------------------------------------------------------------------------------")
		table.AddRow("Status:", status)
		table.AddRow("Name:", failedAudit.Name)
		table.AddRow("Description:", failedAudit.Description)
		table.AddRow("Audit:", failedAudit.AuditCommand)
		table.AddRow("Remediation:", failedAudit.Remediation)
		table.AddRow("References:", failedAudit.References)
		table.AddRow("") // blank
	}
	return table
}
