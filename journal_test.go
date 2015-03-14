package dayone

import (
	"testing"
)

// no error reading empty journal dir
// no error reading journal with empty entries dir
// reads expected entries from default journal
//

func noopRead(e *Entry, err error) error {
	return err
}

func TestReadingEmptyJournalDir(t *testing.T) {
	j := NewJournal("./test_journals/empty_no_dirs")

	err := j.Read(noopRead)

	if err == nil {
		t.Fatal("expected path not found error")
	}

	//err.(os.PathError)
}

func TestReadingEmptyEntriesDir(t *testing.T) {
	j := NewJournal("./test_journals/empty_with_dirs")

	err := j.Read(noopRead)

	if err != nil {
		t.Fatal(err)
	}
}
