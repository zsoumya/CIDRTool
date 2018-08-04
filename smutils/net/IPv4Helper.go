package net

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/zsoumya/smutils/regx"
)

func NewIPv4(octet [4]byte) *IPv4 {
	ip := IPv4{}

	for i := 0; i < 4; i++ {
		ip.octet[i] = octet[i]
	}

	return &ip
}

func ParseIPv4Str(data string) (*IPv4, error) {
	const pattern = `^(?P<octet1>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet2>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet3>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet4>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`

	match := regx.MatchNamedGroups(pattern, data)
	if match == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IPv4: %s", data))
	}

	ip := IPv4{}

	for i := 0; i < 4; i++ {
		octet := match[fmt.Sprintf("octet%d", i+1)]
		b, _ := strconv.ParseUint(octet, 10, 8)
		ip.octet[i] = byte(b)
	}

	return &ip, nil
}

func ParseIPv4Bin(data string) (*IPv4, error) {
	const pattern = `^(?P<octstr1>[01]{8})( |.|-)?(?P<octstr2>[01]{8})( |.|-)?(?P<octstr3>[01]{8})( |.|-)?(?P<octstr4>[01]{8})$`

	match := regx.MatchNamedGroups(pattern, data)
	if match == nil {
		return nil, errors.New("invalid IPv4 format")
	}

	ip := IPv4{}

	for i := 0; i < 4; i++ {
		octstr := match[fmt.Sprintf("octstr%d", i+1)]
		b, _ := strconv.ParseUint(octstr, 2, 8)
		ip.octet[i] = byte(b)
	}

	return &ip, nil
}

func ParseIPv4Int(data uint32) (*IPv4, error) {
	ip := IPv4{}

	ip.octet[0] = byte(data >> 24)
	ip.octet[1] = byte((data & 0xff0000) >> 16)
	ip.octet[2] = byte((data & 0xff00) >> 8)
	ip.octet[3] = byte(data & 0xff)

	return &ip, nil
}
