package models

import "time"

type MapData struct {
	ID              string `json:"id"`
	SummaryPolyline string `json:"summary_polyline"`
	ResourceState   int    `json:"resource_state"`
}

type AthleteData struct {
	ID            int `json:"id"`
	ResourceState int `json:"resource_state"`
}

type Activity struct {
	ResourceState              int         `json:"resource_state"`
	Athlete                    AthleteData `json:"athlete"`
	Name                       string      `json:"name"`
	Distance                   float64     `json:"distance"`
	MovingTime                 int         `json:"moving_time"`
	ElapsedTime                int         `json:"elapsed_time"`
	TotalElevationGain         float64     `json:"total_elevation_gain"`
	Type                       string      `json:"type"`
	SportType                  string      `json:"sport_type"`
	WorkoutType                interface{} `json:"workout_type"`
	ID                         int         `json:"id"`
	StartDate                  time.Time   `json:"start_date"`
	StartDateLocal             time.Time   `json:"start_date_local"`
	Timezone                   string      `json:"timezone"`
	UTCOffset                  float64     `json:"utc_offset"`
	LocationCity               interface{} `json:"location_city"`
	LocationState              interface{} `json:"location_state"`
	LocationCountry            string      `json:"location_country"`
	AchievementCount           int         `json:"achievement_count"`
	KudosCount                 int         `json:"kudos_count"`
	CommentCount               int         `json:"comment_count"`
	AthleteCount               int         `json:"athlete_count"`
	PhotoCount                 int         `json:"photo_count"`
	Map                        MapData     `json:"map"`
	Trainer                    bool        `json:"trainer"`
	Commute                    bool        `json:"commute"`
	Manual                     bool        `json:"manual"`
	Private                    bool        `json:"private"`
	Visibility                 string      `json:"visibility"`
	Flagged                    bool        `json:"flagged"`
	GearID                     interface{} `json:"gear_id"`
	StartLatLng                []float64   `json:"start_latlng"`
	EndLatLng                  []float64   `json:"end_latlng"`
	AverageSpeed               float64     `json:"average_speed"`
	MaxSpeed                   float64     `json:"max_speed"`
	HasHeartrate               bool        `json:"has_heartrate"`
	AverageHeartrate           float64     `json:"average_heartrate"`
	MaxHeartrate               float64     `json:"max_heartrate"`
	HeartrateOptOut            bool        `json:"heartrate_opt_out"`
	DisplayHideHeartrateOption bool        `json:"display_hide_heartrate_option"`
	UploadID                   int64       `json:"upload_id"`
	UploadIDStr                string      `json:"upload_id_str"`
	ExternalID                 string      `json:"external_id"`
	FromAcceptedTag            bool        `json:"from_accepted_tag"`
	PRCount                    int         `json:"pr_count"`
	TotalPhotoCount            int         `json:"total_photo_count"`
	HasKudoed                  bool        `json:"has_kudoed"`
}
