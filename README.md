# frps-panel

[README](README.md) | [中文文档](README_zh.md)

frp server plugin to support multiple users for [frp](https://github.com/fatedier/frp).

frps-multiuser will run as one single process and accept HTTP requests from frps.

![用户列表](screenshots/user%20list.png)
![新增列表](screenshots/new%20user.png)
![支持英文](screenshots/i18n.png)
![自动深色模式](screenshots/dark%20mode.png)

## Update Notes

+ **the default tokens file is frps-multiuser.ini now,ini file support comment**
+ **remove `-l`,it configure in `frps-multiuser.ini` now**
+ **change `-f` to `-c`,the same as `frps`**
+ **if \[users\] section is empty,the authentication will only be handle by frps**
+ **if user under \[disabled\] section ,and the value is `disable`, it means that user is be disabled, and can not connect to server**
+ **add a manage ui, and change color mode base on browser**
+ **you can dynamic `add`,`remove`,`disable` or `enable` user now**
+ **you can limit `ports`,`domains` and `subdomains` for each user now**

***when a user is dynamic been `remove` or `disable`,it will take some time to be effective***

***the limit of `ports`、`domains`、`subdomains` only effective at `NewProxy`***

### Features

* Support multiple user authentication by tokens saved in file.
* Support dynamic `add`,`remove`,`disable` or `enable` user
* Limit `ports`,`domains` and `subdomains` for each user

### Download

Download frps-multiuser binary file from [Release](../../releases).

### Requirements

frp version >= v0.31.0

### Usage

1. Create file `frps-multiuser.ini` including all support usernames and tokens.

```ini
[common]
;plugin listen ip
plugin_addr = 127.0.0.1
;plugin listen port
plugin_port = 7200
;the username of manage ui,optional
admin_user  = admin
;the password of manage ui,optional
admin_pwd   = admin

[users]
;user user1 with meta_token 123
user1 = 123
;user user2 with meta_token abc
user2 = abc

[ports]
;user1 can only use ports 8080,9090 to 9010 ,other ports will fail to create proxy (frpc can normally startup)
user1=8080,9090-9010

[domains]
;user1 can only use domain web01.user1.com ,other domain will fail to create proxy (frpc can normally startup)
user1=web01.user1.com

[subdomains]
;user1 can only use subdomain web01 ,other subdomain will fail to create proxy (frpc can normally startup)
user1=web01

[disabled]
;user2 is disabled,when frpc use this user to connect with frps,if frpc is not startup,it cannot startup,if it's already startup,it will always show error logs on console
user2 = disable
```

   One user each line. Username and token are split by `=`.

2. Run frps-multiuser:

   `./frps-multiuser -c ./frps-multiuser.ini`

3. Register plugin in frps.

```ini
# frps.ini
[common]
bind_port = 7000

[plugin.multiuser]
addr = 127.0.0.1:7200
path = /handler
ops = Login,NewWorkConn,NewUserConn,NewProxy,Ping
```

4. Specify username and meta_token in frpc configure file.

   For user1:

```ini
# frpc.ini
[common]
server_addr = x.x.x.x
server_port = 7000
user = user1
meta_token = 123

[ssh]
type = tcp
local_port = 22
remote_port = 8080
```

   For user2:(user2 cannot connect to server,because it is disabled)

```ini
# frpc.ini
[common]
server_addr = x.x.x.x
server_port = 7000
user = user2
meta_token = abc

[ssh]
type = tcp
local_port = 22
remote_port = 6000
```

## Run as service

this example is for `ubuntu` and with `root` user

+ 1.unzip `frps-multiuser.zip` to dir `/root/frps-multiuser`
+ 2.touch a file with command `touch frps-multiuser.service` in dir `/root/frps-multiuser`.the file content is:
```ini
[Unit]
Description = frp multiuser service
After = network.target syslog.target
Wants = network.target

[Service]
Type = simple
# config of frps-multiuser.ini,you should change the file path
Environment=FRPS_MULTIUSER_OPTS="-c /root/frps-multiuser/frps-multiuser.ini"
# command of run frps-multiuser,you should change the file path
ExecStart = /root/frps-multiuser/frps-multiuser $FRPS_MULTIUSER_OPTS

[Install]
WantedBy = multi-user.target
```
+ 3.copy `frps-multiuser.service` to `/etc/systemd/system/` with command `cp /root/frps-multiuser.service /etc/systemd/system/`
+ 4.reload service with command `systemctl daemon-reload`
+ 5.start service with command `service frps-multiuser start`

## Issues & Ideas

___If you want visit mange ui from internet, you should change `plugin_addr` to `0.0.0.0`___

If you have any issues or ideas, put it on [issues](https://github.com/yhl452493373/frps-multiuser/issues). I will try my best to achieve it.

## Credits

+ [frp](https://github.com/fatedier/frp)
+ [fp-multiuser](https://github.com/gofrp/fp-multiuser)
+ [layui](https://github.com/layui/layui)
+ [layui-theme-dark](https://github.com/Sight-wcg/layui-theme-dark)
