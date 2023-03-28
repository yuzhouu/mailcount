package main

import (
	"log"

	"github.com/emersion/go-imap/client"
)

type mailClient struct {
	client *client.Client
}

func newMailClient(conf mailConfig) *mailClient {
	c, err := client.DialTLS(conf.Remote, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Login(conf.Username, conf.Password); err != nil {
		log.Fatal(err)
	}

	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	return &mailClient{
		client: c,
	}
}

func (mc *mailClient) subscribe(stop chan struct{}) {
	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	mc.client.Updates = updates

	// Start idling
	stopped := false
	done := make(chan error, 1)
	go func() {
		done <- mc.client.Idle(stop, nil)
	}()

	// Listen for updates
	for {
		select {
		case update := <-updates:
			log.Println("New update:", update)
			if !stopped {
				close(stop)
				stopped = true
			}
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Not idling anymore")
			return
		}
	}
}

func (mc *mailClient) unSubscribe() {
	log.Println("unsubscribe")
	mc.client.Logout()
}
