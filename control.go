package sleepiq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type controlResponse struct {
	Error ServiceError `json:"Error"`
}

// ============================================================================
// FOOT WARMER
// ============================================================================

// Footwarmer temperatures
const (
	TempOff    = 0
	TempLow    = 31
	TempMedium = 57
	TempHigh   = 72
)

// ControlFootWarmer sets the foot warmer temperature and duration
// for the given bed and side of bed
func (s SleepIQ) ControlFootWarmer(bedID string, side string, temperature int, duration int) (FootWarmingStatus, error) {
	var response FootWarmingStatus

	// Validate parameters
	if strings.ToLower(side) != "left" && strings.ToLower(side) != "right" {
		return response, errors.New("parameter 'side' must be 'left' or 'right'")
	}

	if temperature != TempLow && temperature != TempMedium && temperature != TempHigh && temperature != TempOff {
		return response, errors.New("parameter 'temperature' must be 'TempOff', 'TempLow', 'TempMedium' or 'TempHigh'")
	}

	if duration < 1 || duration > 360 {
		return response, errors.New("parameter 'duration' must be between 0 and 360 inclusive")
	}

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Create JSON payload
	payload := "{ \"footWarmingTemp{{side}}\": {{temp}}, \"footWarmingTimer{{side}}\": {{duration}} }"
	payload = strings.Replace(payload, "{{side}}", strings.Title(side), -1)
	payload = strings.Replace(payload, "{{temp}}", strconv.Itoa(temperature), -1)
	payload = strings.Replace(payload, "{{duration}}", strconv.Itoa(duration), -1)

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/footwarming?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, _, err := httpPut(url, []byte(payload), s.cookies)
	if err != nil {
		return response, fmt.Errorf("unable to set foot warmer - %s", err)
	}

	// Marshal the response to a loginResponse object
	var footWarmingResponse controlResponse
	err = json.Unmarshal(responseBytes, &footWarmingResponse)
	if err != nil {
		return response, fmt.Errorf("could not read foot warmer response - %s", err)
	}

	// Check for an error returned from the service
	if footWarmingResponse.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", footWarmingResponse.Error.Code, footWarmingResponse.Error.Message)
	}

	// Get the update foot warmer status
	response, err = s.BedFootWarmerStatus(bedID)
	if err != nil {
		return response, fmt.Errorf("could not get bed foot warmer status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read foot warmer response - %s", err)
	}

	// Check for an error returned from the service
	if footWarmingResponse.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", footWarmingResponse.Error.Code, footWarmingResponse.Error.Message)
	}

	return response, nil
}

// ControlFootWarmerOff turns the footwarmer off
func (s SleepIQ) ControlFootWarmerOff(bedID string) (FootWarmingStatus, error) {
	var response FootWarmingStatus

	// Left Side
	response, err := s.ControlFootWarmer(bedID, "Left", TempOff, 120)
	if err != nil {
		return response, fmt.Errorf("could not turn left footwarmer off - %s", err)
	}

	// Right Side
	response, err = s.ControlFootWarmer(bedID, "Right", TempOff, 120)
	if err != nil {
		return response, fmt.Errorf("could not turn right footwarmer off - %s", err)
	}

	return response, nil
}

// ============================================================================
// BED POSITION
// ============================================================================

// bedPosition describes the properties used to set the bed position
type bedPosition struct {
	Speed  int    `json:"speed"`
	Side   string `json:"side"`
	Preset int    `json:"preset"`
}

// Bed Preset Positions
const (
	PositionFavorite = 1
	PositionRead     = 2
	PositionWatchTV  = 3
	PositionFlat     = 4
	PositionZeroG    = 5
	PositionSnore    = 6
)

