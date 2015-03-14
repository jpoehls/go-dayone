package dayone

import (
	"errors"
	"github.com/DHowett/go-plist"
	"github.com/twinj/uuid"
	"io"
	"time"
)

type Entry struct {
	uuid      string
	EntryText string

	Activity        string
	IgnoreStepCount bool
	StepCount       uint64

	Starred bool

	Tags    []string
	Weather *Weather

	TimeZone     string
	Creator      *Creator
	CreationDate time.Time
}

type Creator struct {
	DeviceAgent    string
	GenerationDate time.Time
	HostName       string
	OSAgent        string
	SoftwareAgent  string
}

type Location struct {
	AdministrativeArea string
	Country            string
	Locality           string
	PlaceName          string
	Region             *Region

	Coordinate
}

type Region struct {
	Center *Coordinate
	Radius float64
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Weather struct {
	Celsius          string
	Fahrenheit       string
	Description      string
	IconName         string
	PressureMB       int
	RelativeHumidity int
	Service          string
	SunriseDate      time.Time
	SunsetDate       time.Time
	VisibilityKM     float64
	WindBearing      int
	WindChillCelsius int
	WindSpeedKPH     int
}

func (e *Entry) UUID() string {
	return e.uuid
}

func newEntry() *Entry {
	id := uuid.NewV4()

	return &Entry{
		uuid: uuid.Formatter(id, uuid.Clean), // e.g. FF755C6D7D9B4A5FBC4E41C07D622C65
	}
}

func (e *Entry) validate() error {
	if e.uuid == "" {
		return errors.New("missing uuid")
	}

	return nil
}

func (e *Entry) parse(r io.ReadSeeker) error {
	dec := plist.NewDecoder(r)

	var dict map[string]interface{}
	if err := dec.Decode(&dict); err != nil {
		return err
	}
	for k, v := range dict {
		switch k {
		case "UUID":
			e.uuid = v.(string)
		case "Entry Text":
			e.EntryText = v.(string)
		case "Activity":
			e.Activity = v.(string)
		case "Time Zone":
			e.TimeZone = v.(string)
		case "Ignore Step Count":
			e.IgnoreStepCount = v.(bool)
		case "Starred":
			e.Starred = v.(bool)
		case "Step Count":
			e.StepCount = v.(uint64)
		case "Tags":
			v2 := v.([]interface{})
			e.Tags = e.Tags[:0]
			for _, av := range v2 {
				e.Tags = append(e.Tags, av.(string))
			}
		case "Creator":
			if e.Creator == nil {
				e.Creator = &Creator{}
			}
			err := e.Creator.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Creator) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Device Agent":
			c.DeviceAgent = v.(string)
		}
	}

	return nil
}
