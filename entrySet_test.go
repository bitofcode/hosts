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
	"testing"
)

func TestEntrySet_AddEntry(t *testing.T) {
	tests := []struct {
		entry Entry
	}{
		{
			entry: NewEntry("127.0.0.1", []string{"localhost"}),
		},
		{
			entry: NewEntry("127.0.0.1", []string{"localhost", "example.com"}),
		},
	}
	for _, test := range tests {
		t.Run(test.entry.String(), func(t *testing.T) {
			entries := NewEntrySet()
			entries.AddEntry(test.entry)
			if !entries.Contains(test.entry) {
				t.Errorf("expected EntrySet(%v) contains '%v'", entries, test.entry)
			}
		})
	}
}

func TestEntrySet_AllEntries(t *testing.T) {
	entries := NewEntrySet()
	entries.AddEntry(NewEntry("127.0.0.1", []string{"hallo", "hello"}))
	entries.AddEntry(NewEntry("127.0.0.1", []string{"example.de"}))
	entries.AddEntry(NewEntry("192.0.0.1", []string{"example.com"}))

	ents, err := entries.AllEntries()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ents) != 2 {
		t.Errorf("expected a list of length 2")
	}

	entriesMap := make(map[string]Entry)
	for _, entry := range ents {
		entriesMap[entry.Ip()] = entry
	}

	hostNames1 := []string{"hallo", "hello", "example.de"}
	ent := entriesMap["127.0.0.1"]
	for _, h := range hostNames1 {
		if !ent.Contains(h) {
			t.Errorf("expect Entry(%v) contains '%s'", ent, h)
		}
	}

	ent = entriesMap["192.0.0.1"]
	hostNames2 := []string{"example.com"}
	for _, h := range hostNames2 {
		if !ent.Contains(h) {
			t.Errorf("expect Entry(%v) contains '%s'", ent, h)
		}
	}
}

func TestEntrySet_EntriesOfIP(t *testing.T) {
	entries := NewEntrySet()
	entries.AddEntry(NewEntry("127.0.0.1", []string{"hallo", "hello"}))
	entries.AddEntry(NewEntry("127.0.0.1", []string{"example.de"}))
	entries.AddEntry(NewEntry("192.0.0.1", []string{"example.com"}))

	ents, err := entries.AllEntries()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(ents) != 2 {
		t.Errorf("expected a list of length 2")
	}

	expected := NewEntry("127.0.0.1", []string{"hallo", "hello", "example.de"})
	hosts, found := entries.EntriesOfIP("127.0.0.1")
	if !found {
		t.Errorf("expected to have '%v'", expected.String())
	}

	for _, h := range hosts {
		if !expected.Contains(h) {
			t.Errorf("expected %v.EntriesOfIP(127.0.0.1) to contain '%s'", entries, h)
		}
	}

}
