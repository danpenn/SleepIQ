package sleepiq

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestBedsSuccess(t *testing.T) {
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

	beds, err := sleepiq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	for _, bed := range beds.Beds {
		fmt.Printf("%s (%s)\n", bed.Size, bed.Name)
	}
}

func TestBedsUserNotLoggedIn(t *testing.T) {
	sleepiq := New()

	_, err := sleepiq.Beds()
	if err != nil {
		if strings.Contains(err.Error(), "user is not logged-in") {
			return
		}
	}
	t.Error("user shouldn't have been logged in")
}

func TestBedPrivacyModeSuccess(t *testing.T) {
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

	_, err = sleepiq.BedPrivacyMode(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed pause mode - %s", err)
		return
	}
}

func TestBedFamilyStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedFamilyStatus()
	if err != nil {
		t.Errorf("could not get bed family status - %s", err)
		return
	}
}

func TestBedDetailedStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedDetailedStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedNodesSuccess(t *testing.T) {
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

	_, err = sleepiq.BedNodes(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedResponsiveAirSuccess(t *testing.T) {
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

	_, err = sleepiq.BedResponsiveAir(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedFootWarmerStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedFootWarmerStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed foot warmer status - %s", err)
		return
	}
}

func TestBedSystemStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedSystemStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedSystemPinchSuccess(t *testing.T) {
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

	_, err = sleepiq.BedPinchStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedLightStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedLightStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestBedFoundationStatusSuccess(t *testing.T) {
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

	_, err = sleepiq.BedFoundationStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}
}

func TestUnderbedOutletSystemSuccess(t *testing.T) {
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

	_, err = sleepiq.BedLightingOutletStatus(beds.Beds[0].BedID, 3)
	if err != nil {
		t.Errorf("could not get bed light outlet status - %s", err)
		return
	}

	_, err = sleepiq.BedLightingSystemStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed light system status - %s", err)
		return
	}
}
