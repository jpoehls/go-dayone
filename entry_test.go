package dayone

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewEntryHasUUID(t *testing.T) {
	e := newEntry()

	if e.UUID() == "" {
		t.Error("missing uuid")
	}

	if len(e.UUID()) != 32 {
		t.Error("uuid too short")
	}

	if strings.ToUpper(e.UUID()) != e.UUID() {
		t.Error("uuid should be upper")
	}
}

func TestValidateMissingUUID(t *testing.T) {
	e := &Entry{}

	err := e.validate()
	if err == nil || err.Error() != "missing uuid" {
		t.Fail()
	}
}

func TestValidatePasses(t *testing.T) {
	e := newEntry()

	err := e.validate()
	if err != nil {
		t.Error(err)
	}
}

func TestParsingEntry(t *testing.T) {
	d := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Activity</key>
	<string>Automotive</string>
	<key>Creation Date</key>
	<date>2014-09-24T01:52:11Z</date>
	<key>Creator</key>
	<dict>
		<key>Device Agent</key>
		<string>iPhone/iPhone7,2</string>
		<key>Generation Date</key>
		<date>2014-09-24T01:52:11Z</date>
		<key>Host Name</key>
		<string>Joshua's iPhone</string>
		<key>OS Agent</key>
		<string>iOS/8.0</string>
		<key>Software Agent</key>
		<string>Day One iOS/1.15</string>
	</dict>
	<key>Entry Text</key>
	<string>#title line

body line</string>
	<key>Ignore Step Count</key>
	<true/>
	<key>Location</key>
	<dict>
		<key>Administrative Area</key>
		<string>TX</string>
		<key>Country</key>
		<string>United States</string>
		<key>Latitude</key>
		<real>39.988470999551925</real>
		<key>Locality</key>
		<string>Somewhere</string>
		<key>Longitude</key>
		<real>-87.8764692974257</real>
		<key>Foursquare ID</key>
		<string>372837</string>
		<key>Place Name</key>
		<string>199 Address Ln</string>
		<key>Region</key>
		<dict>
			<key>Center</key>
			<dict>
				<key>Latitude</key>
				<real>39.98898599897538</real>
				<key>Longitude</key>
				<real>-87.876626999999985</real>
			</dict>
			<key>Radius</key>
			<real>70.891240772618431</real>
		</dict>
	</dict>
	<key>Starred</key>
	<true/>
	<key>Step Count</key>
	<integer>1043</integer>
	<key>Tags</key>
	<array>
		<string>bjj</string>
		<string>fitness</string>
	</array>
	<key>Time Zone</key>
	<string>America/Chicago</string>
	<key>UUID</key>
	<string>FF755C6D7D9B4A5FBC4E41C07D622C65</string>
	<key>Publish URL</key>
	<string>http://google.com</string>
	<key>Music</key>
	<dict>
		<key>Album</key>
		<string>Holy Bible - KJV</string>
		<key>Artist</key>
		<string>Alexander Scourby</string>
		<key>Track</key>
		<string>The Book of Genesis - Chapter 49</string>
		<key>Album Year</key>
		<string>2010</string>
	</dict>
	<key>Weather</key>
	<dict>
		<key>Celsius</key>
		<string>24</string>
		<key>Description</key>
		<string>Clear</string>
		<key>Fahrenheit</key>
		<string>75</string>
		<key>IconName</key>
		<string>clearn.png</string>
		<key>Pressure MB</key>
		<integer>1017</integer>
		<key>Relative Humidity</key>
		<integer>47</integer>
		<key>Service</key>
		<string>HAMweather</string>
		<key>Sunrise Date</key>
		<date>2014-09-23T12:20:49Z</date>
		<key>Sunset Date</key>
		<date>2014-09-24T00:26:11Z</date>
		<key>Visibility KM</key>
		<real>16.093440000000001</real>
		<key>Wind Bearing</key>
		<integer>80</integer>
		<key>Wind Chill Celsius</key>
		<integer>24</integer>
		<key>Wind Speed KPH</key>
		<integer>11</integer>
	</dict>
</dict>
</plist>`

	var e Entry
	if err := e.parse(bytes.NewReader([]byte(d))); err != nil {
		t.Fatal(err)
	}

	if e.EntryText != `#title line

body line` {
		t.Error("entry text")
	}

	if e.Activity != "Automotive" {
		t.Error("activity")
	}

	if e.IgnoreStepCount != true {
		t.Error("ignore step count")
	}

	if e.Starred != true {
		t.Error("starred")
	}

	if e.StepCount != 1043 {
		t.Error("step count")
	}

	if e.TimeZone != "America/Chicago" {
		t.Error("time zone")
	}

	if e.UUID() != "FF755C6D7D9B4A5FBC4E41C07D622C65" {
		t.Error("uuid")
	}

	if e.PublishURL != "http://google.com" {
		t.Error("publish url")
	}

	if len(e.Tags) != 2 || e.Tags[0] != "bjj" || e.Tags[1] != "fitness" {
		t.Error("tags")
	}

	if e.Music == nil {
		t.Error("music")
	} else {
		if e.Music.Album != "Holy Bible - KJV" {
			t.Error("music album")
		}
		if e.Music.Artist != "Alexander Scourby" {
			t.Error("music artist")
		}
		if e.Music.Track != "The Book of Genesis - Chapter 49" {
			t.Error("music track")
		}
		if e.Music.AlbumYear != "2010" {
			t.Error("music album year")
		}
	}

	creationDate, _ := time.Parse(time.RFC3339, "2014-09-24T01:52:11Z")
	if e.CreationDate != creationDate {
		t.Errorf("creation date, actual: %s, expected: %s", e.CreationDate, creationDate)
	}

	if e.Creator == nil {
		t.Error("creator")
	} else {
		if e.Creator.DeviceAgent != "iPhone/iPhone7,2" {
			t.Error("creator device agent")
		}

		genDate, _ := time.Parse(time.RFC3339, "2014-09-24T01:52:11Z")
		if e.Creator.GenerationDate != genDate {
			t.Error("creator generation date")
		}

		if e.Creator.HostName != "Joshua's iPhone" {
			t.Error("creator host name")
		}

		if e.Creator.OSAgent != "iOS/8.0" {
			t.Error("creator os agent")
		}

		if e.Creator.SoftwareAgent != "Day One iOS/1.15" {
			t.Error("creator software agent")
		}
	}

	if e.Location == nil {
		t.Error("location")
	} else {
		if e.Location.AdministrativeArea != "TX" {
			t.Error("location administrative area")
		}

		if e.Location.Country != "United States" {
			t.Error("location country")
		}

		if e.Location.Latitude != 39.988470999551925 {
			t.Error("location latitude")
		}

		if e.Location.Longitude != -87.8764692974257 {
			t.Error("location longitude")
		}

		if e.Location.PlaceName != "199 Address Ln" {
			t.Error("location place name")
		}

		if e.Location.Locality != "Somewhere" {
			t.Error("location locality")
		}

		if e.Location.FoursquareID != "372837" {
			t.Error("location foursquare id")
		}

		if e.Location.Region == nil {
			t.Error("location region")
		} else {
			if e.Location.Region.Radius != 70.891240772618431 {
				t.Error("location region radius")
			}

			if e.Location.Region.Center == nil {
				t.Error("location region center")
			} else {
				if e.Location.Region.Center.Latitude != 39.98898599897538 {
					t.Error("location region center latitude")
				}

				if e.Location.Region.Center.Longitude != -87.876626999999985 {
					t.Error("location region center longitude")
				}
			}
		}

		if e.Weather == nil {
			t.Error("weather")
		} else {
			if e.Weather.Celsius != "24" {
				t.Error("weather celsius")
			}
			if e.Weather.Description != "Clear" {
				t.Error("weather description")
			}
			if e.Weather.Fahrenheit != "75" {
				t.Error("weather fahrenheit")
			}
			if e.Weather.IconName != "clearn.png" {
				t.Error("weather icon name")
			}
			if e.Weather.PressureMB != 1017 {
				t.Error("weather pressure mb")
			}
			if e.Weather.RelativeHumidity != 47 {
				t.Error("weather relative humidity")
			}
			if e.Weather.Service != "HAMweather" {
				t.Error("weather service")
			}

			sunrise, _ := time.Parse(time.RFC3339, "2014-09-23T12:20:49Z")
			if e.Weather.SunriseDate != sunrise {
				t.Error("weather sunrise date")
			}

			sunset, _ := time.Parse(time.RFC3339, "2014-09-24T00:26:11Z")
			if e.Weather.SunsetDate != sunset {
				t.Error("weather sunset date")
			}

			if e.Weather.VisibilityKM != 16.093440000000001 {
				t.Error("weather visibility km")
			}
			if e.Weather.WindBearing != 80 {
				t.Error("weather wind bearing")
			}
			if e.Weather.WindChillCelsius != 24 {
				t.Error("weather wind chill celsius")
			}
			if e.Weather.WindSpeedKPH != 11 {
				t.Error("weather wind speed kph")
			}
		}
	}
}
