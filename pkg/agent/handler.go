package agent

/*
 这个文件可以拆分到两个独立package: handler, model
 如果要这样做，需要调整handler的格式为func(responseWriter, request)
*/

import (
	"encoding/json"
	"fmt"

	"github.com/a2659802/window-agent/pkg/message"
)

// 验证密码、修改密码消息
type PwdAction int

const (
	ActionVerify = PwdAction(iota)
	ActionChange
)

type PasswordMessage struct {
	Action   PwdAction `json:"action"`
	ID       string    `json:"id"`
	UserName string    `json:"username"`
	Password string    `json:"password"`
}

func (a Agent) processPwdMessage(data []byte) {
	var pwdMsg PasswordMessage
	var err error

	if err = json.Unmarshal(data, &pwdMsg); err != nil {
		a.response(pwdMsg.ID, message.StatusBadRequest, err.Error())
		return
	}

	// do something
	switch pwdMsg.Action {
	case ActionChange:
		err = a.processChangePassword(pwdMsg)
	case ActionVerify:
		err = a.processVerifyPassword(pwdMsg)
	default:
		err = fmt.Errorf("unknown action")
	}

	if err != nil {
		a.response(pwdMsg.ID, message.StatusFailed, err.Error())
		return
	}
	// success
	a.response(pwdMsg.ID, message.StatusOK)
}

func (a *Agent) processChangePassword(msg PasswordMessage) error {
	return fmt.Errorf("unimplement")
}

func (a *Agent) processVerifyPassword(msg PasswordMessage) error {
	return fmt.Errorf("unimplement")
}
