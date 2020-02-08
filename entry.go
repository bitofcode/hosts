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
	"sort"
)

var ErrorNilEntry = errors.New("entry is nil")

// An Entry represent a line in /etc/hosts with multiple hosts associate to one ip.
type Entry interface {
	Ip() string
	HostNames() []string
	AddHostName(hostName string)
	String() string
	Contains(hostName string) bool
}

type simpleEntry struct {
	ip        string
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
func NewEntry(ip string, hostNames []string) Entry {
	entry := NewEntryIp(ip)
	for _, h := range hostNames {
		entry.AddHostName(h)
	}
	return entry
}

// NewEntry returns an Entry which contains the given ip.
func NewEntryIp(ip string) Entry {
	return &simpleEntry{ip: ip, hostNames: make(map[string]bool)}
}

func (s *simpleEntry) AddHostName(hostName string) {
	s.hostNames[hostName] = true
}

func (s *simpleEntry) Ip() string {
	return s.ip
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
	return NewEntry(entry.Ip(), entry.HostNames()), nil
}
