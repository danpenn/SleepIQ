package sleepiq

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestBedSuccess(t *testing.T) {
	siq := New()

	// Test Beds() - Not Logged In
	_, err := siq.Beds()
	if err != nil {
		if strings.Contains(err.Error(), "user is not logged-in") {
			return
		}
	}
	t.Error("user shouldn't have been logged in")

	response, err := siq.Login(os.Getenv("sleepiq_username"), os.Getenv("sleepiq_password"))
	if err != nil {
		t.Error("login failed - expected success", err)
		return
	}

	if response.Error.Code > 0 {
		t.Errorf("login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
		return
	}

	// Get the beds so we can retrieve the bedID
	beds, err := siq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	if len(beds.Beds) == 0 {
		t.Error("no beds were found in the account")
		return
	}

	// Test Beds()
	beds, err = siq.Beds()
	if err != nil {
		t.Errorf("could not get beds - %s", err)
		return
	}

	for _, bed := range beds.Beds {
		fmt.Printf("%s (%s)\n", bed.Size, bed.Name)
	}

	// Test BedPrivacyMode()
	privacyMode, err := siq.BedPrivacyMode(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed pause mode - %s", err)
		return
	}

	if privacyMode.PauseMode != "off" && privacyMode.PauseMode != "on" {
		t.Error("Privacy mode was neither off nor on")
	}

	// Test BedFamilyStatus
	familyStatus, err := siq.BedFamilyStatus()
	if err != nil {
		t.Errorf("could not get bed family status - %s", err)
		return
	}

	if familyStatus.Beds[0].Status < 0 {
		t.Errorf("bed family status is invalid")
	}

	// Test BedDetailedStatus()
	detailedStatus, err := siq.BedDetailedStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
		return
	}

	if detailedStatus.BedID == "" {
		t.Errorf("could not get a valid detailed status")
	}

	// Test BedNodes()
	bedNodes, err := siq.BedNodes(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed nodes status - %s", err)
	}

	if bedNodes.BedID == "" {
		t.Errorf("could not get a valid bed nodes info")
	}

	// Test BedResponsiveAir()
	responsiveAir, err := siq.BedResponsiveAir(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed responsive air status - %s", err)
	}

	if responsiveAir.InBedTimeout == 0 {
		t.Errorf("could not get bed responsive air status data - %s", err)
	}

	// Test BedFootWarmerStatus()
	_, err = siq.BedFootWarmerStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed foot warmer status - %s", err)
	}

	// Test BedSystemStatus()
	bedSystemStatus, err := siq.BedSystemStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed system status - %s", err)
	}

	if bedSystemStatus.BoardHWRevisionCode == 0 {
		t.Errorf("bed system status is invalid")
	}

	// Test BedPinchStatus()
	_, err = siq.BedPinchStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed pinch status - %s", err)
	}

	// Test BedLightStatus
	_, err = siq.BedLightStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed light status - %s", err)
	}

	// Test BedFoundationStatus()
	foundationStatus, err := siq.BedFoundationStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed detailed status - %s", err)
	}

	if foundationStatus.CurrentPositionPresetRight == "" {
		t.Errorf("foundation status is invalid")
	}

	// Test BedLightingOutletStatus
	outletStatus, err := siq.BedLightingOutletStatus(beds.Beds[0].BedID, 3)
	if err != nil {
		t.Errorf("could not get bed light outlet status - %s", err)
		return
	}

	if outletStatus.BedID == "" {
		t.Errorf("bed light outlet status is invalid")
	}

	// Test BedLightingSystemStatus()
	outletSystem, err := siq.BedLightingSystemStatus(beds.Beds[0].BedID)
	if err != nil {
		t.Errorf("could not get bed light system status - %s", err)
		return
	}

	if outletSystem.BoardHWRevisionCode == 0 {
		t.Errorf("bed light outlet system status is invalid")
	}
}
