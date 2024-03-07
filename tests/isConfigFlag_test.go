package tests

import (
	"codeAudit/utils"
	"testing"
)

func TestIsConfigFlagHappyPath(t *testing.T) {
	mockTerminalArg := "r=true"
	valueIsConfig := utils.IsConfigFlag(mockTerminalArg)
	if !valueIsConfig {
		t.Fatalf("The argument used wasn't a valid config options argument. Found %s", mockTerminalArg)
	}
}

func TestIsConfigFlagFailure(t *testing.T) {
	mockTerminalArg := "-f"
	valueIsConfig := utils.IsConfigFlag(mockTerminalArg)
	if valueIsConfig {
		t.Fatalf("The argument used was unexpectedly a valid config options argument. Found %s", mockTerminalArg)
	}

}
