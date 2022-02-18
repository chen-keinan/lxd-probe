package benchmark

import (
	"embed"
	"fmt"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"io/ioutil"
)

const LxdFolder = "lxd/v1.0.0"

var (
	//go:embed lxd/v1.0.0
	res embed.FS
)

func LoadLxdSpecs() ([]utils.FilesInfo, error) {
	dir, _ := res.ReadDir(LxdFolder)
	specs := make([]utils.FilesInfo, 0)
	for _, r := range dir {
		file, err := res.Open(fmt.Sprintf("%s/%s", LxdFolder, r.Name()))
		if err != nil {
			return specs, err
		}
		data, err := ioutil.ReadAll(file)
		spec := utils.FilesInfo{Name: r.Name(), Data: string(data)}
		if err != nil {
			return specs, err
		}
		if err != nil {
			return specs, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}
