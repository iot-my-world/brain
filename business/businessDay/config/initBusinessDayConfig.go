package config

import (
	shift "bitbucket.org/gotimekeeper/business/shift/config"
	"time"
)

var sixAM, _ = time.Parse(time.RFC3339, "2010-01-01T06:00:00+02:00")
var twoPM, _ = time.Parse(time.RFC3339, "2010-01-01T14:00:00+02:00")
var elevenPM, _ = time.Parse(time.RFC3339, "2010-01-01T22:00:00+02:00")

var initConfig = Config{
	Id: "",
	Monday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
		{
			StartDateTime: elevenPM.Unix(),
			EndDateTime: sixAM.Unix(),
		},
	},
	Tuesday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
		{
			StartDateTime: elevenPM.Unix(),
			EndDateTime: sixAM.Unix(),
		},
	},
	Wednesday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
		{
			StartDateTime: elevenPM.Unix(),
			EndDateTime: sixAM.Unix(),
		},
	},
	Thursday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
		{
			StartDateTime: elevenPM.Unix(),
			EndDateTime: sixAM.Unix(),
		},
	},
	Friday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
		{
			StartDateTime: elevenPM.Unix(),
			EndDateTime: sixAM.Unix(),
		},
	},
	Saturday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
	},
	Sunday: []shift.Config{
		{
			StartDateTime: sixAM.Unix(),
			EndDateTime: twoPM.Unix(),
		},
		{
			StartDateTime: twoPM.Unix(),
			EndDateTime: elevenPM.Unix(),
		},
	},
}
