package sleepiq

import "net/http"

// Version: 1.0.0

// SleepIQ is the main struct which all methods are associated with
// as well as contains global settings for use by all methods
type SleepIQ struct {
	isLoggedIn         bool
	isInsightsLoggedIn bool
	loginKey           string
	insightsToken      string
	cookies            []*http.Cookie
}

// ServiceError contains error information for calls to the sleepiq service
type ServiceError struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

// New creates a new instance of SleepIQ
func New() SleepIQ {
	s := SleepIQ{}
	return s
}
