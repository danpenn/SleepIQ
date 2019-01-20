package sleepiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// SLEEPER DETAILS
// ============================================================================

// SleeperDetails contains personal information about a sleeper (person)
type SleeperDetails struct {
	Sleepers []struct {
		FirstName      string      `json:"firstName"`
		Active         bool        `json:"active"`
		EmailValidated bool        `json:"emailValidated"`
		IsChild        bool        `json:"isChild"`
		BedID          string      `json:"bedId"`
		BirthYear      string      `json:"birthYear"`
		ZipCode        string      `json:"zipCode"`
		Timezone       string      `json:"timezone"`
		IsMale         bool        `json:"isMale"`
		Weight         int         `json:"weight"`
		Duration       interface{} `json:"duration"`
		SleeperID      string      `json:"sleeperId"`
		Height         int         `json:"height"`
		LicenseVersion int         `json:"licenseVersion"`
		Username       string      `json:"username"`
		BirthMonth     int         `json:"birthMonth"`
		SleepGoal      int         `json:"sleepGoal"`
		IsAccountOwner bool        `json:"isAccountOwner"`
		AccountID      string      `json:"accountId"`
		Email          string      `json:"email"`
		Avatar         string      `json:"avatar"`
		LastLogin      string      `json:"lastLogin"`
		Side           int         `json:"side"`
	} `json:"sleepers"`
	Error ServiceError `json:"Error"`
}

