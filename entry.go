package dayone

import (
	"errors"
	"github.com/DHowett/go-plist"
	"github.com/twinj/uuid"
	"io"
	"strings"
	"time"
)

// Entry is the top-level journal entry type.
type Entry struct {
	uuid      string
	EntryText string `plist:"Entry Text"`

	Activity        string
	IgnoreStepCount bool   `plist:"Ignore Step Count"`
	StepCount       uint64 `plist:"Step Count"`

	Starred    bool
	PublishURL string `plist:"Publish URL"`
	Music      *Music

	Tags     []string
	Weather  *Weather
	Location *Location

	TimeZone     string `plist:"Time Zone"`
	Creator      *Creator
	CreationDate time.Time `plist:"Creation Date"`
}

// Creator is the creator of a journal entry.
type Creator struct {
	DeviceAgent    string    `plist:"Device Agent"`
	GenerationDate time.Time `plist:"Generation Date"`
	HostName       string    `plist:Host Name"`
	OSAgent        string    `plist:OS Agent"`
	SoftwareAgent  string    `plist: Software Agent"`
}

// Location of a journal entry.
type Location struct {
	AdministrativeArea string `plist:"Adminstrative Area"`
	Country            string
	Locality           string
	PlaceName          string `plist:"Place Name"`
	Region             *Region
	FoursquareID       string

	Coordinate
}

// Region location data.
type Region struct {
	Center *Coordinate
	Radius float64
}

// Coordinate for location data.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Weather data for a journal entry.
type Weather struct {
	Celsius          string
	Fahrenheit       string
	Description      string
	IconName         string
	PressureMB       float64 `plist:"Pressure MB"`
	RelativeHumidity float64 `plist:"Relative Humidity"`
	Service          string
	SunriseDate      time.Time `plist:"Sunrise Date"`
	SunsetDate       time.Time `plist:"Sunset Date"`
	VisibilityKM     float64   `plist:"Visibility KM"`
	WindBearing      uint64    `plist:"Wind Bearing"`
	WindChillCelsius int64     `plist:"Wind Chill Celsius"`
	WindSpeedKPH     float64   `plist:"Wind Speed KPH"`
}

// Music data for a journal entry.
type Music struct {
	Album     string
	Artist    string
	Track     string
	AlbumYear string `plist:"Album Year"`
}

// UUID gets the unique ID of the entry.
func (e *Entry) UUID() string {
	return e.uuid
}

func NewEntry() *Entry {
	id := uuid.NewV4()

	return &Entry{
		uuid: strings.ToUpper(uuid.Formatter(id, uuid.Clean)), // e.g. FF755C6D7D9B4A5FBC4E41C07D622C65
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
		case "Creation Date":
			e.CreationDate = v.(time.Time)
		case "Tags":
			o, err := parseStringArray(v.([]interface{}))
			if err != nil {
				return err
			}
			e.Tags = o
		case "Creator":
			e.Creator = &Creator{}
			err := e.Creator.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		case "Location":
			e.Location = &Location{}
			err := e.Location.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		case "Weather":
			e.Weather = &Weather{}
			err := e.Weather.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		case "Publish URL":
			e.PublishURL = v.(string)
		case "Music":
			e.Music = &Music{}
			err := e.Music.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			return errors.New("unexpected key: " + k)
		}
	}

	return nil
}

func parseStringArray(in []interface{}) ([]string, error) {
	var out []string
	for _, v := range in {
		out = append(out, v.(string))
	}

	return out, nil
}

func (c *Creator) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Device Agent":
			c.DeviceAgent = v.(string)
		case "Generation Date":
			c.GenerationDate = v.(time.Time)
		case "Host Name":
			c.HostName = v.(string)
		case "OS Agent":
			c.OSAgent = v.(string)
		case "Software Agent":
			c.SoftwareAgent = v.(string)
		default:
			return errors.New("unexpected key: " + k)
		}
	}

	return nil
}

func (l *Location) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Administrative Area":
			l.AdministrativeArea = v.(string)
		case "Country":
			l.Country = v.(string)
		case "Locality":
			l.Locality = v.(string)
		case "Place Name":
			l.PlaceName = v.(string)
		case "Latitude":
			l.Latitude = v.(float64)
		case "Longitude":
			l.Longitude = v.(float64)
		case "Foursquare ID":
			l.FoursquareID = v.(string)
		case "Region":
			l.Region = &Region{}
			err := l.Region.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			return errors.New("unexpected key: " + k)
		}
	}
	return nil
}

func (r *Region) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Radius":
			r.Radius = v.(float64)
		case "Center":
			r.Center = &Coordinate{}
			err := r.Center.parse(v.(map[string]interface{}))
			if err != nil {
				return err
			}
		default:
			return errors.New("unexpected key: " + k)
		}
	}
	return nil
}

func (c *Coordinate) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Latitude":
			c.Latitude = v.(float64)
		case "Longitude":
			c.Longitude = v.(float64)
		default:
			return errors.New("unexpected key: " + k)
		}
	}
	return nil
}

func (w *Weather) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Celsius":
			w.Celsius = v.(string)
		case "Description":
			w.Description = v.(string)
		case "Fahrenheit":
			w.Fahrenheit = v.(string)
		case "IconName":
			w.IconName = v.(string)
		case "Pressure MB":
			switch v.(type) {
			case uint64:
				w.PressureMB = float64(v.(uint64))
			case float64:
				w.PressureMB = v.(float64)
			}
		case "Relative Humidity":
			switch v.(type) {
			case float64:
				w.RelativeHumidity = v.(float64)
			case uint64:
				w.RelativeHumidity = float64(v.(uint64))
			}
		case "Service":
			w.Service = v.(string)
		case "Sunrise Date":
			w.SunriseDate = v.(time.Time)
		case "Sunset Date":
			w.SunsetDate = v.(time.Time)
		case "Visibility KM":
			switch v.(type) {
			case float64:
				w.VisibilityKM = v.(float64)
			case uint64:
				w.VisibilityKM = float64(v.(uint64))
			}
		case "Wind Bearing":
			w.WindBearing = v.(uint64)
		case "Wind Chill Celsius":
			switch v.(type) {
			case int64:
				w.WindChillCelsius = v.(int64)
			case uint64:
				w.WindChillCelsius = int64(v.(uint64))
			}
		case "Wind Speed KPH":
			switch v.(type) {
			case uint64:
				w.WindSpeedKPH = float64(v.(uint64))
			case float64:
				w.WindSpeedKPH = v.(float64)
			}
		default:
			return errors.New("unexpected key: " + k)
		}
	}
	return nil
}

func (m *Music) parse(dict map[string]interface{}) error {
	for k, v := range dict {
		switch k {
		case "Album":
			m.Album = v.(string)
		case "Artist":
			m.Artist = v.(string)
		case "Track":
			m.Track = v.(string)
		case "Album Year":
			m.AlbumYear = v.(string)
		default:
			return errors.New("unexpected key: " + k)
		}
	}
	return nil
}
