package raft

import (
	"flag"
	"net"
	"strings"
)

type Config struct {
	Id       string
	Peers    []net.IP
	IsLeader bool
}

func ConfigFromFlags() Config {
	isLeader := flag.Bool("leader", false, "leader")
	peers := flag.String("peers", "", "peers")
	id := flag.String("id", "", "id")
	flag.Parse()

	peerStrs := strings.Split(*peers, ",")
	peerAddrs := make([]net.IP, len(peerStrs))
	for i, peerStr := range peerStrs {
		peerAddrs[i] = net.ParseIP(peerStr)
	}

	return Config{
		Id:       *id,
		Peers:    peerAddrs,
		IsLeader: *isLeader,
	}
}
