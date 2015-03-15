package dayone

import (
	"errors"
	"strings"
	"testing"
)

func TestReadingMyOwnJournal(t *testing.T) {
	j := NewJournal("/Users/joshua/Dropbox/Apps/Day One/Journal.dayone")

	count := 0
	err := j.Read(func(e *Entry, err error) error {
		if err != nil {
			return err
		}
		count++
		return nil
	})

	t.Logf("%v entries read", count)

	if err != nil {
		t.Fatal(err)
	}
}

func TestReadStopsWhenAsked(t *testing.T) {
	j := NewJournal("./test_journals/default")

	count := 0
	err := j.Read(func(e *Entry, err error) error {
		count++
		return ErrStopRead
	})

	if err != nil {
		t.Log(err)
		t.Error("expected no error")
	}

	if count != 1 {
		t.Error("read too many entries")
	}
}

func TestReadingEntry(t *testing.T) {
	j := NewJournal("./test_journals/default")

	e, err := j.ReadEntry("FF755C6D7D9B4A5FBC4E41C07D622C65")

	if err != nil {
		t.Error(err)
	} else {
		if e.UUID() != "FF755C6D7D9B4A5FBC4E41C07D622C65" {
			t.Error("entry doesn't have right uuid")
		}
	}
}

func TestReadingMissingEntry(t *testing.T) {
	j := NewJournal("./test_journals/default")

	e, err := j.ReadEntry("bad uuid")

	if err == nil {
		t.Error("expected an error")
	}

	if e != nil {
		t.Error("expected nil entry")
	}
}

func TestReadBubblesError(t *testing.T) {
	j := NewJournal("./test_journals/default")

	count := 0

	err := j.Read(func(e *Entry, err error) error {
		count++
		return errors.New("boom")
	})

	if count != 1 {
		t.Error("read func called too many times")
	}

	if !strings.HasPrefix(err.Error(), "boom: file:") {
		t.Logf("error: %v", err)
		t.Error("didn't bubble error")
	}
}

func TestReadingAllEntries(t *testing.T) {
	j := NewJournal("./test_journals/default")

	count := 0
	expected := []string{
		"871D0F435D7B469C9429CD441A9E74B5",
		"FF755C6D7D9B4A5FBC4E41C07D622C65",
	}

	err := j.Read(func(e *Entry, err error) error {
		if err != nil {
			return err
		}

		if e == nil {
			t.Fatal("entry was nil")
		}

		if count > len(expected)-1 {
			t.Fatal("more entries than expectd")
		}

		if e.UUID() != expected[count] {
			t.Error("unexpected entry uuid")
		}

		count++
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	if count != len(expected) {
		t.Error("didn't find all entries")
	}
}

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
