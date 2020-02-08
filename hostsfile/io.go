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

package hostsfile

import (
	"github.com/bitofcode/hosts"
	"github.com/bitofcode/hosts/parser"
	"os"
)

// Read read the content of the given path and parse it to an hosts.EntrySet.
func Read(path string) (entries hosts.EntrySet, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return parser.Read(file)
}

// Write writes the given hosts.EntrySet to the given path (create a new file if none exists).
func Write(entries hosts.EntrySet, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	return parser.Write(entries, file)
}
