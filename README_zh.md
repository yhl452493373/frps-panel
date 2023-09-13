# frps-panel

[README](README.md) | [中文文档](README_zh.md)

frps-panel 是 [frp](https://github.com/fatedier/frp) 的一个服务端插件，用于查看服务器信息以及支持多用户鉴权。

frps-panel 会以一个单独的进程运行，并接收 frps 发送过来的 HTTP 请求。

![支持英文](screenshots/i18n.png)
![服务器信息](screenshots/server%20info.png)
![用户列表](screenshots/user%20list.png)
![新增列表](screenshots/new%20user.png)
![代理列表](screenshots/proxy%20list.png)
![代理流量统计](screenshots/proxy%20traffic%20statistics.png)
![自动深色模式](screenshots/dark%20mode.png)


### 功能

+ **支持展示服务器信息**
+ **支持多用户鉴权**
+ **动态`添加`、`删除`、`禁用`、`启用`用户**
+ **对用户的`端口`、`域名`、`二级域名`进行限制**

***用户被`删除`或`禁用`后，不会马上生效，需要等一段时间***

***用户`端口`、`域名`、`二级域名`限制仅在建立新连接(`NewProxy`)时生效***

### 下载

通过 [Release](../../releases) 页面下载对应系统版本的二进制文件到本地。

### 要求

需要 frp 版本 >= v0.31.0

### 使用示例

1. 创建 `frps-panel.ini` 文件，内容为所有支持的用户名和 token。

```ini
[common]
;插件监听地址
;如果上面配置中 tls_mode = true, 则把plugin_addr的值改为 https://127.0.0.1:7200
plugin_addr = 127.0.0.1
;插件端口
plugin_port = 7200
;插件管理页面账号,可选
admin_user  = admin
;插件管理页面密码,与账号一起进行鉴权,可选
admin_pwd   = admin

; frps 面板页面是否启用https访问,如果为true,则只能通过https访问
tls_mode = false
; tls_cert_file = cert.crt
; tls_key_file = cert.key

; frp服务器的看板页面信息，必须配置，且与frp服务器一致，否则无法获取服务器信息
dashboard_addr = 127.0.0.1
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd  = admin
; 如果frps的dashboard启用了tls，则把这个选项置为true，否则无法获取服务器信息
dashboard_tls  = false

[users]
;user1的meta_token为123
user1 = 123
;user2的meta_token为abc
user2 = abc

[ports]
;user1只能使用8080,9090到9010端口,其他端口则建立连接时返回失败(不影响客户端启动)
user1=8080,9090-9010

[domains]
;user1只能使用web01.yyy.zzz域名,配置了其他域名则建立连接时返回失败(不影响客户端启动)
user1=web01.user1.com

[subdomains]
;user1只能使用web01.xxx.yyy.zzz域名,配置了其他三级域名则建立连接时返回失败(不影响客户端启动)
user1=web01

[disabled]
;user2被禁用,frpc使用此账户与frps通信时,如果未启动则无法启动,如果已启动,则会一直打印错误日志
user2 = disable
```

    每一个用户占一行，用户名和 token 之间以 `=` 号分隔。

2. 运行 frps-panel，指定监听地址以及 token 存储文件路径。

    `./frps-panel -c ./frps-panel.ini`

3. 在 frps 的配置文件中注册插件，并启动。

```ini
# frps.ini
[common]
bind_port = 7000

[plugin.multiuser]
addr = 127.0.0.1:7200
path = /handler
ops = Login,NewWorkConn,NewUserConn,NewProxy,Ping
```

4. 在 frpc 中指定用户名，在 meta 中指定 token，用户名以及 `meta_token` 的内容需要和之前创建的 token 文件匹配。

    user1 的配置:

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

    user2 的配置:（由于示例文件中user2被禁用，因此无法连接）

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

## 以服务的形式运行

本实例是在 `ubuntu` 下， 以 `root` 用户执操作

+ 1、解压 `frps-panel.zip` 到目录 `/root/frps-panel`
+ 2、在目录 `/root/frps-panel` 下 用命令创建文件：`touch frps-panel.service`。创建后修改文件内容：
```ini
[Unit]
Description = frp multiuser service
After = network.target syslog.target
Wants = network.target

[Service]
Type = simple
# 启动frps-panel的配置文件路径，需修改为您的frps-panel.ini的路径
Environment=FRPS_PANEL_OPTS="-c /root/frps-panel/frps-panel.ini"
# 启动frps-panel的命令，需修改为您的frps-panel的安装路径
ExecStart = /root/frps-panel/frps-panel $FRPS_PANEL_OPTS

[Install]
WantedBy = multi-user.target
```
+ 3、复制服务文件： `cp /root/frps-panel.service /etc/systemd/system/`
+ 4、重载服务： `systemctl daemon-reload`
+ 5、启动服务： `service frps-panel start`

## 使用

___如果要从外网访问管理界面, 需要把配置中的 `plugin_addr` 改为 `0.0.0.0`___

如果使用中有问题或者有其他想法，在[issues](https://github.com/yhl452493373/frps-panel/issues)上提出来。 如果我能搞定的话，我尽量搞。

## 致谢

+ [frp](https://github.com/fatedier/frp)
+ [fp-multiuser](https://github.com/gofrp/fp-multiuser)
+ [layui](https://github.com/layui/layui)
+ [layui-theme-dark](https://github.com/Sight-wcg/layui-theme-dark)
+ [echarts](https://github.com/apache/echarts)
