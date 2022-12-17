package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// 修改密码
// 三种方案
// 1. net user $username $password
// 2. pspasswd.exe $username $passowrd
// 3. 自行调用win32 dll api
// 方案1无法有效区分password和其他option, 假设password="/delete", 那么会识别成net user xxx /delete, 从而变成删除xxx用户
// 方案3实现起来太麻烦
// 用方案2时需要用/accepteula参数跳过license确认弹窗
func SetUserPassword(username, password string) error {
	// 判空
	if len(username) == 0 || len(password) == 0 {
		return fmt.Errorf("cannot set an empty username or password")
	}
	var stderr bytes.Buffer
	cmd := exec.Command("pspasswd.exe", username, password, "/accepteula")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if len(errMsg) > 0 {
			return errors.New(errMsg)
		}
		return err
	}
	return nil
}

// 校验密码
// 经测试,psTools系列工具里只有PsExec能验证本地用户的密码，其他都无法达到目的
// 但是这个工具能做到的东西比较广，可能存在通过注入实现任意命令执行的效果
func VerifyUserPassword(username, password string) error {
	// 判空
	if len(username) == 0 || len(password) == 0 {
		return fmt.Errorf("cannot set an empty username or password")
	}

	var stderr bytes.Buffer
	cmd := exec.Command("psexec.exe", "-u", username, "-p", password, "ipconfig", "/accepteula")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if len(errMsg) > 0 {
			return errors.New(errMsg)
		}
		return err
	}
	return nil
}
