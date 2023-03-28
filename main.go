package main

import (
	"log"
	"sync"

	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

type mailConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Remote   string `mapstructure:"remote"`
}

type Config struct {
	MailList []mailConfig `mapstructure:"mailList"`
}

var stopCh = make(chan struct{}, 1)
var wg sync.WaitGroup

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("未读邮件")
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

	for _, mailConf := range conf.MailList {
		wg.Add(1)
		go func(mc mailConfig) {
			client := newMailClient(mc)
			client.subscribe(stopCh)

			defer client.unSubscribe()
			defer wg.Done()

		}(mailConf)
	}
}

func onExit() {
	log.Println("quit called")
	close(stopCh)
	wg.Wait()
}
