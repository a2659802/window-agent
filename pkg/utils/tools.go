package utils

import (
	"os"
	"path/filepath"
)

// 检查应用启动权限
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// 设置工作目录为程序所在位置
func SetWorkingDirectory() (err error) {
	var dir string
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	err = os.Chdir(dir)
	return
}
