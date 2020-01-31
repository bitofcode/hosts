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
