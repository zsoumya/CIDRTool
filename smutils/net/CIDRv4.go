package net

import (
	"errors"
	"fmt"
)

type CIDRv4 struct {
	IPv4
	BlockSize byte
}

func (cidr *CIDRv4) String() string {
	return fmt.Sprintf("%s/%d", cidr.NetworkIP().String(), cidr.BlockSize)
}

func (cidr *CIDRv4) Desc() string {
	return fmt.Sprintf("%18v [%15v => %15v](%d)", cidr, cidr.NetworkIP(), cidr.BroadcastIP(), cidr.Count())
}

func (cidr *CIDRv4) NetworkIP() *IPv4 {
	ip, _ := ParseIPv4Int(cidr.IPv4.Value() & cidr.NetMask().Value())
	return ip
}

func (cidr *CIDRv4) BroadcastIP() *IPv4 {
	ip, _ := ParseIPv4Int(cidr.NetworkIP().Value() | cidr.Wildcard().Value())
	return ip
}

func (cidr *CIDRv4) FirstIP() *IPv4 {
	if cidr.Count() > 2 {
		ip, _ := cidr.NetworkIP().Add(1)
		return ip
	} else {
		return nil
	}
}

func (cidr *CIDRv4) LastIP() *IPv4 {
	if cidr.Count() > 2 {
		ip, _ := cidr.BroadcastIP().Subtract(1)
		return ip
	} else {
		return nil
	}
}

func (cidr *CIDRv4) Count() uint32 {
	return cidr.BroadcastIP().Value() - cidr.NetworkIP().Value() + 1
}

func (cidr *CIDRv4) UsableCount() uint32 {
	if cidr.Count() < 3 {
		return 0
	}

	return cidr.Count() - 2
}

func (cidr *CIDRv4) NetMask() *IPv4 {
	ip, _ := ParseIPv4Int(0xffffffff << (32 - cidr.BlockSize))
	return ip
}

func (cidr *CIDRv4) Wildcard() *IPv4 {
	ip, _ := ParseIPv4Int((1 << (32 - cidr.BlockSize)) - 1)
	return ip
}

func (cidr *CIDRv4) ContainsIPv4(ip *IPv4) bool {
	return ip.Value() >= cidr.NetworkIP().Value() && ip.Value() <= cidr.BroadcastIP().Value()
}

func (cidr *CIDRv4) ContainsCIDRv4(child *CIDRv4) bool {
	return child.NetworkIP().Value() >= cidr.NetworkIP().Value() && child.BroadcastIP().Value() <= cidr.BroadcastIP().Value()
}

func (cidr *CIDRv4) Previous() *CIDRv4 {
	ip, err := cidr.NetworkIP().Subtract(cidr.Count())
	if err != nil {
		return nil
	}

	newCidr, _ := NewCIDRv4(ip, cidr.BlockSize)
	return newCidr
}

func (cidr *CIDRv4) Next() *CIDRv4 {
	ip, err := cidr.BroadcastIP().Add(1)
	if err != nil {
		return nil
	}

	newCidr, _ := NewCIDRv4(ip, cidr.BlockSize)
	return newCidr
}

func (cidr *CIDRv4) Partition(ip *IPv4) []*CIDRv4 {
	if !cidr.ContainsIPv4(ip) {
		return nil
	}

	cidrs := make([]*CIDRv4, 0)
	endIP := cidr.BroadcastIP()

	for {
		if ip.Value() > endIP.Value() {
			break
		}

		ipCount := endIP.Value() - ip.Value() + 1
		blkSize := blockSize(ipCount)
		if ipCount == 0 {
			cidr, _ := NewCIDRv4(NewIPv4([4]byte{0, 0, 0, 0}), 0)
			return []*CIDRv4{cidr}
		}

		startIP, _ := ParseIPv4Int(endIP.Value() - pow2(32-blkSize) + 1)
		childCidr, _ := NewCIDRv4(startIP, blkSize)

		cidrs = append([]*CIDRv4{childCidr}, cidrs...)
		if blkSize == 32 {
			break
		}

		endIP = childCidr.Previous().BroadcastIP()
	}

	return cidrs
}

func (cidr *CIDRv4) Split(splitCount uint) ([]*CIDRv4, error) {
	b, err := splitCountToBlockSize(splitCount)
	if err != nil {
		return nil, err
	}

	if cidr.BlockSize+b > 32 {
		return nil, errors.New(fmt.Sprintf("CIDR block %v not wide enough to be split into %d parts", cidr, splitCount))
	}

	startIP := cidr.NetworkIP()
	newBlockSize := cidr.BlockSize + b

	cidrs := make([]*CIDRv4, 0, 32)

	for i := uint(1); i <= splitCount; i++ {
		childCIDR, _ := NewCIDRv4(startIP, newBlockSize)
		cidrs = append(cidrs, childCIDR)

		startIP, _ = childCIDR.BroadcastIP().Add(1)
	}

	return cidrs, nil
}
