package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
	"time"
)

var (
	cBuff []byte
	err   error
)

type ipAddInfo struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

func QueryIPaddress(ip string) ipAddInfo {
	searcher, _ := xdb.NewWithBuffer(cBuff)
	defer searcher.Close()
	region, _ := searcher.SearchByStr(ip)

	regions := strings.Split(region, "|")
	ipInfo := ipAddInfo{
		Country:  regions[0],
		Province: regions[2],
		City:     regions[3],
		ISP:      regions[4],
	}
	return ipInfo
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
	addrInfo := QueryIPaddress(ip)

	format := c.DefaultQuery("format", "string")
	// 根据提交的不同参数返回不同格式数据
	switch format {
	case "json":
		c.JSON(200, gin.H{"time": formatTimeStr, "addr": ip, "country": addrInfo.Country, "province": addrInfo.Province, "city": addrInfo.City, "isp": addrInfo.ISP})
	case "yaml":
		c.YAML(200, gin.H{"time": formatTimeStr, "addr": ip, "country": addrInfo.Country, "province": addrInfo.Province, "city": addrInfo.City, "isp": addrInfo.ISP})
	case "xml":
		c.XML(200, gin.H{"time": formatTimeStr, "addr": ip, "country": addrInfo.Country, "province": addrInfo.Province, "city": addrInfo.City, "isp": addrInfo.ISP})
	default:
		c.String(200, "当前时间："+formatTimeStr+"\n"+"IP地址："+ip+"\n"+"国家："+addrInfo.Country+"\n"+"省份："+addrInfo.Province+"\n"+"城市："+addrInfo.City+"\n"+"运营商："+addrInfo.ISP+"\n")
	}
}

func main() {
	// 定义ip2region.db文件路径全局加载到内存
	var dbPath = "./ip2region.xdb"
	cBuff, err = xdb.LoadContentFromFile(dbPath)

	r := gin.Default()
	// 如果应用程序不在代理之后，“ForwardedByClientIP”应设置为 false，因此“X-Forwarded-For”将被忽略。
	// 如果在代理后面将其设置为true
	r.ForwardedByClientIP = true
	r.GET("/ip", FormatIP)
	r.Run(":18888")
}
