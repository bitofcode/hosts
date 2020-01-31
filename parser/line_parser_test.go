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
			_, err := ParseFromLine(test.line)
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
			_, err := ParseFromLine(test.line)
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
			ent, err := ParseFromLine(test.line)
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

			gotLine, err := ParseToLine(hosts.NewEntry(tt.args.ip, tt.args.hostNames))

			if tt.wantErr && (tt.err != err) {
				t.Errorf("ParseToLine(%v) expected error = '%v', actual = '%v'", tt.args, tt.err, err)
				return
			}

			if !tt.wantErr && (err != nil) {
				t.Errorf("ParseToLine(%v) unexpected error = %v", tt.args, err)
				return
			}

			if gotLine != tt.wantLine {
				t.Errorf("ParseToLine(%v) want %v, gotLine = %v", tt.args, tt.wantLine, gotLine)
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

			gotLine, err := ParseToLine(hosts.NewEntry(tt.args.ip, tt.args.hostNames))

			if tt.wantErr && (tt.err != err) {
				t.Errorf("ParseToLine(%v) expected error = '%v', actual = '%v'", tt.args, tt.err, err)
				return
			}

			if !tt.wantErr && (err != nil) {
				t.Errorf("ParseToLine(%v) unexpected error = %v", tt.args, err)
				return
			}

			if gotLine != tt.wantLine {
				t.Errorf("ParseToLine(%v) want '%v', gotLine = '%v'", tt.args, tt.wantLine, gotLine)
			}
		})
	}
}
