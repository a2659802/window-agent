package utils

import (
	"os"
)

// 检查应用启动权限
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
