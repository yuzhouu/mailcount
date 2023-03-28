package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

func main() {
	// 读取配置文件
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan struct{}, 1)

	// 使用WaitGroup等待所有goroutine结束
	var wg sync.WaitGroup
	for _, mailConf := range conf.MailList {
		wg.Add(1)
		go func(mc mailConfig) {
			client := newMailClient(mc)
			client.subscribe(stopCh)

			defer client.unSubscribe()
			defer wg.Done()

		}(mailConf)
	}

	// 捕捉 Ctrl+C 信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigCh
	log.Println("Received signal to stop")
	close(stopCh)

	// 停止所有goroutine并等待它们退出
	log.Println("Stopping all goroutines")
	wg.Wait()
	log.Println("All goroutines stopped")
}
