package sleepiq

import (
	"os"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	sleepiq := New()
	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("Login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("Login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}
}

func TestLoginBadCredentials(t *testing.T) {
	sleepiq := New()
	response, err := sleepiq.Login("JohnDoe@live.com", "bogusPassword")
	if err == nil {
		t.Error("Login succeeded - expected failure", err)
		return
	}

	if response.Error.Code == 0 {
		t.Errorf("Login succeeded - expect failure. Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}
}

func TestInsightsLoginSuccess(t *testing.T) {
	sleepiq := New()
	response, err := sleepiq.InsightsLogin(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("Login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("Login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}
}
