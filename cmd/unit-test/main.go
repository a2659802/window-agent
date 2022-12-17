package main

import (
	"log"

	"github.com/a2659802/window-agent/pkg/utils"
)

func main() {
	// bootstrap: net user test /add

	// change test -> 123
	if err := utils.SetUserPassword("test", "123"); err != nil {
		log.Printf("set pwd error:%v", err.Error())
	}
	// verify test -> 123
	if err := utils.VerifyUserPassword("test", "123"); err != nil {
		log.Printf("verify pwd error:%v", err.Error())
	}

}
