# ip-info
The query client requests the attribution of the public network IP, queries the attribution of the specified public network IP, and displays the data in the specified format.



## Build

```bash
$ go mod tidy
$ go build -o ipinfo main.go
```



## Run

```bash
$ ./ipinfo

# 默认请求者ip
$ curl http://192.168.0.92:18888/ip
当前时间：2022-01-24 18:29:31
IP地址：192.168.0.92
归属地：内网IP|内网IP

# 查询指定ip
$ curl http://192.168.0.92:18888/ip\?ip\=114.114.114.114
当前时间：2022-01-24 18:29:56
IP地址：114.114.114.114
归属地：中国|江苏|南京

# 指定json格式访问，支持json、xml、yaml
$ curl http://192.168.0.92:18888/ip\?ip\=114.114.114.114\&format\=json
{"addr":"中国|江苏|南京|0","ipaddr":"114.114.114.114","time":"2022-01-24 18:30:26"}
```



