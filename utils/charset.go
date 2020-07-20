package utils

import (
	"github.com/infinit-lab/yolanda/logutils"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func ConvertGBKToUtf8(bytes []byte) string {
	decodeBytes, err := simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
	if err != nil {
		logutils.Error("Failed to Bytes. error: ", err)
		return ""
	}
	return string(decodeBytes)
}
