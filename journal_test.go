package dayone

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestOpenMissingPhoto(t *testing.T) {
	j := NewJournal("./test_journals/default")

	r, err := j.OpenPhoto("bad uuid")
	if !os.IsNotExist(err) {
		t.Log(err)
		t.Error("expected an os not exist error")
	}

	if r != nil {
		r.Close()
		t.Error("expected nil reader")
	}
}

func TestOpenPhoto(t *testing.T) {
	j := NewJournal("./test_journals/default")

	r, err := j.OpenPhoto("871D0F435D7B469C9429CD441A9E74B5")
	if err != nil {
		t.Error(err)
	}

	if r == nil {
		t.Error("expected reader not nil")
	} else {
		r.Close()
	}
}

func TestStatPhoto(t *testing.T) {
	j := NewJournal("./test_journals/default")

	i, err := j.PhotoStat("871D0F435D7B469C9429CD441A9E74B5")
	if err != nil {
		t.Error(err)
	}

	if i == nil {
		t.Error("file stat")
	}
}

func TestStatEntry(t *testing.T) {
	j := NewJournal("./test_journals/default")

	i, err := j.EntryStat("871D0F435D7B469C9429CD441A9E74B5")
	if err != nil {
		t.Error(err)
	}

	if i == nil {
		t.Error("file stat")
	}
}

func TestStatMissingPhoto(t *testing.T) {
	j := NewJournal("./test_journals/default")

	i, err := j.PhotoStat("missing uuid")
	if !os.IsNotExist(err) {
		t.Log(err)
		t.Error("expected an os not exist error")
	}

	if i != nil {
		t.Error("expected nil file info")
	}
}

func TestStatMissingEntry(t *testing.T) {
	j := NewJournal("./test_journals/default")

	i, err := j.EntryStat("missing uuid")
	if !os.IsNotExist(err) {
		t.Log(err)
		t.Error("expected an os not exist error")
	}

	if i != nil {
		t.Error("expected nil file info")
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

	myerr := errors.New("boom")
	err := j.Read(func(e *Entry, err error) error {
		count++
		return myerr
	})

	if count != 1 {
		t.Error("read func called too many times")
	}

	t.Logf("error: %v", err)
	if !strings.Contains(err.Error(), "boom") {
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

func TestWriteEntryWithNoId(t *testing.T) { // dup with EntryNoId test?
	j := NewJournal("./test_journals/default")

	e := &Entry{}

	err := j.Write(e)

	if err == nil || err.Error() != "missing uuid" {
		t.Fail()
	}
}

func TestWriteOverwriteExisting(t *testing.T) {
	j := NewJournal("./test_journals/default")

	e := new(Entry)
	e.uuid = "871D0F435D7B469C9429CD441A9E74B5"

	err := j.Write(e)

	if err == nil || err.Error() != "overwriting existing entry is not supported yet" {
		t.Fail()
	}
}

func TestWriteEntry(t *testing.T) {
	j := NewJournal("./test_journals/write")

	e := NewEntry()
	e.EntryText = "hello\nworld"
	e.Tags = []string{"hello", "world"}
	e.CreationDate = time.Now()

	err := j.Write(e)

	if err != nil {
		t.Fatal(err)
	}

	// do we need to verify we wrote the file?
}
