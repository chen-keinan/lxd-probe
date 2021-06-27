package cli

import (
	"github.com/chen-keinan/lxd-probe/internal/cli/commands"
	"github.com/chen-keinan/lxd-probe/internal/common"
	"github.com/chen-keinan/lxd-probe/internal/mocks"
	"github.com/chen-keinan/lxd-probe/internal/models"
	"github.com/chen-keinan/lxd-probe/internal/shell"
	m2 "github.com/chen-keinan/lxd-probe/pkg/models"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

//Test_StartCli tests
func Test_StartCli(t *testing.T) {
	fm := utils.NewKFolder()
	initBenchmarkSpecData(fm, ArgsData{SpecType: "lxd", SpecVersion: "v1.0.0"})
	files, err := utils.GetLxdBenchAuditFiles("lxd", "v1.0.0", fm)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(files), 8)
	assert.Equal(t, files[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, files[1].Name, common.APIServer)
	assert.Equal(t, files[2].Name, common.ControllerManager)
	assert.Equal(t, files[3].Name, common.Scheduler)
	assert.Equal(t, files[4].Name, common.Etcd)
	assert.Equal(t, files[5].Name, common.ControlPlaneConfiguration)
	assert.Equal(t, files[6].Name, common.WorkerNodes)
	assert.Equal(t, files[7].Name, common.Policies)
}

func Test_ArgsSanitizer(t *testing.T) {
	args := []string{"--a", "-b"}
	ad := ArgsSanitizer(args)
	assert.Equal(t, ad.Filters[0], "a")
	assert.Equal(t, ad.Filters[1], "b")
	assert.False(t, ad.Help)
	args = []string{}
	ad = ArgsSanitizer(args)
	assert.True(t, ad.Filters[0] == "")
	args = []string{"--Help"}
	ad = ArgsSanitizer(args)
	assert.True(t, ad.Help)
}

//Test_LxdProbeHelpFunc test
func Test_LxdProbeHelpFunc(t *testing.T) {
	cm := make(map[string]cli.CommandFactory)
	bhf := LxdProbeHelpFunc(common.LxdProbe)
	helpFile := bhf(cm)
	assert.True(t, strings.Contains(helpFile, "Available commands are:"))
	assert.True(t, strings.Contains(helpFile, "Usage: lxd-probe [--version] [--help] <command> [<args>]"))
}

//Test_createCliBuilderData test
func Test_createCliBuilderData(t *testing.T) {
	cmdArgs := []string{"a"}
	ad := ArgsSanitizer(os.Args[1:])
	cmdArgs = append(cmdArgs, ad.Filters...)
	cmds := make([]cli.Command, 0)
	completedChan := make(chan bool)
	plChan := make(chan m2.LxdAuditResults)
	// invoke cli
	cmds = append(cmds, commands.NewLxdAudit(ad.Filters, plChan, completedChan, []utils.FilesInfo{}))
	c := createCliBuilderData(cmdArgs, cmds)
	_, ok := c["a"]
	assert.True(t, ok)

}

//Test_InvokeCli test
func Test_InvokeCli(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa"}
	ab.EvalExpr = "'$0' != '';"
	ab.CommandParams = map[int][]string{}
	ab.CmdExprBuilder = utils.UpdateCmdExprParam
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "1234"}, nil).Times(1)
	tl := mocks.NewMockTestLoader(ctrl)
	tl.EXPECT().LoadAuditTests(nil).Return([]*models.SubCategory{{Name: "te", AuditTests: []*models.AuditBench{ab}}})
	completedChan := make(chan bool)
	plChan := make(chan m2.LxdAuditResults)
	go func() {
		<-plChan
		completedChan <- true
	}()
	kb := &commands.LxdAudit{Command: executor, ResultProcessor: commands.GetResultProcessingFunction([]string{}), FileLoader: tl, OutputGenerator: commands.ConsoleOutputGenerator, PlChan: plChan, CompletedChan: completedChan}
	cmdArgs := []string{"a"}
	cmds := make([]cli.Command, 0)
	// invoke cli
	cmds = append(cmds, kb)
	c := createCliBuilderData(cmdArgs, cmds)
	a, err := invokeCommandCli(cmdArgs, c)
	assert.NoError(t, err)
	assert.True(t, a == 0)
}

func Test_InitPluginFolder(t *testing.T) {
	fm := utils.NewKFolder()
	initPluginFolders(fm)
}

func Test_InitPluginWorker(t *testing.T) {
	completedChan := make(chan bool)
	plChan := make(chan m2.LxdAuditResults)
	go func() {
		plChan <- m2.LxdAuditResults{}
		completedChan <- true
	}()
	initPluginWorker(plChan, completedChan)

}
