package commands

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/lxd-probe/internal/common"
	"github.com/chen-keinan/lxd-probe/internal/logger"
	"github.com/chen-keinan/lxd-probe/internal/models"
	"github.com/chen-keinan/lxd-probe/internal/reports"
	"github.com/chen-keinan/lxd-probe/internal/shell"
	"github.com/chen-keinan/lxd-probe/internal/startup"
	"github.com/chen-keinan/lxd-probe/pkg/filters"
	m2 "github.com/chen-keinan/lxd-probe/pkg/models"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/chen-keinan/lxd-probe/ui"
	"github.com/mitchellh/colorstring"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
)

//LxdAudit lxd benchmark object
type LxdAudit struct {
	Command         shell.Executor
	ResultProcessor ResultProcessor
	OutputGenerator ui.OutputGenerator
	FileLoader      TestLoader
	PredicateChain  []filters.Predicate
	PredicateParams []string
	PlChan          chan m2.LxdAuditResults
	CompletedChan   chan bool
	FilesInfo       []utils.FilesInfo
	log             *logger.LdxProbeLogger
}

// ResultProcessor process audit results
type ResultProcessor func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench

// ConsoleOutputGenerator print audit tests to stdout
var ConsoleOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.LdxProbeLogger) {
	grandTotal := make([]models.AuditTestTotals, 0)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Status", "Type", "Audit Test Description"})
	table.SetAutoWrapText(false)
	table.SetBorder(true) // Set
	for _, a := range at {
		categoryTotal := printTestResults(a.AuditTests, table, a.Name)
		grandTotal = append(grandTotal, categoryTotal)
	}
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)
	table.Render()
	log.Console(printFinalResults(grandTotal))
}

func printFinalResults(grandTotal []models.AuditTestTotals) string {
	finalTotal := calculateFinalTotal(grandTotal)
	passTest := colorstring.Color("[green]Pass:")
	failTest := colorstring.Color("[red]Fail:")
	warnTest := colorstring.Color("[yellow]Warn:")
	title := "Test Result Total:   "
	return fmt.Sprintf("%s %s %d , %s %d , %s %d ", title, passTest, finalTotal.Pass, warnTest, finalTotal.Warn, failTest, finalTotal.Fail)
}

func calculateFinalTotal(granTotal []models.AuditTestTotals) models.AuditTestTotals {
	var (
		warn int
		fail int
		pass int
	)
	for _, total := range granTotal {
		warn = warn + total.Warn
		fail = fail + total.Fail
		pass = pass + total.Pass
	}
	return models.AuditTestTotals{Pass: pass, Fail: fail, Warn: warn}
}

// ReportOutputGenerator print failed audit test to human report
var ReportOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.LdxProbeLogger) {
	for _, a := range at {
		log.Table(reports.GenerateAuditReport(a.AuditTests))
	}
}

// simpleResultProcessor process audit results to stdout print only
var simpleResultProcessor ResultProcessor = func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	return AddAllMessages(at, NumFailedTest)
}

// ResultProcessor process audit results to std out and failure results
var reportResultProcessor ResultProcessor = func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	// append failed messages
	return AddFailedMessages(at, NumFailedTest)
}

//NewLxdAudit new audit object
func NewLxdAudit(filters []string, plChan chan m2.LxdAuditResults, completedChan chan bool, fi []utils.FilesInfo) *LxdAudit {
	return &LxdAudit{Command: shell.NewShellExec(),
		PredicateChain:  buildPredicateChain(filters),
		PredicateParams: buildPredicateChainParams(filters),
		ResultProcessor: GetResultProcessingFunction(filters),
		OutputGenerator: getOutputGeneratorFunction(filters),
		FileLoader:      NewFileLoader(),
		PlChan:          plChan,
		FilesInfo:       fi,
		CompletedChan:   completedChan}
}

//Help return benchmark command help
func (ldx LxdAudit) Help() string {
	return startup.GetHelpSynopsis()
}

//Run execute the full lxd benchmark
func (ldx *LxdAudit) Run(args []string) int {
	// load audit tests fro benchmark folder
	auditTests := ldx.FileLoader.LoadAuditTests(ldx.FilesInfo)
	// filter tests by cmd criteria
	ft := filteredAuditBenchTests(auditTests, ldx.PredicateChain, ldx.PredicateParams)
	//execute audit tests and show it in progress bar
	completedTest := executeTests(ft, ldx.runAuditTest, ldx.log)
	// generate output data
	ui.PrintOutput(completedTest, ldx.OutputGenerator, ldx.log)
	// send test results to plugin
	sendResultToPlugin(ldx.PlChan, ldx.CompletedChan, completedTest)
	return 0
}

func sendResultToPlugin(plChan chan m2.LxdAuditResults, completedChan chan bool, auditTests []*models.SubCategory) {
	ka := m2.LxdAuditResults{BenchmarkType: "lxd", Categories: make([]m2.AuditBenchResult, 0)}
	for _, at := range auditTests {
		for _, ab := range at.AuditTests {
			var testResult = "FAIL"
			if ab.TestSucceed {
				testResult = "PASS"
			}
			abr := m2.AuditBenchResult{Category: at.Name, ProfileApplicability: ab.ProfileApplicability, Description: ab.Description, AuditCommand: ab.AuditCommand, Remediation: ab.Remediation, Impact: ab.Impact, AdditionalInfo: ab.AdditionalInfo, References: ab.References, TestResult: testResult}
			ka.Categories = append(ka.Categories, abr)
		}
	}
	plChan <- ka
	<-completedChan
}

