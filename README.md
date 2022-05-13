(github.com/rid-lin/go-tg-bots/for_Vasiliy)

# A simple telegram bot for spying on chats

This bot monitors several chats and sends a message if a keyword appears in these chats.

## Install

First you need to compile the tdlib
https://tdlib.github.io/td/build.html?language=Go

Then clone the repository

`git clone github.com/rid-lin/go-tg-bots/for_Vasiliy.git`

`cd for_Vasiliy`

`cp /config /etc/for_Vasiliy/`

Build programm:

`make build`

`make install`

Not relevant

```
Edit file /usr/share/for_Vasiliy/assets/for_Vasiliy.service

`nano /usr/share/for_Vasiliy/assets/for_Vasiliy.service`

and move to /lib/systemd/system

`mv /usr/share/for_Vasiliy/assets/for_Vasiliy.service /lib/systemd/system`

Make sure the log folder exists, If not then

`mkdir -p /var/log/for_Vasiliy/`

Configuring sistemd to automatically start the program

`systemctl daemon-reload`

`systemctl start for_Vasiliy`

`systemctl enable for_Vasiliy`
```