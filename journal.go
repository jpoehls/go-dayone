package dayone

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const entryExt = ".doentry"

// ErrStopRead is an error you can return from a
// ReadFunc to stop reading journal entries.
var ErrStopRead = errors.New("stop reading")

// Journal is the top-level type for reading Day One journal files.
type Journal struct {
	dir string
}

// NewJournal creates a new Journal for the
// specified dir.
func NewJournal(dir string) *Journal {
	return &Journal{
		dir: dir,
	}
}

func (j *Journal) getEntriesDir() string {
	return filepath.Join(j.dir, "entries")
}

/*
func (j Journal) Write(e *Entry) error {
	if e.id == "" {
		// TODO: Add support for writing new entries.
		//       - Ensure unique on FS (probably by just not allowing overwrite file on write)
		//       - Write to file
		return errors.New("cannot write new entries: not supported yet")
	}

	if err := e.validate(); err != nil {
		return err
	}

	// TODO: Overwrite file in journal

	return nil
}
*/

// ReadEntry reads the entry with the specified id.
func (j *Journal) ReadEntry(id string) (*Entry, error) {
	path := filepath.Join(j.getEntriesDir(), id+entryExt)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	e := &Entry{}
	err = e.parse(f)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// ReadFunc is the func to use when enumerating journal entries.
type ReadFunc func(e *Entry, err error) error

// Read enumerates all of the journal entries and calls
// fn with each entry found. Errors returned by fn
// are returned by Read. fn can return StopError
// to halt enumeration at any point.
func (j *Journal) Read(fn ReadFunc) error {

	var err error
	var e *Entry

	files, err := ioutil.ReadDir(j.getEntriesDir())
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !isEntryFile(f.Name()) {
			continue
		}

		uuid := strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name()))
		e, err = j.ReadEntry(uuid)
		err = fn(e, err)

		if err == ErrStopRead {
			return nil
		} else if err != nil {
			return errors.New(err.Error() + ": file: " + f.Name())
		}
	}

	return nil
}

func isEntryFile(name string) bool {
	if strings.EqualFold(filepath.Ext(name), entryExt) {
		return true
	}

	return false
}
