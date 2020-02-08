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