// runAuditTest execute category of audit tests
func (ldx *LxdAudit) runAuditTest(at *models.AuditBench) []*models.AuditBench {
	auditRes := make([]*models.AuditBench, 0)
	if at.NonApplicable {
		auditRes = append(auditRes, at)
		return auditRes
	}
	cmdTotalRes := make([]string, 0)
	// execute audit test command
	for index := range at.AuditCommand {
		res := ldx.execCommand(at, index, cmdTotalRes, make([]IndexValue, 0))
		cmdTotalRes = append(cmdTotalRes, res)
	}
	// evaluate command result with expression
	NumFailedTest := ldx.evalExpression(at, cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0)
	// continue with result processing
	auditRes = append(auditRes, ldx.ResultProcessor(at, NumFailedTest)...)
	return auditRes
}

func (ldx *LxdAudit) addDummyCommandResponse(expr string, index int, n string) string {
	if n == "[^\"]\\S*'\n" || n == "" || n == common.EmptyValue {
		spExpr := utils.SeparateExpr(expr)
		for _, expr := range spExpr {
			if expr.Type == common.SingleValue {
				if !strings.Contains(expr.Expr, fmt.Sprintf("'$%d'", index)) {
					if strings.Contains(expr.Expr, fmt.Sprintf("$%d", index)) {
						return common.NotValidNumber
					}
				}
			}
		}
		return common.EmptyValue
	}
	return n
}

//IndexValue hold command index and result
type IndexValue struct {
	index int
	value string
}

func (ldx *LxdAudit) execCommand(at *models.AuditBench, index int, prevResult []string, newRes []IndexValue) string {
	cmd := at.AuditCommand[index]
	paramArr, ok := at.CommandParams[index]
	if ok {
		for _, param := range paramArr {
			paramNum, err := strconv.Atoi(param)
			if err != nil {
				ldx.log.Console(fmt.Sprintf("failed to convert param for command %s", cmd))
				continue
			}
			if paramNum < len(prevResult) {
				n := ldx.addDummyCommandResponse(at.EvalExpr, index, prevResult[paramNum])
				newRes = append(newRes, IndexValue{index: paramNum, value: n})
			}
		}
		commandRes := ldx.execCmdWithParams(newRes, len(newRes), make([]IndexValue, 0), cmd, make([]string, 0))
		sb := strings.Builder{}
		for _, cr := range commandRes {
			sb.WriteString(utils.AddNewLineToNonEmptyStr(cr))
		}
		return sb.String()
	}
	result, _ := ldx.Command.Exec(cmd)
	if result.Stderr != "" {
		ldx.log.Console(fmt.Sprintf("Failed to execute command %s\n %s", result.Stderr, cmd))
	}
	return ldx.addDummyCommandResponse(at.EvalExpr, index, result.Stdout)
}

func (ldx *LxdAudit) execCmdWithParams(arr []IndexValue, index int, prevResHolder []IndexValue, currCommand string, resArr []string) []string {
	if len(arr) == 0 {
		return ldx.execShellCmd(prevResHolder, resArr, currCommand, ldx.Command)
	}
	sArr := strings.Split(utils.RemoveNewLineSuffix(arr[0].value), "\n")
	for _, a := range sArr {
		prevResHolder = append(prevResHolder, IndexValue{index: arr[0].index, value: a})
		resArr = ldx.execCmdWithParams(arr[1:index], index-1, prevResHolder, currCommand, resArr)
		prevResHolder = prevResHolder[:len(prevResHolder)-1]
	}
	return resArr
}

func (ldx *LxdAudit) execShellCmd(prevResHolder []IndexValue, resArr []string, currCommand string, se shell.Executor) []string {
	for _, param := range prevResHolder {
		if param.value == common.EmptyValue || param.value == common.NotValidNumber || param.value == "" {
			resArr = append(resArr, param.value)
			break
		}
		cmd := strings.ReplaceAll(currCommand, fmt.Sprintf("#%d", param.index), param.value)
		result, _ := se.Exec(cmd)
		if result.Stderr != "" {
			ldx.log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
		}
		if len(strings.TrimSpace(result.Stdout)) == 0 {
			result.Stdout = common.EmptyValue
		}
		resArr = append(resArr, result.Stdout)
	}
	return resArr
}

//evalExpression expression eval as cartesian product
func (ldx *LxdAudit) evalExpression(at *models.AuditBench,
	commandRes []string, commResSize int, permutationArr []string, testFailure int) int {
	if len(commandRes) == 0 {
		return ldx.evalCommand(at, permutationArr, testFailure)
	}
	outputs := strings.Split(utils.RemoveNewLineSuffix(commandRes[0]), "\n")
	for _, o := range outputs {
		permutationArr = append(permutationArr, o)
		testFailure = ldx.evalExpression(at, commandRes[1:commResSize], commResSize-1, permutationArr, testFailure)
		permutationArr = permutationArr[:len(permutationArr)-1]
	}
	return testFailure
}

func (ldx *LxdAudit) evalCommand(at *models.AuditBench, permutationArr []string, testExec int) int {
	// build command expression with params
	expr := at.CmdExprBuilder(permutationArr, at.EvalExpr)
	testExec++
	// eval command expression
	testSucceeded, err := evalCommandExpr(strings.ReplaceAll(expr, common.EmptyValue, ""))
	if err != nil {
		ldx.log.Console(fmt.Sprintf("failed to evaluate command expr %s for audit test %s : err %s", expr, at.Name, err.Error()))
	}
	return testExec - testSucceeded
}

func evalCommandExpr(expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}

//Synopsis for help
func (ldx *LxdAudit) Synopsis() string {
	return ldx.Help()
}
