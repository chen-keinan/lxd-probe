package bplugin

import (
	"fmt"
	"github.com/chen-keinan/lxd-probe/pkg/models"
	"go.uber.org/zap"
)

//PluginWorker instance which match command data to specific pattern
type PluginWorker struct {
	cmd *PluginWorkerData
	log *zap.Logger
}

//NewPluginWorker return new plugin worker instance
func NewPluginWorker(commandMatchData *PluginWorkerData, log *zap.Logger) *PluginWorker {
	return &PluginWorker{cmd: commandMatchData, log: log}
}

//NewPluginWorkerData return new plugin worker instance
func NewPluginWorkerData(plChan chan models.LxdAuditResults, hook LxdBenchAuditResultHook, completedChan chan bool) *PluginWorkerData {
	return &PluginWorkerData{plChan: plChan, plugins: hook, completedChan: completedChan}
}

//PluginWorkerData encapsulate plugin worker properties
type PluginWorkerData struct {
	plChan        chan models.LxdAuditResults
	completedChan chan bool
	plugins       LxdBenchAuditResultHook
}

//Invoke invoke plugin accept audit bench results
func (pm *PluginWorker) Invoke() {
	go func() {
		ae := <-pm.cmd.plChan
		if len(pm.cmd.plugins.Plugins) > 0 {
			for _, pl := range pm.cmd.plugins.Plugins {
				err := ExecuteLxdAuditResults(pl, ae)
				if err != nil {
					pm.log.Error(fmt.Sprintf("failed to execute plugins %s", err.Error()))
				}
			}
		}
		pm.cmd.completedChan <- true
	}()
}
