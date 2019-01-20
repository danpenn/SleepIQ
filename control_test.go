package sleepiq

import (
	"os"
	"testing"
)

func TestControlFootWarmerSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	duration := 3
	temp := TempLow
	footWarmer, err := sleepiq.ControlFootWarmer(beds.Beds[0].BedID, "left", temp, duration)
	if err != nil {
		t.Errorf("could not set bed foot warmer - %s", err)
		return
	}

	if footWarmer.FootWarmingStatusLeft != temp {
		t.Errorf("failed to verify foot warmer temperature control. Expect=%d, Actual=%d", temp, footWarmer.FootWarmingStatusLeft)
		return
	}

	if footWarmer.FootWarmingTimerLeft != duration {
		t.Errorf("failed to verify foot warmer duration control. Expect=%d, Actual=%d", duration, footWarmer.FootWarmingTimerLeft)
		return
	}

	footWarmer, err = sleepiq.ControlFootWarmerOff(beds.Beds[0].BedID) // Off
	if err != nil {
		t.Errorf("could not set bed foot warmer off - %s", err)
		return
	}
}

func TestControlBedPositionSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	_, err = sleepiq.ControlBedPosition(beds.Beds[0].BedID, "Right", PositionFlat)
	if err != nil {
		t.Errorf("could not set bed position - %s", err)
		return
	}
}

func TestControlBedLightSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	err = sleepiq.ControlUnderbedLight(beds.Beds[0].BedID, LightLevelMedium, 5)
	if err != nil {
		t.Errorf("could not set bed light - %s", err)
		return
	}

	err = sleepiq.ControlUnderbedLightOff(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not set bed light - %s", err)
		return
	}
}

func TestControlLightAutoModeSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	err = sleepiq.ControlUnderbedLightAutoMode(beds.Beds[0].BedID, false)
	if err != nil {
		t.Errorf("could not set bed light auto mode - %s", err)
		return
	}
}

func TestControlResponsiveAirModeSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	err = sleepiq.ControlResponsiveAirMode(beds.Beds[0].BedID, false)
	if err != nil {
		t.Errorf("could not set responsive air mode - %s", err)
		return
	}
}

func TestControlForcePumpIdleSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	err = sleepiq.ControlPumpForceIdle(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not set pump to idle - %s", err)
		return
	}
}

func TestControlSleepNumberSuccess(t *testing.T) {
	sleepiq := New()

	response, err := sleepiq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	err = sleepiq.ControlSleepNumber(beds.Beds[0].BedID, "Left", 55)
	if err != nil {
		t.Errorf("could not set sleep number - %s", err)
		return
	}
}
