package filters

import (
	"github.com/chen-keinan/lxd-probe/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_IncludeTestPredict text
func Test_IncludeTestPredict(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}}
	abp := IncludeAudit(ab, "1.2.1")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.1 abc")
	assert.True(t, len(abp.AuditTests) == 1)
}

//Test_IncludeTestPredicateNoValidArg text
func Test_IncludeTestPredicateNoValidArg(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}}
	abp := IncludeAudit(ab, "1.2.5")
	assert.True(t, len(abp.AuditTests) == 0)
}

//Test_ExcludeTestPredict text
func Test_ExcludeTestPredict(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}}
	abp := ExcludeAudit(ab, "1.2.1")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.3 eft")
	assert.True(t, len(abp.AuditTests) == 1)
}

//Test_ExcludeTestPredicateNoValidArg text
func Test_ExcludeTestPredicateNoValidArg(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}}
	abp := ExcludeAudit(ab, "1.2.5")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.1 abc")
	assert.Equal(t, abp.AuditTests[1].Name, "1.2.3 eft")
	assert.True(t, len(abp.AuditTests) == 2)
}

//Test_BasicPredicate text
func Test_BasicPredicate(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}}
	abp := Basic(ab, "")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.1 abc")
	assert.Equal(t, abp.AuditTests[1].Name, "1.2.3 eft")
	assert.True(t, len(abp.AuditTests) == 2)
}

//Test_NodePredicateMaster text
func Test_NodePredicateMaster(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc", ProfileApplicability: "Master"}, {Name: "1.2.3 eft", ProfileApplicability: "Worker"}}}
	abp := NodeAudit(ab, "n=master")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.1 abc")
	assert.True(t, len(abp.AuditTests) == 1)
}

//Test_NodePredicateWorker text
func Test_NodePredicateWorker(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc", ProfileApplicability: "Master"}, {Name: "1.2.3 eft", ProfileApplicability: "Worker"}}}
	abp := NodeAudit(ab, "n=worker")
	assert.Equal(t, abp.AuditTests[0].Name, "1.2.3 eft")
	assert.True(t, len(abp.AuditTests) == 1)
}

//Test_NodePredicateNone text
func Test_NodePredicateNone(t *testing.T) {
	ab := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 abc", ProfileApplicability: "Master"}, {Name: "1.2.3 eft", ProfileApplicability: "Worker"}}}
	abp := NodeAudit(ab, "n=abd")
	assert.Equal(t, len(abp.AuditTests), 0)
}
