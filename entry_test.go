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
	"net"
	"strings"
	"testing"
)

func TestNewEntryIpWithNilIp(t *testing.T) {
	_, err := NewEntryIp(nil)
	if err != ErrorInvalidIp {
		t.Errorf("expected error= %v, actual: %v", ErrorInvalidIp, err)
	}
}

func TestNewEntryWithNilIp(t *testing.T) {
	_, err := NewEntry(nil, []string{"localhost"})
	if err != ErrorInvalidIp {
		t.Errorf("expected error= %v, actual: %v", ErrorInvalidIp, err)
	}
}

func TestCloneEntry(t *testing.T) {
	_, err := CloneEntry(nil)
	if err != ErrorNilEntry {
		t.Errorf("expected error= %v, actual: %v", ErrorNilEntry, err)
	}
}

func TestNewEntryUnsafePanics(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrorInvalidIp {
			t.Errorf("expected error: '%#v', but actual: '%#v'", ErrorInvalidIp, err)
		}
	}()

	NewEntryUnsafe(nil, []string{"localhost"})
}

func TestNewEntryUnsafe(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("unexpected error: '%#v'", err)
		}
	}()

	NewEntryUnsafe(net.ParseIP("127.0.0.1"), []string{"localhost"})
}

func TestSimpleEntry_AddHostName(t *testing.T) {
	name := "localhost"
	ent, _ := NewEntryIp(net.ParseIP("127.0.0.1"))
	ent.AddHostName(name)
	ent.AddHostName(name)
	ent.AddHostName(strings.ToUpper(name))

	found := false
	for _, hostName := range ent.HostNames() {
		isIn := hostName == name
		if found && isIn {
			t.Errorf("Unexpected duplicate of '%s' in %v", name, ent.HostNames())
		}
		found = found || isIn
	}

	if !found {
		t.Errorf("AddHostName want '%v', got = '%v'", []string{name}, ent.HostNames())
	}
}

func TestSimpleEntry_AddHostNameInvalid(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "local host"},
		{name: "local#host"},
		{name: "localhost#"},
		{name: "#localost"},
	}
	for _, test := range tests {
		entry := NewEntryUnsafe(net.ParseIP("127.0.0.1"), make([]string, 0))
		t.Run(test.name, func(t *testing.T) {
			err := entry.AddHostName(test.name)
			if err != ErrorInvalidHostName {
				t.Errorf("expected error: '%#v', but actual: '%#v'", ErrorInvalidHostName, err)
			}
		})
	}
}
