package main

type mailConfig struct {
	Title    string `mapstructure:"title"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Remote   string `mapstructure:"remote"`
	URL      string `mapstructure:"url"`
}

type Config struct {
	MailConfList []mailConfig `mapstructure:"mailList"`
}
