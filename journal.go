package dayone

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const entryExt = ".doentry"

var StopRead = errors.New("stop reading")

type Journal struct {
	dir string
}

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

// TODO: add ability to attach a JPEG photo
// func (j Journal) WritePhoto(e Entry, path string)

type ReadFunc func(e *Entry, err error) error

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

		e, err = j.ReadEntry(filepath.Base(f.Name()))
		err = fn(e, err)

		if err == StopRead {
			return nil
		} else if err != nil {
			return err
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
