package stream

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/vmxy/go-ioi/ioi/util"
)

func connect(conn net.Conn, sorageId string) error {
	fmt.Println("s1 write")
	req := fmt.Sprintf("vmfs/1/%s", sorageId)
	conn.Write([]byte(req))
	bs := make([]byte, 128)
	size, err := conn.Read(bs)
	if err != nil {
		return err
	}
	res := util.String(bs[0:size])
	fmt.Println("res==========>", res)
	if !res.Match(`(?i)^vmfs/\d+ 200`) {
		return errors.New("connect scheme error " + res.String())
	}
	return nil
}

func IsWebSocket(input []byte) bool {
	headStr := string(input[:12])
	if ok, err := regexp.MatchString(`(?i)^GET `, headStr); !ok || err != nil {
		return false
	}

	str := string(input)
	headers := ParseHttpHeaders(str)

	return util.String(headers.Get("connection")).Match(`(?i)Upgrade`) &&
		util.String(headers.Get("upgrade")).Match(`(?i)websocket`) &&
		util.String(headers.Get("url")).Match(`(?i)^\/ws-`)
}

// vmfs/1
func IsVMFS(input []byte) bool {
	str := util.String(input[:6])
	return str.Match(`^vmfs\/[0-9]+`)
}

func parseVMFSRequest(input []byte) (version int, id string) {
	line := string(input)
	ps := strings.Split(line, "/")
	version = 1
	id = ""
	fmt.Println("parseVMFSRequest", line)
	switch len(ps) {
	case 2:
		return
	case 3:
		id = ps[2]
		return
	}
	return
}

func ParseHttpHeaders(input string) http.Header {
	str := util.String(input).Trim()
	lines := str.Split("\n") // str.split("\n").map((v) => v.trim());
	headers := http.Header{}
	for i, v := range lines {
		sv := util.String(v).Trim()
		if sv == "" {
			break
		}
		if i == 0 {
			hp := sv.Split(" ") //v.split(" ");
			headers.Set("method", string(hp.Get(0)))
			headers.Set("url", hp.Get(1).String())

		} else {
			kv := sv.Split(":")                     // v.split(":");
			key := kv.Get(0).ToLower().String()     //strings.ToLower(kv[0])
			val := util.String(":").Join(kv[1:]...) //strings.Join(kv[1:], ":")
			headers.Set(key, val.String())
		}
		//headers["protocol"] = headers["protocol"]
	}
	return headers
}