// Sleepers retrieves detailed information about all sleepers (people)
func (s SleepIQ) Sleepers() (SleeperDetails, error) {
	var response SleeperDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleeper?_k={{key}}", "{{key}}", s.loginKey, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper details - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper details - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER ACTIVITY
// ============================================================================

// SleeperActivityDetails describes details about a persons sleep quality.
type SleeperActivityDetails struct {
	Sleepers []struct {
		SleeperID             string `json:"sleeperId"`
		Message               string `json:"message"`
		Tip                   string `json:"tip"`
		AvgHeartRate          int    `json:"avgHeartRate"`
		AvgRespirationRate    int    `json:"avgRespirationRate"`
		TotalSleepSessionTime int    `json:"totalSleepSessionTime"`
		InBed                 int    `json:"inBed"`
		OutOfBed              int    `json:"outOfBed"`
		Restful               int    `json:"restful"`
		Restless              int    `json:"restless"`
		AvgSleepIQ            int    `json:"avgSleepIQ"`
		SleepData             []struct {
			Tip      string `json:"tip"`
			Message  string `json:"message"`
			Date     string `json:"date"`
			Sessions []struct {
				StartDate             string `json:"startDate"`
				Longest               bool   `json:"longest"`
				SleepIQCalculating    bool   `json:"sleepIQCalculating"`
				OriginalStartDate     string `json:"originalStartDate"`
				Restful               int    `json:"restful"`
				OriginalEndDate       string `json:"originalEndDate"`
				SleepNumber           int    `json:"sleepNumber"`
				TotalSleepSessionTime int    `json:"totalSleepSessionTime"`
				AvgHeartRate          int    `json:"avgHeartRate"`
				Restless              int    `json:"restless"`
				AvgRespirationRate    int    `json:"avgRespirationRate"`
				IsFinalized           bool   `json:"isFinalized"`
				SleepQuotient         int    `json:"sleepQuotient"`
				EndDate               string `json:"endDate"`
				OutOfBed              int    `json:"outOfBed"`
				InBed                 int    `json:"inBed"`
			} `json:"sessions"`
			GoalEntry interface{}   `json:"goalEntry"`
			Tags      []interface{} `json:"tags"`
		} `json:"sleepData"`
	} `json:"sleepers"`
	Error ServiceError `json:"Error"`
}

// SleepActivity obtains detailed information about a persons sleep
// quality on a daily basis. Information is returned for a given date
// as well as the length of time. The 'date' must be formatted as
// 'YYYY-MM-DD' but aliases can also be specified. Date aliases
// include: 'today' and 'yesterday'. The 'timeLength' supports the
// following options: 'd1' (one day), 'w1' (one week) or 'm1'
// (one year). Defaults for 'date' are "today" and for 'timeLength'
// are 'd1'.
func (s SleepIQ) SleepActivity(date string, timeLength string) (SleeperActivityDetails, error) {
	var response SleeperActivityDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleepData/?_k={{key}}&date={{date}}&interval={{interval}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{date}}", convertDateAlias(date), -1)
	url = strings.Replace(url, "{{interval}}", convertTimeLength(timeLength), -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper activity - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper activity - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER PREFERENCES
// ============================================================================

// SleeperPreferences contains notifications for sleepers
type SleeperPreferences struct {
	Preferences struct {
		Notifications []interface{} `json:"notifications"`
	} `json:"preferences"`
	SleeperID string       `json:"sleeperId"`
	Error     ServiceError `json:"Error"`
}

// SleeperPreference retrieves preference information for a given sleeper
func (s SleepIQ) SleeperPreference(sleeperID string) (SleeperPreferences, error) {
	var response SleeperPreferences

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleeper/{{sleeperId}}/preferences?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper preferences - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper preferences - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER MONTHLY SUMMARY DETAILS
// ============================================================================

// SleeperMonthlySummaryDetails contains a monthly summary by day for
// each sleeper.
type SleeperMonthlySummaryDetails struct {
	MonthSleepData struct {
		Date string `json:"date"`
		Days []struct {
			Date     string `json:"date"`
			Sleepers []struct {
				SleeperID string `json:"sleeperId"`
				Name      string `json:"name"`
				Session   struct {
					StartDate             string `json:"startDate"`
					Longest               bool   `json:"longest"`
					SleepIQCalculating    bool   `json:"sleepIQCalculating"`
					OriginalStartDate     string `json:"originalStartDate"`
					Restful               int    `json:"restful"`
					OriginalEndDate       string `json:"originalEndDate"`
					SleepNumber           int    `json:"sleepNumber"`
					TotalSleepSessionTime int    `json:"totalSleepSessionTime"`
					AvgHeartRate          int    `json:"avgHeartRate"`
					Restless              int    `json:"restless"`
					AvgRespirationRate    int    `json:"avgRespirationRate"`
					IsFinalized           bool   `json:"isFinalized"`
					SleepQuotient         int    `json:"sleepQuotient"`
					EndDate               string `json:"endDate"`
					OutOfBed              int    `json:"outOfBed"`
					InBed                 int    `json:"inBed"`
				} `json:"session"`
			} `json:"sleepers"`
		} `json:"days"`
		Sleepers []struct {
			SleeperID             string `json:"sleeperId"`
			Message               string `json:"message"`
			AvgSleepIQ            int    `json:"avgSleepIQ"`
			Restful               int    `json:"restful"`
			Tip                   string `json:"tip"`
			TotalSleepSessionTime int    `json:"totalSleepSessionTime"`
			AvgHeartRate          int    `json:"avgHeartRate"`
			Restless              int    `json:"restless"`
			AvgRespirationRate    int    `json:"avgRespirationRate"`
			OutOfBed              int    `json:"outOfBed"`
			InBed                 int    `json:"inBed"`
		} `json:"sleepers"`
	} `json:"monthSleepData"`
	Error ServiceError `json:"Error"`
}

// SleeperMonthlySummary contains a monthly summary by day for each sleeper.
// The 'date' parameter should be in the format of 'YYYY-MM'. However, it
// supports aliases in the form of 'this' (this month), 'last' (last month)
// or specifing the month name (i.e., 'June'). When specifying the month
// name, it will revert to the previous year if the month specified is in
// the future for the current year.
func (s SleepIQ) SleeperMonthlySummary(date string) (SleeperMonthlySummaryDetails, error) {
	var response SleeperMonthlySummaryDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleepData/byMonth?startDate={{date}}&_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{date}}", convertMonthlyDateAlias(date), -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper monthly summary - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper monthly summary - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SLEEPER EDITED SLEEP SESSIONS
// ============================================================================

// EditedSleepSessions describes sleep sessions that have been manually edited
type EditedSleepSessions struct {
	Sleepers []struct {
		EditedSleepSessions []struct {
			EndDate           string `json:"endDate"`
			OriginalEndDate   string `json:"originalEndDate"`
			OriginalStartDate string `json:"originalStartDate"`
			StartDate         string `json:"startDate"`
		} `json:"editedSleepSessions"`
		HiddenSleepSessions []interface{} `json:"hiddenSleepSessions"`
		SleeperID           string        `json:"sleeperId"`
	} `json:"sleepers"`
	Error ServiceError `json:"Error"`
}

// SleeperEditedSessions retrieves information about manually edited
// sleeps sessions for a given sleeper by the provided date range.
func (s SleepIQ) SleeperEditedSessions(sleeperID string, startDate time.Time, endDate time.Time) (EditedSleepSessions, error) {
	var response EditedSleepSessions

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleepData/editedHidden?startDate={{startDate}}&endDate={{endDate}}&sleeperId={{sleeperId}}&_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{startDate}}", startDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{endDate}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper edited sessions - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper edited sessions - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// SleeperNighlyTimeSeriesActivity describes nightly sleep activity in 2.4
// second increments.
type SleeperNighlyTimeSeriesActivity struct {
	Sleepers []struct {
		Days []struct {
			Date      string `json:"date"`
			SliceList []struct {
				OutOfBedTime int `json:"outOfBedTime"`
				RestfulTime  int `json:"restfulTime"`
				RestlessTime int `json:"restlessTime"`
				Type         int `json:"type"`
			} `json:"sliceList"`
		} `json:"days"`
		SleeperID string `json:"sleeperId"`
		SliceSize int    `json:"sliceSize"`
	} `json:"sleepers"`
	Error ServiceError `json:"Error"`
}

// SleeperNightlyDetailedActivity retrieves detailed nightly sleep activity
// for a given sleeper for a specific date. 600 'slices' of time are
// provided per day which equates to 2.4 minutes.
func (s SleepIQ) SleeperNightlyDetailedActivity(sleeperID string, date time.Time) (SleeperNighlyTimeSeriesActivity, error) {
	var response SleeperNighlyTimeSeriesActivity

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/sleepSliceData?_k={{key}}&date={{date}}&sleeper={{sleeperId}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{date}}", date.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{sleeperId}}", sleeperID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve sleeper nightly detailed activity - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read sleeper nightly detailed activity - %s", err)
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

// convertDateAlias converts date aliases to actual dates
func convertDateAlias(alias string) string {
	now := time.Now()

	switch strings.ToLower(alias) {
	case "today":
		return now.Format("2006-01-02")
	case "yesterday":
		dateYesterday := now.AddDate(0, 0, -1)
		return dateYesterday.Format("2006-01-02")
	case "":
		return now.Format("2006-01-02")
	}
	return alias
}

// convertTimeLength converts time lengths to valid durations
func convertTimeLength(alias string) string {
	switch strings.ToLower(alias) {
	case "":
		return "D1"
	}
	return strings.ToUpper(alias)
}

// convertMonthlyDateAlias converts a month date alias to actual dates
func convertMonthlyDateAlias(alias string) string {
	now := time.Now()
	nowMonth := int(now.Month())
	thisYear := now.Format("2006") + "-"
	lastYear := now.AddDate(-1, 0, 0).Format("2006") + "-"

	switch strings.ToLower(alias) {
	case "this":
		return now.Format("2006-01")
	case "last":
		lastMonth := now.AddDate(0, -1, 0)
		return lastMonth.Format("2006-01")
	case "january":
		return thisYear + "01"
	case "february":
		if nowMonth >= 2 {
			return thisYear + "02"
		}
		return lastYear + "02"
	case "march":
		if nowMonth >= 3 {
			return thisYear + "03"
		}
		return lastYear + "03"
	case "april":
		if nowMonth >= 4 {
			return thisYear + "04"
		}
		return lastYear + "04"
	case "may":
		if nowMonth >= 5 {
			return thisYear + "05"
		}
		return lastYear + "05"
	case "june":
		if nowMonth >= 6 {
			return thisYear + "06"
		}
		return lastYear + "06"
	case "july":
		if nowMonth >= 7 {
			return thisYear + "07"
		}
		return lastYear + "07"
	case "august":
		if nowMonth >= 8 {
			return thisYear + "08"
		}
		return lastYear + "08"
	case "september":
		if nowMonth >= 9 {
			return thisYear + "09"
		}
		return lastYear + "09"
	case "october":
		if nowMonth >= 10 {
			return thisYear + "10"
		}
		return lastYear + "10"
	case "november":
		if nowMonth >= 11 {
			return thisYear + "11"
		}
		return lastYear + "11"
	case "december":
		return thisYear + "12"
	}

	return alias
}
