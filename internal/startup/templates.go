package startup

import (
	"fmt"
	"github.com/chen-keinan/lxd-probe/internal/common"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
)

//GenerateLxdBenchmarkFiles use packr to load benchmark audit test yaml
func GenerateLxdBenchmarkFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../benchmark/lxd/v1.0.0/")
	// Add Master Node Configuration tests
	mnc, err := box.FindString(common.FilesystemConfiguration)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load lxd benchmarks audit tests %s  %s", common.FilesystemConfiguration, err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.FilesystemConfiguration, Data: mnc})
	su, err := box.FindString(common.ConfigureSoftwareUpdates)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load lxd benchmarks audit tests %s %s", common.ConfigureSoftwareUpdates, err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ConfigureSoftwareUpdates, Data: su})
	cs, err := box.FindString(common.ConfigureSudo)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load lxd benchmarks audit tests %s %s", common.ConfigureSudo, err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ConfigureSudo, Data: cs})
	return fileInfo, nil
}

//GetHelpSynopsis get help synopsis file
func GetHelpSynopsis() string {
	box := packr.NewBox("./../cli/commands/help/")
	// Add Master Node Configuration tests
	hs, err := box.FindString(common.Synopsis)
	if err != nil {
		panic(fmt.Sprintf("faild to load cli help synopsis %s", err.Error()))
	}
	return hs
}

//SaveBenchmarkFilesIfNotExist create benchmark audit file if not exist
func SaveBenchmarkFilesIfNotExist(spec, version string, filesData []utils.FilesInfo) error {
	fm := utils.NewKFolder()
	folder, err := utils.GetBenchmarkFolder(spec, version, fm)
	if err != nil {
		return err
	}
	for _, fileData := range filesData {
		filePath := filepath.Join(folder, fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			_, err = f.WriteString(fileData.Data)
			if err != nil {
				return fmt.Errorf("failed to write benchmark file")
			}
			err = f.Close()
			if err != nil {
				return fmt.Errorf("faild to close file %s", filePath)
			}
		}
	}
	return nil
}
