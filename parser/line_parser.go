// MIT License
//
// Copyright (c) 2020 Bit of Code
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

// Parse entries/lines from the /etc/hosts file.

package parser

import (
	"errors"
	"fmt"
	"github.com/bitofcode/hosts"
	"regexp"
	"strings"
)

var (
	spaceRegexp         = regexp.MustCompile(`\s+`)
	emptyLineError      = errors.New("empty line")
	InvalidLineError    = errors.New("invalid line")
	invalidIp           = errors.New("invalid ip")
	invalidHostName     = errors.New("invalid host name")
	invalidHostNameList = errors.New("invalid host name list")
)

const commentSign = "#"

// ReadFromLine convert a given string to hostsfile.Entry.
func ReadFromLine(line string) (ent hosts.Entry, err error) {
	if isEmptyOrComment(line) {
		return nil, emptyLineError
	}

	trimmedLine := strings.Trim(line, " ")
	if isEmptyOrComment(trimmedLine) {
		return nil, emptyLineError
	}

	commentFreeLine := extractCommentFreeLine(trimmedLine)
	if isEmptyOrComment(commentFreeLine) {
		return nil, emptyLineError
	}

	lineItems := spaceRegexp.Split(commentFreeLine, -1)
	if len(lineItems) <= 1 {
		return nil, InvalidLineError
	}

	return hosts.NewEntry(lineItems[0], lineItems[1:]), nil
}

func extractCommentFreeLine(trimmedLine string) string {
	commentSignPosition := strings.Index(trimmedLine, commentSign)
	commentFreeLine := trimmedLine

	if commentSignPosition == 0 {
		commentFreeLine = ""
	}
	if commentSignPosition > 0 {
		commentFreeLine = trimmedLine[:commentSignPosition]
	}

	return strings.Trim(commentFreeLine, " ")
}

func isEmptyOrComment(line string) bool {
	return len(line) <= 0 || strings.HasPrefix(line, commentSign)
}

// WriteToLine converts a given hostsfile.Entry to etc/hosts line without line-separator.
func WriteToLine(ent hosts.Entry) (line string, err error) {
	if strings.Contains(ent.Ip(), commentSign) {
		return "", invalidIp
	}
	trimmerIp := strings.Trim(ent.Ip(), " ")
	if len(trimmerIp) <= 0 {
		return "", invalidIp
	}

	if ent.HostNames() == nil || len(ent.HostNames()) <= 0 {
		return "", invalidHostNameList
	}

	for _, hostName := range ent.HostNames() {
		if strings.Contains(hostName, commentSign) {
			return "", invalidHostName
		}
	}

	return fmt.Sprintf("%s  %s", ent.Ip(), strings.Join(ent.HostNames(), "  ")), nil
}
