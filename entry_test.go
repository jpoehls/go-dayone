package dayone

import (
	"bytes"
	"testing"
)

// no error reading empty journal dir
// no error reading journal with empty entries dir
// reads expected entries from default journal
//

func TestUnmarshalXML(t *testing.T) {
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
		<string>Joshua Poehls’s iPhone</string>
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
		<real>29.988470999551925</real>
		<key>Locality</key>
		<string>Kyle</string>
		<key>Longitude</key>
		<real>-97.8764692974257</real>
		<key>Place Name</key>
		<string>206 W Center St</string>
		<key>Region</key>
		<dict>
			<key>Center</key>
			<dict>
				<key>Latitude</key>
				<real>29.98898599897538</real>
				<key>Longitude</key>
				<real>-97.876626999999985</real>
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
	</array>
	<key>Time Zone</key>
	<string>America/Chicago</string>
	<key>UUID</key>
	<string>FF755C6D7D9B4A5FBC4E41C07D622C65</string>
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

	if len(e.Tags) != 1 || e.Tags[0] != "bjj" {
		t.Error("tags")
	}

	if e.Creator == nil || e.Creator.DeviceAgent != "iPhone/iPhone7,2" {
		t.Error("creator device agent")
	}
}
