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
	"bufio"
	"fmt"
	"github.com/bitofcode/hosts"
	"io"
	"sort"
)

// Read reads the hosts file from the provided io.Reader and returns an EntrySet.
func Read(reader io.Reader) (entrySet hosts.EntrySet, err error) {
	entrySet = hosts.NewEntrySet()
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		entry, err := ReadFromLine(line)
		if err != nil && err == InvalidLineError {
			return nil, err
		}
		if err == nil {
			entrySet.AddEntry(entry)
		}
	}
	return entrySet, nil
}

// Write writes the EntrySet to a well formatted hosts file into io.Write.
func Write(entrySet hosts.EntrySet, writer io.Writer) error {
	return WriteWith(entrySet, writer, WriteToLine)
}

// WriteWith writes all entries formatted with the provided formatter from the provided EntrySet into io.Write.
func WriteWith(entrySet hosts.EntrySet, writer io.Writer, formatter func(ent hosts.Entry) (line string, err error)) error {
	entries := entrySet.AllEntries()

	lines := make([]string, 0)
	for _, entry := range entries {
		line, err := formatter(entry)
		if err != nil {
			return err
		}
		lines = append(lines, line)
	}

	sort.Strings(lines)
	for _, l := range lines {
		_, err := io.WriteString(writer, fmt.Sprintln(l))
		if err != nil {
			return err
		}
	}
	return nil
}
