package bplugin

import "plugin"

//LxdBenchAuditResultHook hold the plugin symbol for Lxd bench audit result Hook
type LxdBenchAuditResultHook struct {
	Plugins []plugin.Symbol
}
