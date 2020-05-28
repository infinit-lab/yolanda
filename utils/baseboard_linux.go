package utils

import (
	"github.com/infinit-lab/yolanda/logutils"
	"os/exec"
)

func GetBaseBoardUUID() (string, error) {
	cmd := exec.Command("dmidecode", "--string", "system-uuid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logutils.Error("Failed to CombineOutput. error: ", err)
		return "", err
	}
	return string(out), nil
}
