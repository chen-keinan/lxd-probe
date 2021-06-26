package ui

import (
	"github.com/cheggaaa/pb"
	"github.com/chen-keinan/lxd-probe/internal/logger"
	"github.com/chen-keinan/lxd-probe/internal/models"
	"time"
)

// OutputGenerator for  audit results
type OutputGenerator func(at []*models.SubCategory, log *logger.LdxProbeLogger)

//PrintOutput print audit test result to console
func PrintOutput(auditTests []*models.SubCategory, outputGenerator OutputGenerator, log *logger.LdxProbeLogger) {
	log.Console(auditResult)
	outputGenerator(auditTests, log)
}

//ShowProgressBar execute audit test and show progress bar
func ShowProgressBar(a *models.SubCategory, execTestFunc func(ad *models.AuditBench) []*models.AuditBench, log *logger.LdxProbeLogger) *models.SubCategory {
	if len(a.AuditTests) == 0 {
		return a
	}
	completedTest := make([]*models.AuditBench, 0)
	log.Console(a.Name)
	bar := pb.StartNew(len(a.AuditTests))
	for _, test := range a.AuditTests {
		ar := execTestFunc(test)
		completedTest = append(completedTest, ar...)

		bar.Increment()
		time.Sleep(time.Millisecond * 20)
	}
	bar.Finish()
	return &models.SubCategory{Name: a.Name, AuditTests: completedTest}
}
