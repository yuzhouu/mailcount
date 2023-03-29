package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/getlantern/systray"
)

type menuItem struct {
	title       string
	menu        *systray.MenuItem
	refreshMenu *systray.MenuItem
	mailClient  *client.Client
	stopCh      chan struct{}
	collectCh   chan struct{}
	unreadCount int
}

func newMenuItem(conf mailConfig, collectCh chan struct{}) *menuItem {
	c, err := client.DialTLS(conf.Remote, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Login(conf.Username, conf.Password); err != nil {
		log.Println(err)
	}

	menu := systray.AddMenuItem(conf.Title, "")
	refreshMenu := menu.AddSubMenuItem("刷新", "立即刷新未读邮件数量")
	openMenu := menu.AddSubMenuItem("打开邮箱", "在浏览器中打开邮箱")
	go func() {
		for {
			<-openMenu.ClickedCh
			openbrowser(conf.URL)
		}
	}()
	stopCh := make(chan struct{}, 1)

	return &menuItem{
		title:       conf.Title,
		mailClient:  c,
		menu:        menu,
		refreshMenu: refreshMenu,
		stopCh:      stopCh,
		collectCh:   collectCh,
	}
}

func (m *menuItem) loop() {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := m.mailClient.Search(criteria)
	if err != nil {
		log.Println(err)
	}

	m.menu.SetTitle(fmt.Sprintf("%s 未读邮件数量 %d", m.title, len(ids)))
	m.unreadCount = len(ids)
	m.collectCh <- struct{}{}

	m.subscribe()
}

func (m *menuItem) subscribe() {
	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	m.mailClient.Updates = updates

	// Start idling
	stopped := false
	stop := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- m.mailClient.Idle(stop, nil)
	}()

	// Listen for updates
	for {
		select {
		case <-updates:
			if !stopped {
				close(stop)
				stopped = true
			}
		case <-m.refreshMenu.ClickedCh:
			systray.SetTitle("loading")
			if !stopped {
				close(stop)
				stopped = true
			}
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
			m.loop()
			return

		case <-m.stopCh:
			close(stop)
			return
		}
	}
}

func (m *menuItem) unSubscribe() {
	close(m.stopCh)
}
