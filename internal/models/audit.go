package models

import (
	"github.com/chen-keinan/lxd-probe/internal/common"
	"github.com/chen-keinan/lxd-probe/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

//Audit data model
type Audit struct {
	BenchmarkType string     `yaml:"benchmark_type"`
	Categories    []Category `yaml:"categories"`
}

//AuditTestTotals model
type AuditTestTotals struct {
	Warn int
	Pass int
	Fail int
}

//Category data model
type Category struct {
	Name        string       `yaml:"name"`
	SubCategory *SubCategory `yaml:"sub_category"`
}

//SubCategory data model
type SubCategory struct {
	Name       string        `yaml:"name"`
	AuditTests []*AuditBench `yaml:"audit_tests"`
}

//AuditBench data model
type AuditBench struct {
	Name                 string   `mapstructure:"name" yaml:"name"`
	ProfileApplicability string   `mapstructure:"profile_applicability" yaml:"profile_applicability"`
	Description          string   `mapstructure:"description" yaml:"description"`
	AuditCommand         []string `mapstructure:"audit" json:"audit"`
	CheckType            string   `mapstructure:"check_type" yaml:"check_type"`
	Remediation          string   `mapstructure:"remediation" yaml:"remediation"`
	Impact               string   `mapstructure:"impact" yaml:"impact"`
	AdditionalInfo       string   `mapstructure:"additional_info" yaml:"additional_info"`
	References           []string `mapstructure:"references" yaml:"references"`
	EvalExpr             string   `mapstructure:"eval_expr" yaml:"eval_expr"`
	CmdExprBuilder       utils.CmdExprBuilder
	TestSucceed          bool
	CommandParams        map[int][]string
	Category             string
	NonApplicable        bool
	TestType             string `mapstructure:"type" yaml:"type"`
}

//AuditResult data
type AuditResult struct {
	NumOfExec    int
	NumOfSuccess int
}

//UnmarshalYAML over unmarshall to add logic
func (at *AuditBench) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var res map[string]interface{}
	if err := unmarshal(&res); err != nil {
		return err
	}
	err := mapstructure.Decode(res, &at)
	if err != nil {
		return err
	}
	switch at.CheckType {
	case "multi_param":
		at.CmdExprBuilder = utils.UpdateCmdExprParam
	}
	at.CommandParams = make(map[int][]string)
	for index, command := range at.AuditCommand {
		findIndex(command, "#", index, at.CommandParams)
	}
	if at.TestType == common.NonApplicableTest || at.TestType == common.ManualTest {
		at.NonApplicable = true
	}
	return nil
}

// find all params in command to be replace with output
func findIndex(s, c string, commandIndex int, locations map[int][]string) {
	b := strings.Index(s, c)
	var err error
	if b != -1 && len(s) >= b+2 {
		s2 := s[b+1 : b+2]
		_, err = strconv.Atoi(s2)
	}
	if b == -1 || err != nil {
		return
	}
	if locations[commandIndex] == nil {
		locations[commandIndex] = make([]string, 0)
	}
	locations[commandIndex] = append(locations[commandIndex], s[b+1:b+2])
	findIndex(s[b+2:], c, commandIndex, locations)
}
