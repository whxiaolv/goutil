package goutil

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"io"
	random "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Tool struct {
}

func (this *Tool) GetUUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	id := this.MD5(base64.URLEncoding.EncodeToString(b))
	return fmt.Sprintf("%s-%s-%s-%s-%s", id[0:8], id[8:12], id[12:16], id[16:20], id[20:])
}
func (this *Tool) CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
func (this *Tool) RandInt(min, max int) int {
	return func(min, max int) int {
		r := random.New(random.NewSource(time.Now().UnixNano()))
		if min >= max {
			return max
		}
		return r.Intn(max-min) + min
	}(min, max)
}
func (this *Tool) GetToDay() string {
	return time.Now().Format("20060102")
}
func (this *Tool) UrlEncode(v interface{}) string {
	switch v.(type) {
	case string:
		m := make(map[string]string)
		m["name"] = v.(string)
		return strings.Replace(this.UrlEncodeFromMap(m), "name=", "", 1)
	case map[string]string:
		return this.UrlEncodeFromMap(v.(map[string]string))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (this *Tool) UrlEncodeFromMap(m map[string]string) string {
	vv := url.Values{}
	for k, v := range m {
		vv.Add(k, v)
	}
	return vv.Encode()
}
func (this *Tool) UrlDecodeToMap(body string) (map[string]string, error) {
	var (
		err error
		m   map[string]string
		v   url.Values
	)
	m = make(map[string]string)
	if v, err = url.ParseQuery(body); err != nil {
		return m, err
	}
	for _k, _v := range v {
		if len(_v) > 0 {
			m[_k] = _v[0]
		}
	}
	return m, nil
}
func (this *Tool) GetDayFromTimeStamp(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format("20060102")
}
func (this *Tool) StrToMapSet(str string, sep string) mapset.Set {
	result := mapset.NewSet()
	for _, v := range strings.Split(str, sep) {
		result.Add(v)
	}
	return result
}
func (this *Tool) MapSetToStr(set mapset.Set, sep string) string {
	var (
		ret []string
	)
	for v := range set.Iter() {
		ret = append(ret, v.(string))
	}
	return strings.Join(ret, sep)
}
func (this *Tool) GetPulicIP() string {
	var (
		err  error
		conn net.Conn
	)
	if conn, err = net.Dial("udp", "8.8.8.8:80"); err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
func (this *Tool) MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return fmt.Sprintf("%x", md.Sum(nil))
}

func (this *Tool) Sha1(str string) string {
	md5h := sha1.New()
	md5h.Write([]byte(str))
	return fmt.Sprintf("%x", md5h.Sum(nil))
}

func (this *Tool) Contains(obj interface{}, arrayobj interface{}) bool {
	targetValue := reflect.ValueOf(arrayobj)
	kind := reflect.TypeOf(arrayobj).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func (this *Tool) Match(matcher string, content string) []string {
	var result []string
	if reg, err := regexp.Compile(matcher); err == nil {
		result = reg.FindAllString(content, -1)
	}
	return result
}

func (this *Tool) JsonEncodePretty(o interface{}) string {
	resp := ""
	switch o.(type) {
	case map[string]interface{}:
		if data, err := json.Marshal(o); err == nil {
			resp = string(data)
		}
	case map[string]string:
		if data, err := json.Marshal(o); err == nil {
			resp = string(data)
		}
	case []interface{}:
		if data, err := json.Marshal(o); err == nil {
			resp = string(data)
		}
	case []string:
		if data, err := json.Marshal(o); err == nil {
			resp = string(data)
		}
	case string:
		resp = o.(string)
	default:
		if data, err := json.Marshal(o); err == nil {
			resp = string(data)
		}
	}
	var v interface{}
	if ok := json.Unmarshal([]byte(resp), &v); ok == nil {
		if buf, ok := json.MarshalIndent(v, "", "  "); ok == nil {
			resp = string(buf)
		}
	}
	return resp
}
func (this *Tool) GetClientIp(r *http.Request) string {
	client_ip := ""
	headers := []string{"X_Forwarded_For", "X-Forwarded-For", "X-Real-Ip",
		"X_Real_Ip", "Remote_Addr", "Remote-Addr"}
	for _, v := range headers {
		if _v, ok := r.Header[v]; ok {
			if len(_v) > 0 {
				client_ip = _v[0]
				break
			}
		}
	}
	if client_ip == "" {
		clients := strings.Split(r.RemoteAddr, ":")
		client_ip = clients[0]
	}
	return client_ip
}

func (this *Tool) GetAllIpsV4() ([]string, error) {
	var (
		ips    []string
		allIps []string
		err    error
	)
	if allIps, err = this.GetAllIps(); err != nil {
		return ips, err
	}
	for _, ip := range allIps {
		i := this.Match("\\d+\\.\\d+\\.\\d+\\.\\d+", ip)
		if len(i) > 0 {
			ips = append(ips, i[0])
		}
	}
	return ips, nil
}

func (this *Tool) GetAllIpsV6() ([]string, error) {
	var (
		ips    []string
		allIps []string
		err    error
	)
	if allIps, err = this.GetAllIps(); err != nil {
		return ips, err
	}
	for _, ip := range allIps {
		i := this.Match("[0-9a-z:]{15,}", ip)
		if len(i) > 0 {
			ips = append(ips, i[0])
		}
	}
	return ips, nil
}

func (this *Tool) GetAllIps() ([]string, error) {
	var (
		err        error
		interfaces []net.Interface
		ips        []string
		addrs      []net.Addr
	)
	if interfaces, err = net.Interfaces(); err != nil {
		return ips, err
	}
	for _, v := range interfaces {
		if addrs, err = v.Addrs(); err == nil {
			for _, addr := range addrs {
				ips = append(ips, addr.String())
			}
		}
	}
	return ips, nil
}

func (this *Tool) GetAllMacs() ([]string, error) {
	var (
		err        error
		interfaces []net.Interface
		macs       []string
	)
	if interfaces, err = net.Interfaces(); err != nil {
		return macs, err
	}
	for _, v := range interfaces {
		macs = append(macs, v.HardwareAddr.String())
	}
	return macs, nil
}
