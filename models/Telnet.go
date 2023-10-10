package models

import (
	"fmt"
	"net"
	"strings"
)

var conn *net.TCPConn

func InitConn() {
	server := fmt.Sprintf("%s:%s", Env.GetString("openvpn.ip"), Env.GetString("openvpn.port"))

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)

	if err != nil {
		Logger.Fatalln(err)
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		Logger.Fatalln(err)
	}

	//接收初次消息
	_, _, _ = Receive()
}

func Send(msg string) (int, error) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	sent, err := conn.Write([]byte(msg))
	if err != nil {
		Logger.Fatalln(err)
		return 0, err
	}
	if sent != len(msg) {
		Logger.Warning("Expected send length ", len(msg), ", but fact length ", sent)
	}
	return sent, nil
}

func Receive() (int, []byte, error) {
	buf := make([]byte, 2048)
	res, err := conn.Read(buf)
	if err != nil {
		Logger.Warning(err)
		return 0, nil, err
	}
	return res, buf, nil
}
