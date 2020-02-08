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

package parser

import (
	"bytes"
	"fmt"
	"github.com/bitofcode/hosts"
	"io"
	"sort"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	reader := bytes.NewBufferString(
		`127.0.0.1  localhost localhost.local
# a comment
192.168.10.10 example.com
192.168.10.10 example.io # another comment
192.168.15.15 hello.world# a dummy host`)

	entrySet, err := Read(reader)

	assertNoError(err, t)

	assertCorrectEntrySet(entrySet, t, err)

	expectedEntries := []hosts.Entry{
		hosts.NewEntry("127.0.0.1", []string{"localhost", "localhost.local"}),
		hosts.NewEntry("192.168.10.10", []string{"example.com", "example.io"}),
		hosts.NewEntry("192.168.15.15", []string{"hello.world"}),
	}

	assertEntries(expectedEntries, entrySet, t)
}

func assertEntries(expectedEntries []hosts.Entry, entrySet hosts.EntrySet, t *testing.T) {
	for _, expectedEntry := range expectedEntries {
		if !entrySet.Contains(expectedEntry) {
			t.Errorf("expected EntrySet(%v) contains %v", entrySet, expectedEntry)
		}
	}
}

func assertCorrectEntrySet(entrySet hosts.EntrySet, t *testing.T, err error) {
	if entrySet == nil {
		t.Fatalf("unexpected to get EntrySet=nil")
	}

	entries, err := entrySet.AllEntries()
	if len(entries) != 3 {
		t.Fatalf("unexpected number of entries %d", len(entries))
	}
}

func assertNoError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteFileWith(t *testing.T) {
	entrySet := hosts.NewEntrySet()
	entrySet.AddEntry(hosts.NewEntry("127.0.0.1", []string{"localhost", "localhost.local"}))
	entrySet.AddEntry(hosts.NewEntry("10.0.0.10", []string{"example.com", "example.io"}))

	buffer := bytes.NewBuffer(make([]byte, 0))

	err := WriteWith(entrySet, buffer, WriteToLine)

	if err != nil {
		t.Errorf("unexpected error '%v'", err)
	}

	content := buffer.String()

	entries, _ := entrySet.AllEntries()
	for _, entry := range entries {
		line, err := WriteToLine(entry)
		if !strings.Contains(content, line) || err != nil {
			t.Errorf("unexpected erro (%v) and expected '%s' to conatain '%s'", err, content, line)
		}
	}
}

func TestWriteFile(t *testing.T) {
	entrySet := hosts.NewEntrySet()
	entrySet.AddEntry(hosts.NewEntry("127.0.0.1", []string{"localhost", "localhost.local"}))
	entrySet.AddEntry(hosts.NewEntry("10.0.0.10", []string{"example.com", "example.io"}))
	entries, _ := entrySet.AllEntries()
	buffer := bytes.NewBuffer(make([]byte, 0))

	err := Write(entrySet, buffer)

	if err != nil {
		t.Errorf("unexpected error '%v'", err)
	}

	content := buffer.String()

	expectedContent := getExpectedContent(entries, t)
	if content != expectedContent {
		t.Fatalf("expected '%s' actual '%s'", expectedContent, content)
	}
}

func getExpectedContent(entries []hosts.Entry, t *testing.T) string {
	expectedBuffer := bytes.NewBuffer(make([]byte, 0))
	lines := make([]string, 0)
	for _, entry := range entries {
		line, err := WriteToLine(entry)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		lines = append(lines, line)
	}

	sort.Strings(lines)
	for _, l := range lines {
		_, err := io.WriteString(expectedBuffer, fmt.Sprintln(l))
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	expectedContent := expectedBuffer.String()
	return expectedContent
}
