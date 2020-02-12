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

import "net"

type EntrySet interface {
	AddEntry(entry Entry, entries ...Entry)
	Contains(entry Entry) bool
	EntriesOfIP(ip net.IP) (hosts []string, ok bool)
	AllEntries() []Entry
}

type entrySet struct {
	entries map[string]Entry
}

func (e *entrySet) AddEntry(entry Entry, entries ...Entry) {
	e.addEntry(entry)
	for _, en := range entries {
		e.addEntry(en)
	}
}

func (e *entrySet) addEntry(entry Entry) {

	en := e.getOrCreateEntry(entry.Ip())
	for _, h := range entry.HostNames() {
		en.AddHostName(h)
	}
}

func (e *entrySet) Contains(entry Entry) bool {
	internalEntry, ok := e.entries[entry.IpString()]
	if !ok {
		return ok
	}

	for _, h := range entry.HostNames() {
		if !internalEntry.Contains(h) {
			return false
		}
	}
	return true
}

func (e *entrySet) EntriesOfIP(ip net.IP) (hosts []string, ok bool) {
	ent, ok := e.entries[ip.String()]
	if !ok {
		return nil, ok
	}
	hosts = make([]string, len(ent.HostNames()))
	copy(hosts, ent.HostNames())
	return hosts, true
}

func (e *entrySet) AllEntries() []Entry {
	entries := make([]Entry, 0)
	for _, ent := range e.entries {
		entry, _ := CloneEntry(ent)
		entries = append(entries, entry)
	}
	return entries
}

func (e *entrySet) getOrCreateEntry(ip net.IP) Entry {
	en, ok := e.entries[ip.String()]
	if !ok {
		en, _ = NewEntryIp(ip)
		e.entries[en.IpString()] = en
	}
	return en
}

func NewEntrySet() EntrySet {
	return &entrySet{
		entries: make(map[string]Entry),
	}
}
