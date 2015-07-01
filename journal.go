package dayone

import (
	"errors"
	"fmt"
	"github.com/juju/errgo"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const entryExt = ".doentry"
const photoExt = ".jpg"

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

func (j *Journal) getPhotosDir() string {
	return filepath.Join(j.dir, "photos")
}

// Write creates a new entry
func (j Journal) Write(e *Entry) error {
	//TODO: add overwrite existing journal entry
	//TODO: add photo support

	if err := e.validate(); err != nil {
		return err
	}

	// stat file on, for now err if file exists
	path := filepath.Join(j.getEntriesDir(), e.UUID()+entryExt)
	f, err := os.Stat(path)

	if !os.IsNotExist(err) {
		return errors.New("overwriting existing entry is not supported yet")
	} else {
		return errors.New("something else f: " + fmt.Sprintf("%v", f))
	}

	// write new entry created

	return nil
}

// PhotoStat returns the result of os.Stat() for the
// photo associated with the entry uuid.
func (j *Journal) PhotoStat(uuid string) (os.FileInfo, error) {
	path := filepath.Join(j.getPhotosDir(), uuid+photoExt)

	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		} else {
			return nil, errgo.Mask(err)
		}
	}

	return f, nil
}

// OpenPhoto opens an io.ReadCloser for the photo file
// associated with the specified entry uuid or returns an error.
func (j *Journal) OpenPhoto(uuid string) (io.ReadCloser, error) {
	path := filepath.Join(j.getPhotosDir(), uuid+photoExt)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		} else {
			return nil, errgo.Mask(err)
		}
	}

	return f, nil
}

// EntryStat returns the result of os.Stat() for the
// entry with the specified uuid.
func (j *Journal) EntryStat(uuid string) (os.FileInfo, error) {
	path := filepath.Join(j.getEntriesDir(), uuid+entryExt)

	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		} else {
			return nil, errgo.Mask(err)
		}
	}

	return f, nil
}

// ReadEntry reads the entry with the specified id.
func (j *Journal) ReadEntry(uuid string) (*Entry, error) {
	path := filepath.Join(j.getEntriesDir(), uuid+entryExt)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		} else {
			return nil, errgo.Mask(err)
		}
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
		if os.IsNotExist(err) {
			return err
		} else {
			return errgo.Mask(err)
		}
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
			return errgo.NoteMask(err, "file: "+f.Name())
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
