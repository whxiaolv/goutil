package goutil

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_tool(t *testing.T) {
	t.Run("testTool", testTool)
	//t.Run("testUrlEncode", testUrlEncode)
}

var (
	tool Tool
)

func testTool(t *testing.T) {

	//获取uuid，36位
	if len(tool.GetUUID()) != 36 {
		t.Error("testTool fail")
	}

	//随机范围正整数
	if tool.RandInt(10, 100) < 10 {
		t.Errorf("testTool fail ")
	}

	// 获取当前时间， 格式：20060102
	if len(tool.GetToDay()) != 8 {
		t.Errorf("testTool fail ")
	}

	// 时间戳转化为 格式：20060102
	tool.GetDayFromTimeStamp(1693555418)

	// 字符串 => 不可重复map 【更多方法参考源码 （github.com/deckarep/golang-set）】
	setMap := tool.StrToMapSet("a,a,b,c", ",")
	if len(setMap.ToSlice()) != 3 {
		t.Errorf("testTool fail ")
	}

	// map => 字符串,_ 符号隔开 strMap.Union(strMap1)
	strMap := tool.MapSetToStr(setMap, "_")
	if len(strMap) != 5 {
		t.Errorf("testTool fail")
	}

	//获取公网ip
	tool.GetPulicIP()

	//获取字符串 md5/sha1
	tool.MD5("a")
	tool.Sha1("a")

	//包含 Slice/Array/Map
	if !tool.Contains("1", []string{"1", "test"}) ||
		!tool.Contains("1", map[string]string{"1": "test"}) {
		t.Errorf("testTool fail")
	}

	ips, _ := tool.GetAllIps()
	ipv4, _ := tool.GetAllIpsV4()
	ipv6, _ := tool.GetAllIpsV6()
	macs, _ := tool.GetAllMacs()

	fmt.Println(ips)
	fmt.Println(ipv4)
	fmt.Println(ipv6)
	fmt.Println(macs)

	// 正则匹配
	marchIp := tool.Match("\\d+\\.\\d+\\.\\d+\\.\\d+", ipv4[0])
	fmt.Println(marchIp)

	//获取客户端ip
	request, _ := http.NewRequest("GET", "www.baidu.com", nil)
	request.Header.Add("X-Real-Ip", "127.0.0.2")

	tool.GetClientIp(request)

	// 转化为json
	jsonString := tool.JsonEncodePretty(map[string]string{"tst": "1"})
	fmt.Println(jsonString)

}

func TestPartEncode(t *testing.T) {

	// [Encode] string
	stringUrlEncode := tool.UrlEncode("test")
	fmt.Println(stringUrlEncode)

	// [Encode] string => Encode
	string1UrlEncode := tool.UrlEncode("c=hello&a=golang")
	fmt.Println(string1UrlEncode)

	// [Encode] map => 排序 => string
	urlMap := map[string]string{
		"c": "hello",
		"a": "golang",
	}
	fmt.Println(tool.UrlEncode(urlMap))

	// [Encode]  map => 排序 => sort => Encode
	fmt.Println(tool.UrlEncode(tool.UrlEncode(urlMap)))

	// [Encode]  map => 排序 => string
	fmt.Println(tool.UrlEncodeFromMap(urlMap))

	// [Decode] encode string => map
	urlDeMap := "c%3Dhello%26a%3Dgolang"
	fmt.Println(tool.UrlDecodeToMap(urlDeMap))

}
