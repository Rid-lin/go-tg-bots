(github.com/rid-lin/go-tg-bots)

# A simple telegram bot for spying on chats

This bot monitors several chats and sends a message if a keyword appears in these chats.

## Install

First you need to compile the tdlib
https://tdlib.github.io/td/build.html?language=Go

Then clone the repository

`git clone github.com/rid-lin/go-tg-bots`

`cd for_Vasiliy`

`cp /config /etc/go-tg-bots`

Build programm:

`make build`

`make install`

Not relevant

```
Edit file /usr/share/for_Vasiliy/go-tg-bots.service

`nano /usr/share/for_Vasiliy/go-tg-bots.service`

and move to /lib/systemd/system

`mv /usr/share/for_Vasiliy/go-tg-bots.service /lib/systemd/system`

Make sure the log folder exists, If not then

`mkdir -p /var/log/go-tg-bots/`

Configuring sistemd to automatically start the program

`systemctl daemon-reload`

`systemctl start go-tg-bots`

`systemctl enable go-tg-bots`
```