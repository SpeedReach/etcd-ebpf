package raft

import "testing"

func TestEncodeDecode(t *testing.T) {
	var entry1 [64]byte
	str := "hello"
	copy(entry1[:], str)
	entry2 := [64]byte{}
	str = "world"
	copy(entry2[:], str)
	args := AppendEntriesArgs{
		Term:         1,
		LeaderId:     2,
		PrevLogIndex: 3,
		PrevLogTerm:  4,
		LeaderCommit: 5,
		EntryCount:   2,
		Entries: []LogEntry{
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

	data, err := args.ToBytes()
	if err != nil {
		t.Fatal(err)
	}
	decodedArgs, err := AppendEntriesArgsFromBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if decodedArgs.Term != args.Term {
		t.Fatalf("expected term 1, got %d", decodedArgs.Term)
	}
	if decodedArgs.LeaderId != args.LeaderId {
		t.Fatalf("expected leader id 2, got %d", decodedArgs.LeaderId)
	}

	if decodedArgs.PrevLogIndex != args.PrevLogIndex {
		t.Fatalf("expected prev log index 3, got %d", decodedArgs.PrevLogIndex)
	}
	if decodedArgs.PrevLogTerm != args.PrevLogTerm {
		t.Fatalf("expected prev log term 4, got %d", decodedArgs.PrevLogTerm)
	}

	if decodedArgs.LeaderCommit != args.LeaderCommit {
		t.Fatalf("expected leader commit 5, got %d", decodedArgs.LeaderCommit)
	}
	if decodedArgs.EntryCount != args.EntryCount {
		t.Fatalf("expected entry count 2, got %d", decodedArgs.EntryCount)
	}
	if len(decodedArgs.Entries) != len(args.Entries) {
		t.Fatalf("expected 2 entries, got %d", len(decodedArgs.Entries))
	}

	for i, entry := range decodedArgs.Entries {
		if entry.Term != args.Entries[i].Term {
			t.Fatalf("expected term %d, got %d", args.Entries[i].Term, entry.Term)
		}
		if string(entry.Command[:]) != string(args.Entries[i].Command[:]) {
			t.Fatalf("expected command %s, got %s", args.Entries[i].Command[:], entry.Command[:])
		}
	}

}
