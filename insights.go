package sleepiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// SLEEPER ACTIVITIES
// ============================================================================

// SleeperActivities describes sleeper activites that are sourced
// from external health monitors such as Apple Watch and Nest
// thermostats.
type SleeperActivities struct {
	Activities []struct {
		SleeperID    string `json:"sleeperId"`
		ActivityDate string `json:"activityDate"`
		Partner      struct {
			Nest struct {
				SummaryData string      `json:"summary_data"`
				Status      interface{} `json:"status"`
			} `json:"nest"`
			Apple struct {
				SummaryData interface{} `json:"summary_data"`
				GoalSteps   interface{} `json:"goal_steps"`
				DailySteps  interface{} `json:"daily_steps"`
				Status      interface{} `json:"status"`
			} `json:"apple"`
		} `json:"partner"`
	} `json:"activities"`
	Statuses struct {
		Fitbit      bool `json:"fitbit"`
		Underarmour bool `json:"underarmour"`
		Nest        bool `json:"nest"`
		Withings    bool `json:"withings"`
		Health      bool `json:"health"`
		Apple       bool `json:"apple"`
		Honeywell   bool `json:"honeywell"`
		Google      bool `json:"google"`
	} `json:"statuses"`
	Error ServiceError `json:"Error"`
}

// InsightsActiviy obtains activites that are sourced from external
// monitors such as Apple Watch and Nest thermostats
func (s SleepIQ) InsightsActiviy(sleeperID string, startDate time.Time, endDate time.Time) (SleeperActivities, error) {
	var response SleeperActivities

	// Bail if there is not an active logged-in session
	if !s.isInsightsLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://sleepiqapi.azure-api.net/prod/activities?sleeperId={{sleeperId}}&startDate={{startDate}}&endDate={{endDate}}&access_token={{token}}", "{{token}}", s.insightsToken, -1)
	url = strings.Replace(url, "{{startDate}}", startDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{endDate}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getInsightsHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve Insights activities - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read Insights activities - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// PROVIDERS
// ============================================================================

// InsightProvidersStatus describes the status of the various activity
// monitors that are supported by the Insights service.
type InsightProvidersStatus struct {
	Providers []struct {
		ID          string        `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Image       string        `json:"image"`
		Scope       []string      `json:"scope"`
		Platforms   []string      `json:"platforms"`
		DataTypes   []string      `json:"data_types"`
		Connected   bool          `json:"connected"`
		Order       int           `json:"order"`
		ConnectedAt time.Time     `json:"connectedAt"`
		LastSync    interface{}   `json:"last_sync"`
		IsValid     interface{}   `json:"is_valid"`
		Permissions []interface{} `json:"permissions"`
	} `json:"providers"`
	Error ServiceError `json:"Error"`
}

// InsightsProviders retrieves information about the status of the
// various activity monitors that are supported by SleepIQ
func (s SleepIQ) InsightsProviders() (InsightProvidersStatus, error) {
	var response InsightProvidersStatus

	// Bail if there is not an active logged-in session
	if !s.isInsightsLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://sleepiqapi.azure-api.net/prod/providers/?access_token={{token}}", "{{token}}", s.insightsToken, -1)

	responseBytes, err := httpGet(url, s.cookies, getInsightsHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve Insights providers - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read Insights providers - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER RELATIVE INSIGHTS
// ============================================================================

// RelativeInsights provides a historical view of data relative to ones self
type RelativeInsights struct {
	Data []struct {
		Count       int    `json:"count"`
		Date        string `json:"date"`
		SiqScore    int    `json:"siqScore"`
		SleepNumber int    `json:"sleepNumber"`
		TimeInBed   int    `json:"timeInBed"`
	} `json:"data"`
	Error ServiceError `json:"Error"`
}

// InsightsLikeMe retrieves historical data for people with similar sleep
// patterns to yourself
func (s SleepIQ) InsightsLikeMe(sleeperID string, startDate time.Time, endDate time.Time) (RelativeInsights, error) {
	var response RelativeInsights

	// Bail if there is not an active logged-in session
	if !s.isInsightsLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://sleepiqapi.azure-api.net/prod/insights/historical/likeme/{{sleeperId}}?start={{startDate}}&end={{endDate}}&access_token={{token}}", "{{token}}", s.insightsToken, -1)
	url = strings.Replace(url, "{{startDate}}", startDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{endDate}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getInsightsHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve Insights like me - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read Insights like me - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// InsightsNearMe retrieves historical data for people near my location
func (s SleepIQ) InsightsNearMe(sleeperID string, startDate time.Time, endDate time.Time) (RelativeInsights, error) {
	var response RelativeInsights

	// Bail if there is not an active logged-in session
	if !s.isInsightsLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://sleepiqapi.azure-api.net/prod/insights/historical/nearme/{{sleeperId}}?start={{startDate}}&end={{endDate}}&access_token={{token}}", "{{token}}", s.insightsToken, -1)
	url = strings.Replace(url, "{{startDate}}", startDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{endDate}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getInsightsHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve Insights like me - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read Insights like me - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER PERSONAL INSIGHTS
// ============================================================================

// MyInsights contains monthly insight data about ones self
type MyInsights struct {
	Data []struct {
		Count            int    `json:"count"`
		Date             string `json:"date"`
		MaxScore         int    `json:"maxScore"`
		MaxScoreDate     string `json:"maxScoreDate"`
		MaxTimeInBed     int    `json:"maxTimeInBed"`
		MaxTimeInBedDate string `json:"maxTimeInBedDate"`
		SiqScore         int    `json:"siqScore"`
		SleepNumber      int    `json:"sleepNumber"`
		TimeInBed        int    `json:"timeInBed"`
		TotalTimeInBed   int    `json:"totalTimeInBed"`
	} `json:"data"`
	Error ServiceError `json:"Error"`
}

// InsightsMe retrieves historical insight data about ones self
func (s SleepIQ) InsightsMe(sleeperID string, startDate time.Time, endDate time.Time) (MyInsights, error) {
	var response MyInsights

	// Bail if there is not an active logged-in session
	if !s.isInsightsLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://sleepiqapi.azure-api.net/prod/insights/historical/sleeper/{{sleeperId}}?start={{startDate}}&end={{endDate}}&access_token={{token}}", "{{token}}", s.insightsToken, -1)
	url = strings.Replace(url, "{{startDate}}", startDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{endDate}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getInsightsHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve Insights like me - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read Insights like me - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SUPPORTING FUNCTIONS
// ============================================================================

// getInsightsHeaders retrieves the set of headers required to make API
// calls against the Insights service
func getInsightsHeaders() map[string]string {
	var headers map[string]string
	headers = make(map[string]string)

	headers["Accept"] = "application/json, text/javascript, */*; q=0.01"
	headers["Content-Type"] = "application/json"
	headers["Ocp-Apim-Subscription-Key"] = "3c924e14923642baa1c4ad1d5096a1c5"

	return headers
}
