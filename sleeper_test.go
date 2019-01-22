package sleepiq

import (
	"os"
	"testing"
	"time"
)

func TestSleepersSuccess(t *testing.T) {
	siq := New()

	response, err := siq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Test Sleepers()
	sleepers, err := siq.Sleepers()
	if err != nil {
		t.Errorf("could not get sleeper details - %s", err)
		return
	}

	if len(sleepers.Sleepers) == 0 {
		t.Error("No sleepers were found")
	}

	// Test SleepActivity() -- TODAY
	activity, err := siq.SleepActivity("today", "")
	if err != nil {
		t.Errorf("could not get sleeper activity - %s", err)
		return
	}

	if len(activity.Sleepers) == 0 {
		t.Error("warning: no sleep data was available for the date test. Skipping test!")
	}

	if len(activity.Sleepers[0].SleepData) == 0 {
		t.Error("warning: no sleep data was available for the date test. Skipping test!")
	}

	// Validate date - date returned should be 'today'
	today := convertDateAlias("today")
	actualDate := activity.Sleepers[0].SleepData[0].Date
	if actualDate != today {
		t.Errorf("failed to get correct date for sleeper activity. Expected=%s, Actual=%s", today, actualDate)
	}

	// Validate date - days returned should be one days worth of data
	if len(activity.Sleepers[0].SleepData) != 1 {
		t.Errorf("failed to get correct date for sleeper activity. Expected=%d, Actual=%d", 1, len(activity.Sleepers[0].SleepData))
	}

	// Test SleepActivity() -- YESTERDAY
	activity, err = siq.SleepActivity("yesterday", "W1")
	if err != nil {
		t.Errorf("could not get sleeper activity - %s", err)
		return
	}

	if len(activity.Sleepers) == 0 {
		t.Error("warning: no sleep data was available for the date test. Skipping test!")
	}

	if len(activity.Sleepers[0].SleepData) == 0 {
		t.Error("warning: no sleep data was available for the date test. Skipping test!")
	}

	// Validate date - date returned should be 'yesterday'
	today = convertDateAlias("yesterday")
	actualDate = activity.Sleepers[0].SleepData[0].Date
	if actualDate != today {
		t.Errorf("failed to get correct date for sleeper activity. Expected=%s, Actual=%s", today, actualDate)
	}

	// Validate date - days returned should be one days worth of data
	if len(activity.Sleepers[0].SleepData) != 7 {
		t.Errorf("failed to get correct date for sleeper activity. Expected=%d, Actual=%d", 1, len(activity.Sleepers[0].SleepData))
	}

	// Test SleeperPreferences()
	preferences, err := siq.SleeperPreference(sleepers.Sleepers[0].SleeperID)
	if err != nil {
		t.Errorf("could not get sleeper preferences - %s", err)
		return
	}

	if sleepers.Sleepers[0].SleeperID != preferences.SleeperID {
		t.Errorf("failed to verify sleeperID. Expect=%s, Actual=%s", sleepers.Sleepers[0].SleeperID, preferences.SleeperID)
	}

	// Test SleeperMonthlySummary
	monthlySummary, err := siq.SleeperMonthlySummary("this") // current month
	if err != nil {
		t.Errorf("could not get sleeper monthly summary - %s", err)
		return
	}

	if monthlySummary.MonthSleepData.Date != convertMonthlyDateAlias("this") {
		t.Errorf("failed to verify sleeper monthly summary date. Expect=%s, Actual=%s", convertMonthlyDateAlias("this"), monthlySummary.MonthSleepData.Date)
	}

	now := time.Now()
	monthlySummary, err = siq.SleeperMonthlySummary(now.Format("January")) // current month name
	if err != nil {
		t.Errorf("could not get sleeper monthly summary - %s", err)
		return
	}

	if monthlySummary.MonthSleepData.Date != convertMonthlyDateAlias(now.Format("January")) {
		t.Errorf("failed to verify sleeper monthly summary date. Expect=%s, Actual=%s", convertMonthlyDateAlias(now.Format("January")), monthlySummary.MonthSleepData.Date)
	}

	// Test SleeperEditedSessions()
	editedSessions, err := siq.SleeperEditedSessions(sleepers.Sleepers[0].SleeperID, now.AddDate(0, 0, -1), now)
	if err != nil {
		t.Errorf("could not get sleeper edited sessions - %s", err)
		return
	}

	if editedSessions.Sleepers[0].SleeperID != sleepers.Sleepers[0].SleeperID {
		t.Errorf("failed to verify sleeper ID's match. Expect=%s, Actual=%s", sleepers.Sleepers[0].SleeperID, editedSessions.Sleepers[0].SleeperID)
	}

	// Test SleeperNightlyDetailedActivity()
	nightlyActivity, err := siq.SleeperNightlyDetailedActivity(sleepers.Sleepers[0].SleeperID, now.AddDate(0, 0, -1)) // yesterday
	if err != nil {
		t.Errorf("could not get sleeper edited sessions - %s", err)
		return
	}

	if nightlyActivity.Sleepers[0].SleeperID != sleepers.Sleepers[0].SleeperID {
		t.Errorf("failed to verify sleeper ID's match. Expect=%s, Actual=%s", sleepers.Sleepers[0].SleeperID, activity.Sleepers[0].SleeperID)
	}
}

