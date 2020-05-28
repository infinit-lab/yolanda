package utils

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/infinit-lab/yolanda/logutils"
	"strings"
)

func GetMachineFingerprint() (string, error) {
	baseboardUuid, err := GetBaseBoardUUID()
	if err != nil {
		logutils.Error("Failed to GetBaseBoardUUID. error: ", err)
		return "", err
	}
	cpuId, err := GetCpuID()
	if err != nil {
		logutils.Error("Failed to GetCpuID. error: ", err)
		return "", err
	}
	diskIds, err := GetDiskSerialNumber()
	if err != nil {
		logutils.Error("Failed to GetDiskSerialNumber. error: ", err)
		return "", err
	}
	str := baseboardUuid + cpuId
	for _, id := range diskIds {
		str += id
	}
	encode := base64.StdEncoding.EncodeToString([]byte(str))
	has := md5.Sum([]byte(encode))
	encode = base64.StdEncoding.EncodeToString(has[:])
	encode = strings.ReplaceAll(encode, "+", "-")
	encode = strings.ReplaceAll(encode, "/", "_")
	return encode, nil
}
