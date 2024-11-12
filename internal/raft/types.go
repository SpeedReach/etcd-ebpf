package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type LogEntry struct {
	Term    int32
	Command [64]byte
}

type AppendEntriesArgs struct {
	Term         int32
	LeaderId     int32
	PrevLogIndex int32
	PrevLogTerm  int32
	LeaderCommit int32
	EntryCount   int32
	Entries      []LogEntry
}

type AppendEntriesReply struct {
	PeerId  int32
	Term    int32
	Success bool
}

func (a AppendEntriesReply) String() string {
	return fmt.Sprintf("AppendEntriesReply{PeerId: %d, Term: %d, Success: %v}", a.PeerId, a.Term, a.Success)
}

func (a AppendEntriesArgs) String() string {
	return fmt.Sprintf("AppendEntriesArgs{Term: %d, LeaderId: %d, PrevLogIndex: %d, PrevLogTerm: %d, LeaderCommit: %d, EntryCount: %d, Entries: %v}",
		a.Term, a.LeaderId, a.PrevLogIndex, a.PrevLogTerm, a.LeaderCommit, a.EntryCount, a.Entries)
}

func (log LogEntry) String() string {
	return fmt.Sprintf("LogEntry{Term: %d, Command: %s}", log.Term, string(log.Command[:]))
}

func (a AppendEntriesArgs) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, a.Term); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.LeaderId); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.PrevLogIndex); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.PrevLogTerm); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.LeaderCommit); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.EntryCount); err != nil {
		panic(err)
	}
	for _, entry := range a.Entries {
		if err := binary.Write(&buf, binary.LittleEndian, entry.Term); err != nil {
			panic(err)
		}
		if err := binary.Write(&buf, binary.LittleEndian, entry.Command); err != nil {
			panic(err)
		}
	}
	return buf.Bytes(), nil
}

func AppendEntriesArgsFromBytes(buf io.Reader) (AppendEntriesArgs, error) {
	var args AppendEntriesArgs
	if err := binary.Read(buf, binary.LittleEndian, &args.Term); err != nil {
		return AppendEntriesArgs{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &args.LeaderId); err != nil {
		return AppendEntriesArgs{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &args.PrevLogIndex); err != nil {
		return AppendEntriesArgs{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &args.PrevLogTerm); err != nil {
		return AppendEntriesArgs{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &args.LeaderCommit); err != nil {
		return AppendEntriesArgs{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &args.EntryCount); err != nil {
		return AppendEntriesArgs{}, err
	}
	for i := 0; i < int(args.EntryCount); i++ {
		entry, _ := LogEntryFromBytes(buf)
		args.Entries = append(args.Entries, entry)
	}
	return args, nil
}

func LogEntryFromBytes(buf io.Reader) (LogEntry, error) {
	var entry LogEntry
	if err := binary.Read(buf, binary.LittleEndian, &entry.Term); err != nil {
		panic(err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &entry.Command); err != nil {
		panic(err)
	}
	return entry, nil
}

func (a AppendEntriesReply) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, a.PeerId); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.LittleEndian, a.Term); err != nil {
		panic(err)
	}

	//Align the buffer to 4 bytes
	if a.Success {
		if err := binary.Write(&buf, binary.LittleEndian, int32(1)); err != nil {
			panic(err)
		}
	} else {
		if err := binary.Write(&buf, binary.LittleEndian, int32(0)); err != nil {
			panic(err)
		}
	}

	return buf.Bytes(), nil
}

func AppendEntriesReplyFromBytes(buf io.Reader) (AppendEntriesReply, error) {
	var reply AppendEntriesReply
	if err := binary.Read(buf, binary.LittleEndian, &reply.PeerId); err != nil {
		return AppendEntriesReply{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &reply.Term); err != nil {
		return AppendEntriesReply{}, err
	}
	var success int32
	if err := binary.Read(buf, binary.LittleEndian, &success); err != nil {
		return AppendEntriesReply{}, err
	}
	if success == 1 {
		reply.Success = true
	} else {
		reply.Success = false
	}
	return reply, nil
}
