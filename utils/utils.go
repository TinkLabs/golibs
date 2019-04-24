package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

var letterRunes = []rune("1234567890")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetNowTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func UUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func Quit() chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	return quit
}

func GenerateRandomString(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func ReadBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	newBuf := new(bytes.Buffer)
	if err := json.Compact(newBuf, buf.Bytes()); err != nil {
		return buf.String()
	}

	return newBuf.String()
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func GetPort() int {
	port, err := GetFreePort()
	if err != nil {
		panic(fmt.Sprintf("get free port:%v", err))
	}
	return port
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func SnameMapKey(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		m[snakeString(k)] = v
	}

	return m
}

func Decode(source map[string]interface{}, target interface{}) error {
	customHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(&time.Time{}) && f == reflect.TypeOf("") {
			rv, err := time.Parse(time.RFC3339, data.(string))
			if err != nil {
				return nil, err
			}

			return &rv, nil
		}
		if t == reflect.TypeOf(decimal.Decimal{}) && f == reflect.TypeOf(1.0) {
			return decimal.NewFromFloat(data.(float64)), nil
		}

		return data, nil
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: customHook,
		Result:     &target,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	return decoder.Decode(source)
}

func GetIntranetIp(dontCheckPrefix bool) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(fmt.Sprintf("get interfaces:%v", err))
	}

	fmt.Printf("interfaces = %+v\n", interfaces)
	for _, i := range interfaces {
		if dontCheckPrefix || strings.HasPrefix(i.Name, "eth") {
			itf, _ := net.InterfaceByName(i.Name)
			item, _ := itf.Addrs()
			for _, addr := range item {
				switch v := addr.(type) {
				case *net.IPNet:
					if !v.IP.IsLoopback() {
						if v.IP.To4() != nil {
							return v.IP.String()
						}
					}
				}
			}
		}
	}

	panic("cant find intranet ip")
}

func GetConsulAddressFromMetadata() string {
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		panic(fmt.Sprintf("get meta:%v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("read meta:%v", err))
	}

	return fmt.Sprintf("http://%s", string(body))
}
