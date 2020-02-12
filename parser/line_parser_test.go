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
	"fmt"
	"github.com/bitofcode/hosts"
	"sort"
	"testing"
)

func TestParseFromLineEmptyLine(t *testing.T) {
	tests := []struct {
		line string
	}{
		{line: "# "},
		{line: " # "},
		{line: ""},
		{line: "    "},
		{line: "    #"},
		{line: "    #     "},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			_, err := ReadFromLine(test.line)
			if err == nil || err != emptyLineError {
				t.Errorf("For line '%s' an error is expected: '%v', actual: %v", test.line, emptyLineError, err)
			}
		})
	}
}

func TestParseFromLineInvalidLine(t *testing.T) {
	tests := []struct {
		line string
	}{
		{line: "123# "},
		{line: " 123# "},
		{line: "123"},
		{line: "    123"},
		{line: "    234#"},
		{line: "    234#     "},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			_, err := ReadFromLine(test.line)
			if err == nil || err != InvalidLineError {
				t.Errorf("For line '%s' an error is expected: '%v', actual: %v", test.line, emptyLineError, err)
			}
		})
	}
}
func TestParseFromLineValidLine(t *testing.T) {
	tests := []struct {
		line      string
		ip        string
		hostNames []string
	}{
		{
			line:      "127.0.0.1 localhost",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost#",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost #",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost # ",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost# ",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost#ignore me",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost# ignore me",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		}, {
			line:      "127.0.0.1 localhost # ignore me",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost"},
		},
		{
			line:      "127.0.0.1 localhost example.com",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost", "example.com"},
		},
		{
			line:      "127.0.0.1 localhost example.com # ignore",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost", "example.com"},
		},
		{
			line:      "       127.0.0.1        localhost       example.com # ignore",
			ip:        "127.0.0.1",
			hostNames: []string{"localhost", "example.com"},
		},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			ent, err := ReadFromLine(test.line)
			if err != nil {
				t.Errorf("For line '%s' an unexpected error occurred: '%v'", test.line, err)
			}
			ip := ent.Ip()
			if ip != test.ip {
				t.Errorf("For line '%s', expected ip='%s', actual='%s'", test.line, test.ip, ip)
			}

			hostNames := ent.HostNames()
			if hostNames == nil || len(hostNames) != len(test.hostNames) {
				t.Fatalf("For line '%s', expected number of hosts='%d', actual='%d'",
					test.line, len(test.hostNames), len(hostNames))
			}
			sort.Strings(test.hostNames)
			for index, hostName := range test.hostNames {
				if hostName != hostNames[index] {
					t.Errorf("For line '%s', test.hostNames[%d]='%s', actual hostNames[%d]='%s'",
						test.line, index, hostName, index, hostNames[index])
				}
			}

		})
	}
}

func TestParseToLine(t *testing.T) {
	type args struct {
		ip        string
		hostNames []string
	}

	tests := []struct {
		args     args
		wantLine string
		wantErr  bool
		err      error
	}{
		{
			args: args{
				ip:        "127.0.0.1",
				hostNames: nil,
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidHostNameList,
		}, {

			args: args{
				ip:        "127.0.0.1",
				hostNames: []string{},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidHostNameList,
		},
		{
			args: args{
				ip:        "127.0.0.1",
				hostNames: []string{"#"},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidHostName,
		},
		{
			args: args{
				ip:        "127.0.0.1",
				hostNames: []string{"# hallo"},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidHostName,
		},
		{
			args: args{
				ip:        "127.0.0.1#",
				hostNames: []string{"# hallo"},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidIp,
		},
		{
			args: args{
				ip:        "#127.0.0.1#",
				hostNames: []string{"# hallo"},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidIp,
		},
		{
			args: args{
				ip:        "#127.0.0.1",
				hostNames: []string{"# hallo"},
			},
			wantLine: "",
			wantErr:  true,
			err:      invalidIp,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("'%v'=>'%s' (error='%v')", tt.args, tt.wantLine, tt.err), func(t *testing.T) {

			gotLine, err := WriteToLine(hosts.NewEntry(tt.args.ip, tt.args.hostNames))

			if tt.wantErr && (tt.err != err) {
				t.Errorf("WriteToLine(%v) expected error = '%v', actual = '%v'", tt.args, tt.err, err)
				return
			}

			if !tt.wantErr && (err != nil) {
				t.Errorf("WriteToLine(%v) unexpected error = %v", tt.args, err)
				return
			}

			if gotLine != tt.wantLine {
				t.Errorf("WriteToLine(%v) want %v, gotLine = %v", tt.args, tt.wantLine, gotLine)
			}
		})
	}
}

func TestParseToLineValidEntry(t *testing.T) {
	type args struct {
		ip        string
		hostNames []string
	}

	tests := []struct {
		args     args
		wantLine string
		wantErr  bool
		err      error
	}{
		{
			args: args{
				ip:        "127.0.0.1",
				hostNames: []string{"localhost"},
			},
			wantLine: "127.0.0.1  localhost",
			wantErr:  false,
		}, {
			args: args{
				ip:        "127.0.0.1",
				hostNames: []string{"localhost", "hello.world"},
			},
			wantLine: "127.0.0.1  hello.world  localhost",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("'%v'=>'%s' (error='%v')", tt.args, tt.wantLine, tt.err), func(t *testing.T) {

			gotLine, err := WriteToLine(hosts.NewEntry(tt.args.ip, tt.args.hostNames))

			if tt.wantErr && (tt.err != err) {
				t.Errorf("WriteToLine(%v) expected error = '%v', actual = '%v'", tt.args, tt.err, err)
				return
			}

			if !tt.wantErr && (err != nil) {
				t.Errorf("WriteToLine(%v) unexpected error = %v", tt.args, err)
				return
			}

			if gotLine != tt.wantLine {
				t.Errorf("WriteToLine(%v) want '%v', gotLine = '%v'", tt.args, tt.wantLine, gotLine)
			}
		})
	}
}

func TestLineRegexp(t *testing.T) {
	line := "127.0.0.1 localhost host.local"
	findings := regexpTrimWhitespace.FindStringSubmatch(line)
	t.Logf("Len(findings)=%d\n", len(findings))
	t.Logf("%#v\n", findings)
	t.Logf("%#v\n", regexpTrimWhitespace.SubexpNames())
	for i, r := range findings {
		t.Logf("%d => '%s'\n", i, r)
	}
}
