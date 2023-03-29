package main

import (
	"log"
	"strconv"

	_ "embed"

	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

//go:embed mail.ico
var mailIcon []byte

var menuMap = map[string]*menuItem{}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(mailIcon)
	systray.SetTooltip("未读邮件数量")
	mQuit := systray.AddMenuItem("退出", "退出程序")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	collectCh := make(chan struct{}, len(conf.MailConfList))

	go func() {
		for {
			<-collectCh
			count := 0
			for _, mi := range menuMap {
				count += mi.unreadCount
			}
			systray.SetTitle(strconv.Itoa(count))
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
