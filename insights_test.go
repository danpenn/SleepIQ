package sleepiq

import (
	"os"
	"testing"
	"time"
)

func TestInsightsActivitySuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.InsightsLogin(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Errorf("login failed - expected success. %s", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	now := time.Now()
	activity, err := sleepiq.InsightsActiviy(response.SleeperID, now.AddDate(0, 0, -7), now) // 1 week
	if err != nil {
		t.Error("could not get Insights activity")
		return
	}

	if len(activity.Activities) == 0 {
		t.Error("no statuses were found in Insights activity")
	}
}

func TestInsightsProvidersSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.InsightsLogin(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Errorf("login failed - expected success. %s", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	providers, err := sleepiq.InsightsProviders()
	if err != nil {
		t.Errorf("could not get Insights providers - %s", err)
		return
	}

	if len(providers.Providers) == 0 {
		t.Error("no statuses were found in Insights providers")
	}
}

func TestInsightsLikeMeNearMeAndMeSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.InsightsLogin(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Errorf("login failed - expected success. %s", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	now := time.Now()
	likeMe, err := sleepiq.InsightsLikeMe(response.SleeperID, now.AddDate(-2, 0, 0), now)
	if err != nil {
		t.Errorf("could not get Insights Like Me - %s", err)
		return
	}

	if len(likeMe.Data) == 0 {
		t.Error("no data was found in Insights Like Me")
	}

	nearMe, err := sleepiq.InsightsNearMe(response.SleeperID, now.AddDate(-2, 0, 0), now)
	if err != nil {
		t.Errorf("could not get Insights Near Me - %s", err)
		return
	}

	if len(nearMe.Data) == 0 {
		t.Error("no data was found in Insights Near Me")
	}

	me, err := sleepiq.InsightsMe(response.SleeperID, now.AddDate(-2, 0, 0), now)
	if err != nil {
		t.Errorf("could not get Insights Near Me - %s", err)
		return
	}

	if len(me.Data) == 0 {
		t.Error("no data was found in Insights Near Me")
	}
}
