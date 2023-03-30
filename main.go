package main

import (
	"log"
	"strconv"

	_ "embed"

	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

//go:embed mailbox.ico
var mailBoxIcon []byte

//go:embed warn.ico
var warnIcon []byte

var menuMap = map[string]*menuItem{}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(mailBoxIcon, mailBoxIcon)
	systray.SetTooltip("快速查看未读邮件数量")
	mQuit := systray.AddMenuItem("退出", "退出程序")
	mHelp := systray.AddMenuItem("使用说明", "查看使用说明")
	go func() {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		case <-mHelp.ClickedCh:
			openbrowser("https://github.com/yuzhouu/mailcount")
		}
	}()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/mailcount")
	viper.AddConfigPath("$HOME/.mailcount")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		systray.SetTitle("Err01")
		log.Println(err)
		return
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		systray.SetTitle("Err02")
		log.Println(err)
		return
	}

	collectCh := make(chan struct{}, len(conf.MailConfList))

	go func() {
		for {
			<-collectCh
			count := 0
			for _, mi := range menuMap {
				count += mi.unreadCount
			}
			if count != 0 {
				systray.SetTitle(strconv.Itoa(count))
			} else {
				systray.SetTitle("")
			}
		}

	}()

	for _, mc := range conf.MailConfList {
		menuMap[mc.Title] = newMenuItem(mc, collectCh)
	}

	for _, mi := range menuMap {
		go func(mi *menuItem) {
			mi.loop()
		}(mi)
	}
}

func onExit() {
	for _, mi := range menuMap {
		mi.unSubscribe()
	}
}
