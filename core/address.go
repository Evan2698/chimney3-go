package core

import (
	"bytes"
	"net"
	"strconv"
)

const (
	// Address types used by SOCKS5
	ADDRESSTYPE_IPV4   byte = 0x01
	ADDRESSTYPE_DOMAIN byte = 0x03
	ADDRESSTYPE_IPV6   byte = 0x04
)

type Socks5Address struct {
	Type  byte
	Host  string
	Port  uint16
	IP    []byte
	IPV6  []byte
	Valid bool
}

func NewSocks5Address() *Socks5Address {
	return &Socks5Address{
		Valid: false,
	}
}

func (addr *Socks5Address) GetAddress() string {
	if addr == nil {
		return ""
	}
	if addr.Type == ADDRESSTYPE_DOMAIN {
		return addr.Host
	}
	if addr.Type == ADDRESSTYPE_IPV4 {
		return net.IP(addr.IP).String()
	}
	if addr.Type == ADDRESSTYPE_IPV6 {
		return net.IP(addr.IPV6).String()
	}
	return ""
}

func (addr *Socks5Address) GetPort() uint16 {
	if addr == nil {
		return 0
	}
	return addr.Port
}

func (addr *Socks5Address) String() string {
	if addr == nil {
		return ""
	}
	return net.JoinHostPort(addr.GetAddress(), strconv.Itoa(int(addr.GetPort())))
}

func (addr *Socks5Address) SetDomainAddress(host string, port uint16) {
	if addr == nil {
		return
	}
	addr.Host = host
	addr.Port = port
	addr.Type = ADDRESSTYPE_DOMAIN
	addr.Valid = true
}

func (addr *Socks5Address) SetIPv4Address(ip []byte, port uint16) {
	if addr == nil {
		return
	}
	if ip4 := net.IP(ip).To4(); ip4 != nil {
		addr.IP = ip4
	} else {
		addr.IP = append([]byte(nil), ip...)
	}
	addr.Port = port
	addr.Type = ADDRESSTYPE_IPV4
	addr.Valid = true
}

func (addr *Socks5Address) SetIPv6Address(ip []byte, port uint16) {
	if addr == nil {
		return
	}
	if ip16 := net.IP(ip).To16(); ip16 != nil {
		addr.IPV6 = ip16
	} else {
		addr.IPV6 = append([]byte(nil), ip...)
	}
	addr.Port = port
	addr.Type = ADDRESSTYPE_IPV6
	addr.Valid = true
}

func (addr *Socks5Address) Bytes() []byte {
	if addr == nil {
		return nil
	}
	var hello bytes.Buffer
	hello.WriteByte(addr.Type)
	switch addr.Type {
	case ADDRESSTYPE_DOMAIN:
		hello.WriteByte(byte(len(addr.Host)))
		hello.WriteString(addr.Host)
	case ADDRESSTYPE_IPV4:
		hello.Write(addr.IP)
	case ADDRESSTYPE_IPV6:
		hello.Write(addr.IPV6)
	}
	hello.WriteByte(byte(addr.Port >> 8))
	hello.WriteByte(byte(addr.Port & 0xff))
	return hello.Bytes()
}

func (addr *Socks5Address) Parse(data []byte) error {
	if addr == nil {
		return nil
	}
	if len(data) < 7 {
		addr.Valid = false
		return nil
	}
	addr.Type = data[0]
	if addr.Type == ADDRESSTYPE_DOMAIN {
		addr.Host = string(data[2 : 2+data[1]])
		addr.Port = uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])
	} else if addr.Type == ADDRESSTYPE_IPV4 {
		addr.IP = data[1:5]
		addr.Port = uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])
	} else if addr.Type == ADDRESSTYPE_IPV6 {
		addr.IPV6 = data[1:17]
		addr.Port = uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])
	} else {
		addr.Valid = false
		return nil
	}
	addr.Valid = true
	return nil
}
func ParseTargetAddress(host string) (*Socks5Address, error) {
	s, p, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}

	np, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}

	v := &Socks5Address{}
	ip := net.ParseIP(s)
	if ip == nil {
		v.SetDomainAddress(s, uint16(np))
		return v, nil
	}
	if ip.To4() == nil {
		v.SetIPv6Address(ip, uint16(np))
		return v, nil
	}
	v.SetIPv4Address(ip, uint16(np))
	return v, nil
}
