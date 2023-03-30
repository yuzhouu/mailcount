# Mailcount

Mailcount is a lightweight application that displays the number of unread emails in your mailbox on the system tray. This can be useful if you want to stay updated on your mailbox status without having to open your email client.

## Features

- Small: The app only occupies 8MB of disk space, making it a low-impact addition to your system.
- Simple: The app has no GUI, uses very little CPU, and requires minimal memory usage.
- Secure: Mailcount is open source on GitHub, so you can review the code for added security and transparency.

## How to Use

To use Mailcount, you need to create a configuration file. For example, you can create a file at `$HOME/.mailcount/config.yaml` and input the following content:

```yaml
mailList:
  - title: outlook # id tell you which mailbox it is.
    username: username
    password: pass
    remote: outlook.office365.com:993 # the imap server.
    url: https://outlook.live.com/mail/0/ # optional, the application will open this url if you click on menu.
```

### Common Errors

- Err01: If the application shows Err01 beside the mailbox icon, it means it cannot find the config file. Make sure the config file is located at `$HOME/.mailcount/config.yaml`.

- Err02: If the application shows Err01 beside the mailbox icon, it means the config file is malformed.

## Tested Mail Providers:

- [x] outlook
- [ ] 163

**If your mail provider supports IMAP, this application should work well with it regardless of whether it is listed or not.**
