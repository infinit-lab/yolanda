// +build windows

package utils

import (
	"github.com/infinit-lab/yolanda/logutils"
	"os/exec"
	"strings"
)

func GetBaseBoardUUID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logutils.Error("Failed to CombineOutput. error: ", err)
		return "", nil
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line != "" && strings.Contains(line, "UUID") == false {
			line = strings.ReplaceAll(line, " ", "")
			return line, nil
		}
	}
	return "", nil
}
