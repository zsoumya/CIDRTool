package net

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/zsoumya/smutils/regx"
)

func NewCIDRv4(ip *IPv4, blockSize byte) (*CIDRv4, error) {
	if blockSize > 32 {
		return nil, errors.New(fmt.Sprintf("Invalid CIDRv4 block size: %d", blockSize))
	}

	cidr := CIDRv4{
		BlockSize: blockSize,
		IPv4:      *ip,
	}

	return &cidr, nil
}

func ParseCIDRv4Str(data string) (*CIDRv4, error) {
	pattern := `^(?P<octet1>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet2>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet3>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(?P<octet4>25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])/(?P<blockSize>[0-9]|[1-2][0-9]|3[0-2])?$`

	match := regx.MatchNamedGroups(pattern, data)
	if match == nil {
		return nil, errors.New(fmt.Sprintf("Invalid CIDRv4: %s", data))
	}

	b, _ := strconv.ParseUint(match["blockSize"], 10, 8)

	cidr := CIDRv4{
		BlockSize: byte(b),
	}

	for i := 0; i < 4; i++ {
		b, _ := strconv.ParseUint(match[fmt.Sprintf("octet%d", i+1)], 10, 8)
		cidr.IPv4.octet[i] = byte(b)
	}

	return &cidr, nil
}