func TestDateAliasConversionToday(t *testing.T) {
	actualDate := convertDateAlias("today")

	now := time.Now()
	expectedDate := now.Format("2006-01-02")
	if actualDate != expectedDate {
		t.Errorf("date alias conversion failed. Expected=%s, Actual=%s", expectedDate, actualDate)
	}
}

func TestDateAliasConversionYesterday(t *testing.T) {
	actualDate := convertDateAlias("yesterday")

	now := time.Now()
	yesterdayDate := now.AddDate(0, 0, -1)
	expectedDate := yesterdayDate.Format("2006-01-02")

	if actualDate != expectedDate {
		t.Errorf("date alias conversion failed. Expected=%s, Actual=%s", expectedDate, actualDate)
	}
}

func TestDateAliasConversionDefault(t *testing.T) {
	alias := ""
	conversionResult := convertDateAlias(alias)

	now := time.Now()
	expectedDate := now.Format("2006-01-02")

	if expectedDate != conversionResult {
		t.Errorf("date alias conversion failed. Expected=%s, Actual=%s", expectedDate, conversionResult)
	}
}

func TestDateAliasConversionNoMatch(t *testing.T) {
	alias := "foobar"
	conversionResult := convertDateAlias(alias)

	if alias != conversionResult {
		t.Errorf("date alias conversion failed. Expected=%s, Actual=%s", alias, conversionResult)
	}
}

func TestTimeLengthConversionDefault(t *testing.T) {
	alias := ""
	conversionResult := convertTimeLength(alias)

	if conversionResult != "D1" {
		t.Errorf("time length alias conversion failed. Expected=%s, Actual=%s", "D1", conversionResult)
	}
}

func TestTimeLengthConversionNone(t *testing.T) {
	alias := "M1"
	conversionResult := convertTimeLength(alias)

	if conversionResult != "M1" {
		t.Errorf("time length alias conversion failed. Expected=%s, Actual=%s", "M1", conversionResult)
	}
}

func TestTimeLengthConversionUpperCase(t *testing.T) {
	alias := "w1"
	conversionResult := convertTimeLength(alias)

	if conversionResult != "W1" {
		t.Errorf("time length alias conversion failed. Expected=%s, Actual=%s", "W1", conversionResult)
	}
}

func TestConvertMonthlyDateAliasThis(t *testing.T) {
	now := time.Now()
	expectedDate := now.Format("2006-01")
	actualDate := convertMonthlyDateAlias("this")

	if actualDate != expectedDate {
		t.Errorf("monthly date alias conversion failed. Expected=%s, Actual=%s", expectedDate, actualDate)
	}
}

func TestConvertMonthlyDateAliasLast(t *testing.T) {
	now := time.Now()
	expectedDate := now.AddDate(0, -1, 0).Format("2006-01")
	actualDate := convertMonthlyDateAlias("last")

	if actualDate != expectedDate {
		t.Errorf("monthly date alias conversion failed. Expected=%s, Actual=%s", expectedDate, actualDate)
	}
}

func TestConvertMonthlyDateAliasCurrentMonth(t *testing.T) {
	now := time.Now()
	expectedDate := now.Format("2006") + "-01"
	actualDate := convertMonthlyDateAlias(now.Format("January"))

	if actualDate != expectedDate {
		t.Errorf("monthly date alias conversion failed. Expected=%s, Actual=%s", expectedDate, actualDate)
	}
}
