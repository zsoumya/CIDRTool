package net

import (
	"errors"
	"fmt"

	"strconv"

	"encoding/hex"

	"github.com/zsoumya/smutils/str"
)

type IPv4 struct {
	octet [4]byte
}

func (ip *IPv4) Value() uint32 {
	return uint32(ip.octet[0])<<24 + uint32(ip.octet[1])<<16 + uint32(ip.octet[2])<<8 + uint32(ip.octet[3])
}

func (ip *IPv4) String() string {
	s := ""

	for i := 0; i < 4; i++ {
		if i > 0 {
			s += "."
		}

		s += strconv.Itoa(int(ip.octet[i]))
	}

	return s
}

func (ip *IPv4) Format() string {
	s := ""

	for i := 0; i < 4; i++ {
		if i > 0 {
			s += "."
		}

		s += fmt.Sprintf("%03d", ip.octet[i])
	}

	return s
}

func (ip *IPv4) FormatHex(noDot bool) string {
	if noDot {
		return hex.EncodeToString(ip.octet[:])
	}

	s := ""

	for i := 0; i < 4; i++ {
		if i > 0 {
			s += "."
		}

		s += hex.EncodeToString([]byte{ip.octet[i]})
	}

	return s
}

func (ip *IPv4) Bin(sep string) string {
	s := ""

	for i := 0; i < 4; i++ {
		if i > 0 {
			s += sep
		}

		s += str.PadLeft(strconv.FormatInt(int64(ip.octet[i]), 2), "0", 8)
	}

	return s
}

func (ip *IPv4) Add(n uint32) (*IPv4, error) {
	ipInt := ip.Value()

	if ipInt+n < ipInt {
		return nil, errors.New("IPv4 range overflow")
	}

	newIP, _ := ParseIPv4Int(ipInt + n)
	return newIP, nil
}

func (ip *IPv4) Subtract(n uint32) (*IPv4, error) {
	ipInt := ip.Value()

	if ipInt-n > ipInt {
		return nil, errors.New("IPv4 range underflow")
	}

	newIP, _ := ParseIPv4Int(ip.Value() - n)
	return newIP, nil
}

func (ip *IPv4) CIDRSpan(endIP *IPv4) []*CIDRv4 {
	startIP := NewIPv4(ip.octet)
	if startIP.Value() > endIP.Value() {
		startIP, endIP = endIP, startIP
	}

	cidrs := make([]*CIDRv4, 0, 100)

	for {
		if startIP == nil || startIP.Value() > endIP.Value() {
			break
		}

		ipCount := endIP.Value() - startIP.Value() + 1
		if ipCount == 0 {
			cidr, _ := NewCIDRv4(NewIPv4([4]byte{0, 0, 0, 0}), 0)
			return []*CIDRv4{cidr}
		}

		blockSize := blockSize(ipCount)

		cidr, _ := NewCIDRv4(startIP, blockSize)
		partitionedCIDRs := cidr.Partition(startIP)

		for _, partitionedCIDR := range partitionedCIDRs {
			cidrs = append(cidrs, partitionedCIDR)
		}

		startIP, _ = cidr.BroadcastIP().Add(1)
	}

	return cidrs
}
