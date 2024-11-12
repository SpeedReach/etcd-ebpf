package fast_ack

import (
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"net"
)

type Attachment struct {
	Link    link.Link
	Objects fast_ackObjects
}

func (a Attachment) Close() {
	defer a.Link.Close()
	defer a.Objects.Close()
}

func AttachEbpf() Attachment {
	if err := rlimit.RemoveMemlock(); err != nil {
		panic(err)
	}
	interfaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}

	for _, i := range interfaces {
		println("nic ", i.Name)
	}

	var objs fast_ackObjects
	if err := loadFast_ackObjects(&objs, nil); err != nil {
		panic(err)
	}
	ifname := "lo" // Change this to an interface on your machine.
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		panic(err)
	}

	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Interface: iface.Index,
		Program:   objs.FastReturn,
	})

	if err != nil {
		panic(err)
	}

	return Attachment{Link: xdpLink, Objects: objs}
}
