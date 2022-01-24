package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"strings"
	"time"
)

func QueryIPaddress(ip string) string {
	region, _ := ip2region.New("ip2region.db")
	defer region.Close()
	infodata, _ := region.MemorySearch(ip)
	return infodata.String()
}

func FormatIP(c *gin.Context) {
	var ip string
	// 判断查询指定ip还是查询客户端请求的公网ip（默认不加ip参数为查询请求者的公网ip）
	IpAdd := c.DefaultQuery("ip", "")
	if IpAdd == "" {
		ip = c.ClientIP()
	} else {
		ip = IpAdd
	}

	// 将时间戳转换成时间
	formatTimeStr := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	addrInfo := strings.Replace(QueryIPaddress(ip), "0|", "", -1)

	format := c.DefaultQuery("format", "string")
	// 根据提交的不同参数返回不同格式数据
	if format == "string" {
		c.String(200, "当前时间："+formatTimeStr+"\n"+"IP地址："+ip+"\n"+"归属地："+addrInfo+"\n")
	} else if format == "json" {
		c.JSON(200, gin.H{"time": formatTimeStr, "ipaddr": ip, "addr": addrInfo})
	} else if format == "yaml" {
		c.YAML(200, gin.H{"time": formatTimeStr, "ipaddr": ip, "addr": addrInfo})
	} else if format == "xml" {
		c.XML(200, gin.H{"time": formatTimeStr, "ipaddr": ip, "addr": addrInfo})
	}
}

func main() {
	r := gin.Default()
	// 如果应用程序不在代理之后，“ForwardedByClientIP”应设置为 false，因此“X-Forwarded-For”将被忽略。
	// 如果在代理后面将其设置为true
	r.ForwardedByClientIP = true
	r.GET("/ip", FormatIP)
	r.Run(":18888")
}
