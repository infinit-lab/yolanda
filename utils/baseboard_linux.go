// +build linux

package utils

import (
	"github.com/infinit-lab/yolanda/logutils"
	"os/exec"
	"strings"
)

func GetBaseBoardUUID() (string, error) {
	cmd := exec.Command("dmidecode", "--string", "system-uuid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logutils.Error("Failed to CombineOutput. error: ", err)
		return "", err
	}
	line := strings.ReplaceAll(string(out), "\n", "")
	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, " ", "")
	return line, nil
}
