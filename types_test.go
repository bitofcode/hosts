package hosts

import "testing"

func TestSimpleEntry_AddHostName(t *testing.T) {
	name := "localhost"
	ent := NewEntryIp("127.0.0.1")
	ent.AddHostName(name)
	ent.AddHostName(name)

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
