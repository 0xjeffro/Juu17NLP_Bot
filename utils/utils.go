package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InArrayInt64(needle int64, haystack []int64) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func Int2Date(i int) string {
	// 把int时间戳转换为日期 xx年xx月xx日 xx:xx
	tm := time.Unix(int64(i), 0)
	return tm.Format("2006年01月02日 15:04")
}
