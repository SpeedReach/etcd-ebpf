package main

import (
	"bytes"
	"github.com/SpeedReach/ebpf-etcd/internal/raft"
	"net"
	"time"
)

func main() {
	//config := raft.ConfigFromFlags()
	go RunFollower()
	time.Sleep(1 * time.Second)
	RunLeader()
	time.Sleep(5 * time.Second)
}

func RunFollower() {
	node1, err := raft.NewNode(1, 0, nil)
	if err != nil {
		panic(err)
	}
	defer node1.Close()
	node1.Start()
}

func RunLeader() {
	udp, err := net.DialUDP("udp", nil, &net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}
	args := testArg()
	data, err := args.ToBytes()
	if err != nil {
		panic(err)
	}
	println("Size of data: ", len(data))
	response := raft.AppendEntriesReply{}

	resData, err := response.ToBytes()
	if err != nil {
		panic(err)
	}
	println(len(resData))

	data = append(resData, data...)

	_, err = udp.Write(data)

	resBuf := make([]byte, 1024)
	_, addr, err := udp.ReadFromUDP(resBuf)
	if err != nil {
		panic(err)
	}
	println("Received from: ", addr)
	res, err := raft.AppendEntriesReplyFromBytes(bytes.NewBuffer(resBuf))
	if err != nil {
		panic(err)
	}
	println("www", res.String())
}

func testArg() raft.AppendEntriesArgs {
	var entry1 [64]byte
	str := "hello"
	copy(entry1[:], str)
	entry2 := [64]byte{}
	str = "world"
	copy(entry2[:], str)
	return raft.AppendEntriesArgs{
		Term:         1,
		LeaderId:     2,
		PrevLogIndex: 3,
		PrevLogTerm:  4,
		LeaderCommit: 5,
		EntryCount:   2,
		Entries: []raft.LogEntry{
			{
				Term:    1,
				Command: entry1,
			},
			{
				Term:    2,
				Command: entry2,
			},
		},
	}
}
