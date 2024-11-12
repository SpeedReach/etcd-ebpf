package raft

import (
	"bytes"
	"github.com/SpeedReach/ebpf-etcd/internal/kernel/fast_ack"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/ringbuf"
)

type NodeState struct {
	CurrentTerm int
	VotedFor    int
	Logs        []LogEntry
}

type Node struct {
	State      NodeState
	Attachment fast_ack.Attachment
}

func (n Node) Close() {
	defer n.Attachment.Close()
}

func NewNode(currentTerm int32, votedFor int, logs []LogEntry) (Node, error) {
	attachment := fast_ack.AttachEbpf()
	if err := attachment.Objects.Term.Update(int32(0), currentTerm, ebpf.UpdateAny); err != nil {
		panic(err)
	}
	return Node{
		State: NodeState{
			int(currentTerm),
			votedFor,
			logs,
		},
		Attachment: attachment,
	}, nil
}

func (n Node) AppendEntries() {

}

func (n Node) HandleAppendEntries(args AppendEntriesArgs) {
	println(args.String())
}

func (n Node) Start() {
	argsChan := make(chan AppendEntriesArgs, 10)
	go func() {
		reader, err := ringbuf.NewReader(n.Attachment.Objects.NewEntries)
		if err != nil {
			panic(err)
		}
		for {
			record, err := reader.Read()
			if err != nil {
				panic(err)
			}
			data := record.RawSample
			args, err := AppendEntriesArgsFromBytes(bytes.NewBuffer(data))
			if err != nil {
				panic(err)
			}
			argsChan <- args
		}
	}()

	for {
		select {
		case args := <-argsChan:
			n.HandleAppendEntries(args)
		}
	}
}
