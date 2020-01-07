package astrotime

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type location struct {
	timezone  *time.Location
	latitude  float64
	longitude float64
}

type times struct {
	midnight time.Time
	sunrise  time.Time
	sunset   time.Time
	dawn     time.Time
	dusk     time.Time
}

func mustNewTimes(midnight time.Time, sunrise string, sunset string, dawn string, dusk string) times {
	f := func(t string) time.Time {
		ts, err := time.Parse("15:04", t)
		if err != nil {
			panic(err)
		}

		hh, mm, ss := ts.Clock()
		return midnight.Add(time.Duration(hh) * time.Hour).Add(time.Duration(mm) * time.Minute).Add(time.Duration(ss) * time.Second)
	}

	return times{
		midnight: midnight,
		sunrise:  f(sunrise),
		sunset:   f(sunset),
		dawn:     f(dawn),
		dusk:     f(dusk),
	}
}

type test struct {
	times    times
	location location
	error    time.Duration
}

// test data from http://www.timeanddate.com/worldclock/sunrise.html
var (
	sydney = location{
		timezone:  mustParseTimezone("Australia/Sydney"),
		latitude:  -33.900626,
		longitude: 151.196210,
	}
	stockholm = location{
		timezone:  mustParseTimezone("Europe/Stockholm"),
		latitude:  59.332953,
		longitude: 18.059986,
	}
	newyork = location{
		timezone:  mustParseTimezone("America/New_York"),
		latitude:  40.692197,
		longitude: -73.940547,
	}

	minute = time.Minute
	DATA   = []*test{
		&test{mustNewTimes(november(sydney.timezone), "05:55", "19:23", "05:29", "19:49"), sydney, minute},
		&test{mustNewTimes(november(stockholm.timezone), "07:07", "15:55", "06:23", "16:39"), stockholm, minute},
		&test{mustNewTimes(november(newyork.timezone), "07:26", "17:52", "06:58", "18:21"), newyork, minute},
		&test{mustNewTimes(july(sydney.timezone), "07:01", "16:57", "06:33", "17:25"), sydney, minute},
		&test{mustNewTimes(july(stockholm.timezone), "03:37", "22:06", "02:09", "23:33"), stockholm, minute},
		&test{mustNewTimes(july(newyork.timezone), "05:29", "20:31", "04:55", "21:04"), newyork, minute},
	}
)

func july(location *time.Location) time.Time {
	return time.Date(2014, 7, 1, 0, 0, 0, 0, location)
}

func november(location *time.Location) time.Time {
	return time.Date(2014, 11, 1, 0, 0, 0, 0, location)
}

func mustParseTimezone(tz string) *time.Location {
	tl, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	return tl
}

func TestForFixedLocations(t *testing.T) {
	for _, datum := range DATA {
		t.Run(fmt.Sprintf("%s-%s", datum.location.timezone, datum.times.midnight), func(t *testing.T) {
			checkDate := func(t *testing.T, description string, expected time.Time, calculated time.Time) {
				t.Run(description, func(t *testing.T) {
					require.GreaterOrEqual(t, datum.error.Seconds(), math.Abs(expected.Sub(calculated).Seconds()))
				})
			}

			checkDate(t, "dawn", datum.times.dawn, NextDawn(datum.times.midnight, datum.location.latitude, datum.location.longitude, CIVIL_DAWN))
			checkDate(t, "sunrise", datum.times.sunrise, NextSunrise(datum.times.midnight, datum.location.latitude, datum.location.longitude))
			checkDate(t, "sunset", datum.times.sunset, NextSunset(datum.times.midnight, datum.location.latitude, datum.location.longitude))
			checkDate(t, "dusk", datum.times.dusk, NextDusk(datum.times.midnight, datum.location.latitude, datum.location.longitude, CIVIL_DUSK))
		})
	}
}
