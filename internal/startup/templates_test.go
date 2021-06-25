package startup

import (
	"github.com/chen-keinan/lxd-probe/internal/common"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//Test_CreateLxdBenchmarkFilesIfNotExist test
func Test_CreateLxdBenchmarkFilesIfNotExist(t *testing.T) {
	bFiles, err := GenerateLxdBenchmarkFiles()
	if err != nil {
		t.Fatal(err)
	}
	// generate test with packr
	assert.Equal(t, bFiles[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, bFiles[1].Name, common.APIServer)
	assert.Equal(t, bFiles[2].Name, common.ControllerManager)
	assert.Equal(t, bFiles[3].Name, common.Scheduler)
	assert.Equal(t, bFiles[4].Name, common.Etcd)
	assert.Equal(t, bFiles[5].Name, common.ControlPlaneConfiguration)
	assert.Equal(t, bFiles[6].Name, common.WorkerNodes)
	assert.Equal(t, bFiles[7].Name, common.Policies)
	fm := utils.NewKFolder()
	err = utils.CreateBenchmarkFolderIfNotExist("lxd", "v1.6.0", fm)
	assert.NoError(t, err)
	// save benchmark files to folder
	err = SaveBenchmarkFilesIfNotExist("lxd", "v1.6.0", bFiles)
	assert.NoError(t, err)
	// fetch files from benchmark folder
	bFiles, err = utils.GetLxdBenchAuditFiles("lxd", "v1.6.0", fm)
	assert.Equal(t, bFiles[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, bFiles[1].Name, common.APIServer)
	assert.Equal(t, bFiles[2].Name, common.ControllerManager)
	assert.Equal(t, bFiles[3].Name, common.Scheduler)
	assert.Equal(t, bFiles[4].Name, common.Etcd)
	assert.Equal(t, bFiles[5].Name, common.ControlPlaneConfiguration)
	assert.Equal(t, bFiles[6].Name, common.WorkerNodes)
	assert.Equal(t, bFiles[7].Name, common.Policies)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}

//Test_CreateGkeBenchmarkFilesIfNotExist test
func Test_CreateGkeBenchmarkFilesIfNotExist(t *testing.T) {
	bFiles, err := GenerateGkeBenchmarkFiles()
	if err != nil {
		t.Fatal(err)
	}
	// generate test with packr
	assert.Equal(t, bFiles[0].Name, common.GkeControlPlaneConfiguration)
	assert.Equal(t, bFiles[1].Name, common.GkeWorkerNodes)
	assert.Equal(t, bFiles[2].Name, common.GkePolicies)
	assert.Equal(t, bFiles[3].Name, common.GkeManagedServices)
	fm := utils.NewKFolder()
	err = utils.CreateBenchmarkFolderIfNotExist("gke", "v1.1.0", fm)
	assert.NoError(t, err)
	// save benchmark files to folder
	err = SaveBenchmarkFilesIfNotExist("gke", "v1.1.0", bFiles)
	assert.NoError(t, err)
	// fetch files from benchmark folder
	bFiles, err = utils.GetLxdBenchAuditFiles("gke", "v1.1.0", fm)
	assert.Equal(t, bFiles[0].Name, common.GkeControlPlaneConfiguration)
	assert.Equal(t, bFiles[1].Name, common.GkeWorkerNodes)
	assert.Equal(t, bFiles[2].Name, common.GkePolicies)
	assert.Equal(t, bFiles[3].Name, common.GkeManagedServices)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}

//Test_GetHelpSynopsis test
func Test_GetHelpSynopsis(t *testing.T) {
	hs := GetHelpSynopsis()
	assert.True(t, len(hs) != 0)
}

//Test_SaveBenchmarkFilesIfNotExist test
func Test_SaveBenchmarkFilesIfNotExist(t *testing.T) {
	fm := utils.NewKFolder()
	folder, err2 := utils.GetBenchmarkFolder("lxd", "v1.6.0", fm)
	assert.NoError(t, err2)
	err := os.RemoveAll(folder)
	assert.NoError(t, err)
	filesData := make([]utils.FilesInfo, 0)
	err = utils.CreateBenchmarkFolderIfNotExist("lxd", "v1.6.0", fm)
	assert.NoError(t, err)
	filesData = append(filesData, utils.FilesInfo{Name: common.Scheduler, Data: "bbb"})
	err = SaveBenchmarkFilesIfNotExist("lxd", "v1.6.0", filesData)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}

//Test_SaveBenchmarkFilesIfNotExist test
func Test_SaveGkeBenchmarkFilesIfNotExist(t *testing.T) {
	fm := utils.NewKFolder()
	folder, err2 := utils.GetBenchmarkFolder("gke", "v1.1.0", fm)
	assert.NoError(t, err2)
	err := os.RemoveAll(folder)
	assert.NoError(t, err)
	filesData := make([]utils.FilesInfo, 0)
	err = utils.CreateBenchmarkFolderIfNotExist("gke", "v1.1.0", fm)
	assert.NoError(t, err)
	filesData = append(filesData, utils.FilesInfo{Name: common.Scheduler, Data: "bbb"})
	err = SaveBenchmarkFilesIfNotExist("gke", "v1.1.0", filesData)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}
