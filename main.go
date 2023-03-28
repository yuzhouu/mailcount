package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

type MailConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Remote   string `mapstructure:"remote"`
}

type Config struct {
	MailList []MailConfig `mapstructure:"mailList"`
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

	// 连接到邮件服务器
	c, err := client.DialTLS(conf.Remote, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// 进行身份验证
	if err := c.Login(conf.Username, conf.Password); err != nil {
		log.Fatal(err)
	}

	// 选择收件箱
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// 创建等待新邮件的 channel
	updates := make(chan imap.Update)
	c.Notify(updates, imap.NewSearchCriteria())

	// 创建系统托盘图标
	systray.Run(onReady, onExit)

	// 等待新邮件到达的通知
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case update := <-updates:
				// 有新邮件到达，更新未读邮件数量
				if update.MailboxStatus != nil {
					// 更新系统托盘图标上的未读邮件数量
					systray.SetTitle(fmt.Sprintf("%d", update.MailboxStatus.Unseen))
				}
			case <-done:
				return
			}
		}
	}()

	// 进入 IDLE 状态，等待新邮件到达
	if err := c.Idle(); err != nil {
		log.Fatal(err)
	}
}

func onReady() {
	// 创建系统托盘菜单项
	m := systray.AddMenuItem("退出", "退出程序")

	go func() {
		for {
			// 等待菜单项被点击
			select {
			case <-m.ClickedCh:
				// 点击了退出菜单项，退出程序
				systray.Quit()
				return
			default:
				// 更新系统托盘图标上的未读邮件数量
				title := systray.Title()
				unseen, err := strconv.Atoi(title)
				if err == nil {
					if unseen > 0 {
						systray.SetTooltip(fmt.Sprintf("您有 %d 封未读邮件", unseen))
					} else {
						systray.SetTooltip("您没有未读邮件")
					}
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

func onExit() {
	// 程序退出时关闭系统托盘
	systray.Quit()
}
