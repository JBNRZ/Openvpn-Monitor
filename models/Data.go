package models

import (
	"regexp"
	"strconv"
	"strings"
)

func communicate(msg string) ([]byte, error) {
	_, err := Send(msg)
	if err != nil {
		Logger.Fatalln(err)
	}
	_, res, err := Receive()
	return res, err
}

func GetVersion() string {
	res, _ := communicate("version")
	reg := regexp.MustCompile(`OpenVPN Version: ([[:print:]]+) \[SSL.*?`)
	version := reg.FindStringSubmatch(string(res))
	if len(version) == 0 {
		reg = regexp.MustCompile(`TITLE,([[:print:]]+) \[SSL.*?`)
		version = reg.FindStringSubmatch(string(res))
	}
	if len(version) == 0 {
		return "OpenVPN"
	}
	return version[len(version)-1]
}

//func GetStats() (int64, int64, int64) {
//	res, _ := communicate("load-stats")
//	reg := regexp.MustCompile(`SUCCESS: nclients=(\d+),bytesin=(\d+),bytesout=(\d+)`)
//	data := reg.FindStringSubmatch(string(res))
//	if len(data) == 0 {
//		return 0, 0, 0
//	}
//	n, _ := strconv.ParseInt(data[1], 10, 64)
//	i, _ := strconv.ParseInt(data[2], 10, 64)
//	o, _ := strconv.ParseInt(data[3], 10, 64)
//	return n, i, o
//}

func GetDetail() []User {
	res, _ := communicate("status 2")
	data := string(res)
	lines := strings.Split(data, "\n")
	var clients []User
	for _, line := range lines {
		if strings.HasPrefix(line, "CLIENT_LIST,") {
			line := strings.Split(line, ",")
			client := new(User)
			client.Name = line[1]
			client.From = line[2]
			client.IPv4 = line[3]
			client.IPv6 = line[4]
			client.Status = true
			client.Received, _ = strconv.ParseInt(line[5], 10, 64)
			client.Sent, _ = strconv.ParseInt(line[6], 10, 64)
			client.TotalSent = 0
			client.TotalReceived = 0
			if client.Name == "UNDEF" {
				continue
			}
			clients = append(clients, *client)
		}
	}
	return clients
}
