package commands

import (
	"github.com/chen-keinan/lxd-probe/internal/logger"
	"github.com/chen-keinan/lxd-probe/internal/mocks"
	"github.com/chen-keinan/lxd-probe/internal/models"
	"github.com/chen-keinan/lxd-probe/internal/shell"
	"github.com/chen-keinan/lxd-probe/internal/startup"
	"github.com/chen-keinan/lxd-probe/pkg/filters"
	m2 "github.com/chen-keinan/lxd-probe/pkg/models"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

//Test_AddFailedMessages text
func Test_AddFailedMessages(t *testing.T) {
	atb1 := &models.AuditBench{TestSucceed: false}
	afm := AddFailedMessages(atb1, 1)
	assert.True(t, len(afm) == 1)
	atb2 := &models.AuditBench{TestSucceed: true}
	afm = AddFailedMessages(atb2, 0)
	assert.True(t, len(afm) == 0)
}

//Test_isArgsExist
func Test_isArgsExist(t *testing.T) {
	args := []string{"aaa", "bbb"}
	exist := isArgsExist(args, "aaa")
	assert.True(t, exist)
	exist = isArgsExist(args, "ccc")
	assert.False(t, exist)
}

//Test_isArgsExist
func Test_GetProcessingFunction(t *testing.T) {
	args := []string{"r"}
	a := GetResultProcessingFunction(args)
	name := GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func5"))
	args = []string{"report"}
	a = GetResultProcessingFunction(args)
	name = GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func5"))
	args = []string{"c"}
	a = GetResultProcessingFunction(args)
	name = GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func4"))
	args = []string{"classic"}
	a = GetResultProcessingFunction(args)
	name = GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func4"))
	args = []string{}
	a = GetResultProcessingFunction(args)
	name = GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func4"))
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//Test_getSpecificTestsToExecute test
func Test_getSpecificTestsToExecute(t *testing.T) {
	test := utils.GetAuditTestsList("i", "i=1.2.4,1.2.5")
	assert.Equal(t, test[0], "1.2.4")
	assert.Equal(t, test[1], "1.2.5")
}

//Test_LoadAuditTest test
func Test_LoadAuditTest(t *testing.T) {
	fm := utils.NewKFolder()
	folder, err2 := utils.GetBenchmarkFolder("lxd", "v1.0.0", fm)
	assert.NoError(t, err2)
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = utils.CreateHomeFolderIfNotExist(fm)
	if err != nil {
		t.Fatal(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist("lxd", "v1.0.0", fm)
	if err != nil {
		t.Fatal(err)
	}
	bFiles, err := startup.GenerateLxdBenchmarkFiles()
	if err != nil {
		t.Fatal(err)
	}
	err = startup.SaveBenchmarkFilesIfNotExist("lxd", "v1.0.0", bFiles)
	if err != nil {
		t.Fatal(err)
	}
	at := NewFileLoader().LoadAuditTests(bFiles)
	assert.True(t, len(at) != 0)
	assert.True(t, strings.Contains(at[0].AuditTests[0].Name, "1.1.1"))
}

//Test_LoadGkeAuditTest test
func Test_LoadGkeAuditTest(t *testing.T) {
	fm := utils.NewKFolder()
	folder, err2 := utils.GetBenchmarkFolder("gke", "v1.1.0", fm)
	assert.NoError(t, err2)
	err := os.RemoveAll(folder)
	if err != nil {
		t.Fatal(err)
	}
	err = utils.CreateHomeFolderIfNotExist(fm)
	if err != nil {
		t.Fatal(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist("gke", "v1.1.0", fm)
	if err != nil {
		t.Fatal(err)
	}
	bFiles, err := startup.GenerateLxdBenchmarkFiles()
	if err != nil {
		t.Fatal(err)
	}
	err = startup.SaveBenchmarkFilesIfNotExist("gke", "v1.1.0", bFiles)
	if err != nil {
		t.Fatal(err)
	}
	at := NewFileLoader().LoadAuditTests(bFiles)
	assert.True(t, len(at) != 0)
	assert.True(t, strings.Contains(at[0].AuditTests[0].Name, "1.1.1"))
}

//Test_FilterAuditTests test
func Test_FilterAuditTests(t *testing.T) {
	at := &models.SubCategory{AuditTests: []*models.AuditBench{{Name: "1.2.1 aaa"}, {Name: "2.2.2"}}}
	fab := FilterAuditTests([]filters.Predicate{filters.IncludeAudit}, []string{"1.2.1"}, at)
	assert.Equal(t, fab.AuditTests[0].Name, "1.2.1 aaa")
	assert.True(t, len(fab.AuditTests) == 1)
}

//Test_buildPredicateChain test
func Test_buildPredicateChain(t *testing.T) {
	fab := buildPredicateChain([]string{"a", "i=1.2.1"})
	assert.True(t, len(fab) == 2)
	fab = buildPredicateChain([]string{"a"})
	assert.True(t, len(fab) == 1)
	fab = buildPredicateChain([]string{"i=1.2.1"})
	assert.True(t, len(fab) == 1)
}

//Test_buildPredicateChainParams test
func Test_buildPredicateChainParams(t *testing.T) {
	p := buildPredicateChainParams([]string{"a", "i=1.2.1"})
	assert.True(t, len(p) == 2)
	assert.Equal(t, p[0], "a")
	assert.Equal(t, p[1], "i=1.2.1")
}

//Test_buildPredicateChainParamsExcludeNode test
func Test_buildPredicateChainExcludeNode(t *testing.T) {
	p := buildPredicateChainParams([]string{"a", "n=master"})
	assert.True(t, len(p) == 2)
	assert.Equal(t, p[0], "a")
	assert.Equal(t, p[1], "n=master")
	p = buildPredicateChainParams([]string{"a", "e=1.2.1"})
	assert.True(t, len(p) == 2)
	assert.Equal(t, p[0], "a")
	assert.Equal(t, p[1], "e=1.2.1")
}

func Test_filteredAuditBenchTests(t *testing.T) {
	asc := []*models.SubCategory{{AuditTests: []*models.AuditBench{{Name: "1.1.0 bbb"}}}}
	fp := []filters.Predicate{filters.IncludeAudit, filters.ExcludeAudit}
	st := []string{"i=1.1.0", "e=1.1.0"}
	fr := filteredAuditBenchTests(asc, fp, st)
	assert.True(t, len(fr) == 0)
}

//Test_executeTests test
func Test_executeTests(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa", "bbb"}
	ab.EvalExpr = "'$0' == ''; && '$1' == '';"
	ab.CommandParams = map[int][]string{}
	ab.CmdExprBuilder = utils.UpdateCmdExprParam
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "\n\n\n\n\n"}, nil).Times(1)
	executor.EXPECT().Exec("bbb").Return(&shell.CommandResult{Stdout: "default-token-ppzx7\n\n\n\n\n"}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.LxdAuditResults)
	kb := LxdAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	sc := []*models.SubCategory{{AuditTests: []*models.AuditBench{ab}}}
	executeTests(sc, kb.runAuditTest, logger.GetLog())
	assert.False(t, ab.TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

func Test_printTestResults(t *testing.T) {
	ab := make([]*models.AuditBench, 0)
	ats := &models.AuditBench{Name: "bbb", TestSucceed: true}
	atf := &models.AuditBench{Name: "ccc", TestSucceed: false}
	ata := &models.AuditBench{Name: "ddd", NonApplicable: true}
	ab = append(ab, ats)
	ab = append(ab, atf)
	ab = append(ab, ata)
	tr := printTestResults(ab, tablewriter.NewWriter(os.Stdout), "aaa")
	assert.Equal(t, tr.Warn, 1)
	assert.Equal(t, tr.Pass, 1)
	assert.Equal(t, tr.Fail, 1)
}
