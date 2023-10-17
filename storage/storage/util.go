package storage

import (
	"github.com/freddyfeng-fy/mucy-core/digests"
	"github.com/freddyfeng-fy/mucy-core/utils/dates"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/go-resty/resty/v2"
	"mime"
	"os"
	"strings"
	"time"
)

func NormalizeKey(key string) string {
	key = strings.Replace(key, "\\", "/", -1)
	key = strings.Replace(key, " ", "", -1)
	key = filterNewLines(key)

	return key
}

func filterNewLines(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029:
			return -1
		default:
			return r
		}
	}, s)
}

func OpenAsReadOnly(key string) (*os.File, os.FileInfo, error) {
	fd, err := os.Open(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, FileNotFoundErr
		}
		if os.IsPermission(err) {
			return nil, nil, FileNoPermissionErr
		}
		return nil, nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, nil, err
	}

	return fd, stat, nil
}

// generateKey 生成图片Key
func GenerateImageKey(data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := ""
	if strs.IsNotBlank(contentType) {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil || len(exts) > 0 {
			ext = exts[0]
		}
	}
	return "images/" + dates.Format(time.Now(), "2006/01/02/") + md5 + ext
}

func GenerateVideoKey(data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := ""
	if strs.IsNotBlank(contentType) {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil || len(exts) > 0 {
			ext = exts[0]
		}
	}
	return "video/" + dates.Format(time.Now(), "2006/01/02/") + md5 + ext
}

func Download(url string) ([]byte, string, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, "", err
	}
	return rsp.Body(), rsp.Header().Get("Content-Type"), nil
}
