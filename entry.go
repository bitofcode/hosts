// MIT License
//
// Copyright (c) 2020 Wassim Akachi <wassim@bitofcode.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package hosts

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

var ErrorNilEntry = errors.New("entry is nil")
var ErrorInvalidIp = errors.New("invalid ip")
var ErrorInvalidHostName = errors.New("invalid host-name")

// An Entry represent a line in /etc/hosts with multiple hosts associate to one ip.
type Entry interface {
	Ip() net.IP
	IpString() string
	HostNames() []string
	AddHostName(hostName string) error
	String() string
	Contains(hostName string) bool
}

type simpleEntry struct {
	ip        net.IP
	hostNames map[string]bool
}

func (s *simpleEntry) Contains(hostName string) bool {
	_, ok := s.hostNames[hostName]
	return ok
}

func (s *simpleEntry) String() string {
	return fmt.Sprintf("%v", *s)
}

// NewEntry returns an Entry which contains the given ip and the host names.
func NewEntry(ip net.IP, hostNames []string) (Entry, error) {
	entry, err := NewEntryIp(ip)
	if err != nil {
		return nil, err
	}
	for _, h := range hostNames {
		err = entry.AddHostName(h)
		if err != nil {
			return nil, err
		}
	}
	return entry, nil
}

// NewEntryIp returns an Entry which contains the given ip.
func NewEntryIp(ip net.IP) (Entry, error) {
	if ip == nil {
		return nil, ErrorInvalidIp
	}
	return &simpleEntry{ip: ip, hostNames: make(map[string]bool)}, nil
}

func (s *simpleEntry) AddHostName(hostName string) error {
	if strings.Contains(hostName, " ") {
		return ErrorInvalidHostName
	}
	if strings.Contains(hostName, "#") {
		return ErrorInvalidHostName
	}
	s.hostNames[strings.ToLower(hostName)] = true
	return nil
}

func (s *simpleEntry) Ip() net.IP {
	return s.ip
}

func (s *simpleEntry) IpString() string {
	return s.ip.String()
}

func (s *simpleEntry) HostNames() []string {
	var hostNames []string

	if len(s.hostNames) > 0 {
		for host := range s.hostNames {
			hostNames = append(hostNames, host)
		}
		sort.Strings(hostNames)
	}
	return hostNames
}

func CloneEntry(entry Entry) (Entry, error) {
	if entry == nil {
		return nil, ErrorNilEntry
	}
	return NewEntry(entry.Ip(), entry.HostNames())
}

func NewEntryUnsafe(ip net.IP, hosts []string) Entry {
	h, err := NewEntry(ip, hosts)
	if err != nil {
		panic(err)
	}
	return h
}