// ControlBedPosition controls the position of the bed using preset
// bed positions
func (s SleepIQ) ControlBedPosition(bedID string, side string, position int) (BedFoundationStatus, error) {
	var response BedFoundationStatus

	// Validate parameters
	if strings.ToLower(side) != "left" && strings.ToLower(side) != "right" {
		return response, errors.New("parameter 'side' must be 'left' or 'right'")
	}

	if position < 1 || position > 6 {
		return response, errors.New("parameter 'position' must be between 1 and 6 inclusive. Use 'Position' constants")
	}

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Create JSON payload
	payload := bedPosition{
		Preset: position,
		Side:   strings.ToUpper(side[0:1]),
	}
	payloadBytes := new(bytes.Buffer)
	json.NewEncoder(payloadBytes).Encode(payload)

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/preset?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, _, err := httpPut(url, payloadBytes.Bytes(), s.cookies)
	if err != nil {
		return response, fmt.Errorf("unable to set bed position - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return response, fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	// Get the bed foundation status
	response, err = s.BedFoundationStatus(bedID)
	if err != nil {
		return response, fmt.Errorf("could not get bed foundation status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed foundation response - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// LIGHTING
// ============================================================================

// underbedLightOutlet describes the properties that are sent to control
// the underbed lighting outlet
type underbedLightOutlet struct {
	OutletID int    `json:"outletId"`
	Setting  string `json:"setting"`
	Timer    int    `json:"timer"`
}

// underbedLightSystem describes the properties that are sent to control
// the underbed lighting system
type underbedLightSystem struct {
	RightUnderbedLightPWM int `json:"rightUnderbedLightPWM"`
	LeftUnderbedLightPWM  int `json:"leftUnderbedLightPWM"`
}

// Underbed Lighting Levels
const (
	LightLevelLow    = 1
	LightLevelMedium = 30
	LightLevelHigh   = 100
)

// ControlUnderbedLight controls the underbed lighting system
func (s SleepIQ) ControlUnderbedLight(bedID string, lightLevel int, duration int) error {
	// Validate parameters
	if lightLevel != LightLevelLow && lightLevel != LightLevelMedium && lightLevel != LightLevelHigh {
		return errors.New("parameter 'lightLevel' must be 'LightLevelLow', 'LightLevelMedium' or 'LightLevelHigh'")
	}

	if duration < 0 || duration > 180 {
		return errors.New("parameter 'duration' must be between 0 and 180 minutes inclusive")
	}

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return errors.New("user is not logged-in. Please login and try again")
	}

	// Make request - First we need to set the system status
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/system?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	// Create JSON payload
	payload := underbedLightSystem{
		RightUnderbedLightPWM: lightLevel,
		LeftUnderbedLightPWM:  lightLevel,
	}

	payloadBytes := new(bytes.Buffer)
	json.NewEncoder(payloadBytes).Encode(payload)

	responseBytes, _, err := httpPut(url, payloadBytes.Bytes(), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set light system duration - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	// Make request - Last we need to set the outlet status
	url = strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/outlet?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	// Create JSON payload
	payloadOutlet := underbedLightOutlet{
		OutletID: 3,
		Setting:  "1",
		Timer:    duration,
	}

	payloadOutletBytes := new(bytes.Buffer)
	json.NewEncoder(payloadOutletBytes).Encode(payloadOutlet)
	fmt.Println(string(payloadOutletBytes.Bytes()))
	responseBytes, _, err = httpPut(url, payloadOutletBytes.Bytes(), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set light outlet duration - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	return nil
}

// ControlUnderbedLightOff turns the underbed light off
func (s SleepIQ) ControlUnderbedLightOff(bedID string) error {
	return s.ControlUnderbedLight(bedID, LightLevelHigh, 0)
}

// autoUnderbedLight describes the properties that are sent to control
// the underbed lighting automatic mode when a person leaves the bed
type autoUnderbedLight struct {
	EnableAuto bool `json:"enableAuto"`
}

// ControlUnderbedLightAutoMode controls the auto mode of the underbed light
func (s SleepIQ) ControlUnderbedLightAutoMode(bedID string, enabled bool) error {
	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return errors.New("user is not logged-in. Please login and try again")
	}

	// Make request - First we need to set the system status
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/underbedLight?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	// Create JSON payload
	payload := autoUnderbedLight{
		EnableAuto: enabled,
	}

	payloadBytes := new(bytes.Buffer)
	json.NewEncoder(payloadBytes).Encode(payload)

	responseBytes, _, err := httpPut(url, payloadBytes.Bytes(), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set light auto mode - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	return nil
}

// ============================================================================
// RESPONSIVE AIR
// ============================================================================

// responsiveAirMode describes the properties that are sent to control
// the enabledment of the responsive air system
type responsiveAirMode struct {
	RightSideEnabled bool `json:"rightSideEnabled"`
}

// ControlResponsiveAirMode controls the enablement of the responsive air mode
func (s SleepIQ) ControlResponsiveAirMode(bedID string, enabled bool) error {
	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return errors.New("user is not logged-in. Please login and try again")
	}

	// Make request - First we need to set the system status
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/responsiveAir?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	// Create JSON payload
	payload := autoUnderbedLight{
		EnableAuto: enabled,
	}

	payloadBytes := new(bytes.Buffer)
	json.NewEncoder(payloadBytes).Encode(payload)

	responseBytes, _, err := httpPut(url, payloadBytes.Bytes(), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set responsive air mode - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	return nil
}

// ============================================================================
// SLEEP NUMBER
// ============================================================================

//sleepNumberSettings describes the properties that are sent to control
// the sleep number value
type sleepNumberSettings struct {
	Side        string `json:"side"`
	SleepNumber int    `json:"sleepNumber"`
}

// ControlSleepNumber sets the sleep number for the bed
func (s SleepIQ) ControlSleepNumber(bedID string, side string, sleepNumber int) error {
	// Validate Parameters
	if strings.ToLower(side) != "left" && strings.ToLower(side) != "right" {
		return errors.New("parameter 'side' must be 'left' or 'right'")
	}

	if sleepNumber < 1 || sleepNumber > 100 {
		return errors.New("parameter 'sleepNumber' must be between 1 and 100")
	}

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return errors.New("user is not logged-in. Please login and try again")
	}

	// Make request - First we need to set the pump to idle
	err := s.ControlPumpForceIdle(bedID)
	if err != nil {
		return fmt.Errorf("unable to set pump to idle - %s", err)
	}

	// Make request - Last we need to set the sleep number
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/sleepNumber?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	// Create JSON payload
	payload := sleepNumberSettings{
		Side:        side[0:1],
		SleepNumber: sleepNumber,
	}

	payloadBytes := new(bytes.Buffer)
	json.NewEncoder(payloadBytes).Encode(payload)

	responseBytes, _, err := httpPut(url, payloadBytes.Bytes(), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set sleep number - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	return nil
}

// ============================================================================
// PUMP
// ============================================================================

// ControlPumpForceIdle forces the pump to be idle
func (s SleepIQ) ControlPumpForceIdle(bedID string) error {
	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return errors.New("user is not logged-in. Please login and try again")
	}

	// Make request - First we need to set the system status
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/pump/forceIdle?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, _, err := httpPut(url, []byte(""), s.cookies)
	if err != nil {
		return fmt.Errorf("unable to set pump to idle - %s", err)
	}

	// Marshal the response to a loginResponse object
	var controlResponse controlResponse
	err = json.Unmarshal(responseBytes, &controlResponse)
	if err != nil {
		return fmt.Errorf("could not read control response - %s", err)
	}

	// Check for an error returned from the service
	if controlResponse.Error.Code > 0 {
		return fmt.Errorf("error #%d: %s", controlResponse.Error.Code, controlResponse.Error.Message)
	}

	return nil
}
