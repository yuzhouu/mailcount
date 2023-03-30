# mailcount

show unread mail count in system tray

## features

- simple: small application, no gui, very low cpu, very low memory need
- security: open source on github, you can review the code

## instructions

to use the app, you need create a config file, example `$HOME/.mailcount/config.yaml`, the content like this

```yaml
mailList:
  - title: outlook # id tell you which mailbox it is.
    username: username
    password: pass
    remote: outlook.office365.com:993 # the imap server.
    url: https://outlook.live.com/mail/0/ # optional, the application will open this url if you click on menu.
```

### errors

- Err01: if the application show `Err01` beside mailbox icon means it can't find config file, please make sure you place have config file at `$HOME/.mailcount/config.yaml`.

- Err02: if the application show `Err01` beside mailbox icon means the config file is malformat.

## support

- [x] outlook
- [ ] 163
