package netlink

/*
#include "reg.h"
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
)

type NetlinkProtocol int

const (
	CONNECTOR NetlinkProtocol = syscall.NETLINK_CONNECTOR
	INET_DIAG NetlinkProtocol = syscall.NETLINK_INET_DIAG
)

type NetlinkSocket struct {
	fd       int
	protocol NetlinkProtocol
	groups   uint32
}

func NewNetlinkSocket(protocol NetlinkProtocol, groups uint32) (*NetlinkSocket, error) {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, int(protocol))
	if err != nil {
		return nil, err
	}
	ns := &NetlinkSocket{
		fd:       fd,
		protocol: protocol,
		groups:   groups,
	}
	if ns.protocol == CONNECTOR {
		sockaddr := &syscall.SockaddrNetlink{
			Family: syscall.AF_NETLINK,
			Pid:    uint32(os.Getpid()),
			Groups: groups,
		}
		if err = syscall.Bind(fd, sockaddr); err != nil {
			syscall.Close(fd)
			return nil, err
		}
	}

	if tmp := C.registerNetlink(C.int(ns.fd)); tmp != 0 {
		return nil, fmt.Errorf("can't register to netlink!")
	}
	fmt.Printf("Netlink Registered!\n")
	return ns, nil
}

func (ns *NetlinkSocket) Close() {
	syscall.Close(ns.fd)
}

func (ns *NetlinkSocket) Send(request *NetlinkRequest) error {
	if ns.protocol == CONNECTOR {
		if _, err := syscall.Write(ns.fd, request.Serialize()); err != nil {
			return err
		}
	} else {
		sockaddr := &syscall.SockaddrNetlink{
			Family: syscall.AF_NETLINK,
			Pid:    0,
			Groups: ns.groups,
		}
		if err := syscall.Sendto(ns.fd, request.Serialize(), 0, sockaddr); err != nil {
			return err
		}
	}
	return nil
}

func (s *NetlinkSocket) Receive(numPages int) ([]syscall.NetlinkMessage, error) {
	if numPages <= 0 {
		return nil, fmt.Errorf("too small buffer only %d pages", numPages)
	}
	rb := make([]byte, numPages*syscall.Getpagesize())
	nr, _, err := syscall.Recvfrom(s.fd, rb, 0)
	if err != nil {
		return nil, err
	}
	if nr < syscall.NLMSG_HDRLEN {
		return nil, fmt.Errorf("Got short response from netlink")
	}
	rb = rb[:nr]
	return syscall.ParseNetlinkMessage(rb)
}
