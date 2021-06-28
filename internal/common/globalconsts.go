package common

const (
	//FilesystemConfiguration file name
	FilesystemConfiguration = "1.1_filesystem_configuration.yml"
	//ConfigureSoftwareUpdates file name
	ConfigureSoftwareUpdates = "1.2_configure_software_updates.yml"
	//ConfigureSudo file name
	ConfigureSudo = "1.3_configure_sudo.yml"
	//FilesystemIntegrityChecking file name
	FilesystemIntegrityChecking = "1.4_filesystem_integrity_checking.yml"
	//GrepRegex for tests
	GrepRegex = "[^\"]\\S*'"
	//MultiValue for tests
	MultiValue = "MultiValue"
	//SingleValue for tests
	SingleValue = "SingleValue"
	//EmptyValue for test
	EmptyValue = "EmptyValue"
	//NotValidNumber value
	NotValidNumber = "10000"
	//Report arg
	Report = "r"
	//Synopsis help
	Synopsis = "synopsis"
	//LdxProbeCli Name
	LdxProbeCli = "ldx-probe"
	//LdxProbeVersion version
	LdxProbeVersion = "0.1"
	//IncludeParam param
	IncludeParam = "i="
	//ExcludeParam param
	ExcludeParam = "e="
	//NodeParam param
	NodeParam = "n="
	//LxdProbeHomeEnvVar ldx probe Home env var
	LxdProbeHomeEnvVar = "LXD_PROBE_HOME"
	//LxdProbe binary name
	LxdProbe = "lxd-probe"
	//RootUser process user owner
	RootUser = "root"
	//NonApplicableTest test is not applicable
	NonApplicableTest = "non_applicable"
	//ManualTest test can only be manual executed
	ManualTest = "manual"
	//LxdBenchAuditResultHook hook name
	LxdBenchAuditResultHook = "LxdBenchAuditResultHook"
)
