package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
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

func generateKey(fingerprint string) ([]byte, error) {
	fingerprint = strings.ReplaceAll(fingerprint, "-", "+")
	fingerprint = strings.ReplaceAll(fingerprint, "_", "/")
	decode, err := base64.StdEncoding.DecodeString(fingerprint)
	if err != nil {
		logutils.Error("Failed to DecodeString. error: ", err)
		return nil, err
	}
	has := md5.Sum(decode)
	return has[:], nil
}

func Encode(fingerprint string, content []byte) (string, error) {
	key, err := generateKey(fingerprint)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		logutils.Error("Failed to NewCipher. error: ", err)
		return "", err
	}
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	dst := make([]byte, len(content))
	stream.XORKeyStream(dst, content)
	return hex.EncodeToString(dst), nil
}

func EncodeSelf(content []byte) (string, error) {
	key, err := GetMachineFingerprint()
	if err != nil {
		return "", err
	}
	return Encode(key, content)
}

func Decode(fingerprint string, encryptData string) ([]byte, error) {
	data, err := hex.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}
	key, err := generateKey(fingerprint)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		logutils.Error("Failed to NewCipher. error: ", err)
		return nil, err
	}
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	dst := make([]byte, len(data))
	stream.XORKeyStream(dst, data)
	return dst, nil
}

func DecodeSelf(encryptData string) ([]byte, error) {
	key, err := GetMachineFingerprint()
	if err != nil {
		return nil, err
	}
	return Decode(key, encryptData)
}
