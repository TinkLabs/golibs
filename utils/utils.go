package utils

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
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

	s := buf.String()
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\t", "\\t", -1)
	return s
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
		panic(err)
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
	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(time.RFC3339, data.(string))
		}

		return data, nil
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: stringToDateTimeHook,
		Result:     &target,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		panic(err)
	}

	return decoder.Decode(source)
}

func GetIntranetIp(dontCheckPrefix bool) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
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
