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

import "testing"

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		rawLine    string
		wantedLine string
	}{
		{rawLine: "", wantedLine: ""},
		{rawLine: "a", wantedLine: "a"},
		{rawLine: " a", wantedLine: "a"},
		{rawLine: " a ", wantedLine: "a"},
		{rawLine: "a ", wantedLine: "a"},
		{rawLine: "\ta", wantedLine: "a"},
		{rawLine: "a\t", wantedLine: "a"},
		{rawLine: " hallo ", wantedLine: "hallo"},
		{rawLine: " hallo \t", wantedLine: "hallo"},
		{rawLine: "\t hallo \t", wantedLine: "hallo"},
		{rawLine: "\thallo", wantedLine: "hallo"},
		{rawLine: "hallo", wantedLine: "hallo"},
		{rawLine: "       hallo      ", wantedLine: "hallo"},
		{rawLine: "       hallo      \n\n\t", wantedLine: "hallo"},
	}
	for _, test := range tests {
		t.Run(test.rawLine, func(t *testing.T) {
			trimmed := TrimWhitespace(test.rawLine)
			if trimmed != test.wantedLine {
				t.Errorf("wanted: TrimWhitespace('%s')='%s', but actual: '%s'", test.rawLine, test.wantedLine, trimmed)
			}
		})
	}
}
