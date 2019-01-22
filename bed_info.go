package sleepiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Bed Types
const (
	BedTypeSingle      = 0
	BedTypeSplitHead   = 1
	BedTypeSplitKing   = 2
	BedTypeEasternKing = 3
)

// Sides of the bed
const (
	BedSideLeft  = 0
	BedSideRight = 1
)

// ============================================================================
// BEDS
// ============================================================================

// BedsInfo describes the beds returned from the service
type BedsInfo struct {
	Beds  []Bed        `json:"beds"`
	Error ServiceError `json:"Error"`
}

// Bed describes the details of a bed
type Bed struct {
	RegistrationDate    time.Time   `json:"registrationDate"`
	SleeperRightID      string      `json:"sleeperRightId"`
	Base                interface{} `json:"base"`
	ReturnRequestStatus int         `json:"returnRequestStatus"`
	Size                string      `json:"size"`
	Name                string      `json:"name"`
	Serial              string      `json:"serial"`
	IsKidsBed           bool        `json:"isKidsBed"`
	DualSleep           bool        `json:"dualSleep"`
	BedID               string      `json:"bedId"`
	Status              int         `json:"status"`
	SleeperLeftID       string      `json:"sleeperLeftId"`
	Version             string      `json:"version"`
	AccountID           string      `json:"accountId"`
	Timezone            string      `json:"timezone"`
	Generation          string      `json:"generation"`
	Model               string      `json:"model"`
	PurchaseDate        time.Time   `json:"purchaseDate"`
	MacAddress          string      `json:"macAddress"`
	Sku                 string      `json:"sku"`
	Zipcode             string      `json:"zipcode"`
	Reference           string      `json:"reference"`
}

// Beds returns properties about all beds associated with the account
func (s SleepIQ) Beds() (BedsInfo, error) {
	var response BedsInfo

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed?_k={{key}}", "{{key}}", s.loginKey, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed details - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read beds response - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// PRIVACY MODE
// ============================================================================

// BedPrivacyModeDetails describes whether privacy mode is on or off
type BedPrivacyModeDetails struct {
	AccountID string       `json:"accountId"`
	BedID     string       `json:"bedId"`
	PauseMode string       `json:"pauseMode"`
	Error     ServiceError `json:"Error"`
}

// BedPrivacyMode gets the privacy mode for the specified bed. The bedID can be obtained via the call
// to Beds().
func (s SleepIQ) BedPrivacyMode(bedID string) (BedPrivacyModeDetails, error) {
	var response BedPrivacyModeDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/pauseMode?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed pause mode - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed pause mode - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// FAMILY STATUS
// ============================================================================

// FamilyStatusDetails describes the bed settings for each bed
type FamilyStatusDetails struct {
	Beds []struct {
		Status   int    `json:"status"`
		BedID    string `json:"bedId"`
		LeftSide struct {
			IsInBed              bool   `json:"isInBed"`
			AlertDetailedMessage string `json:"alertDetailedMessage"`
			SleepNumber          int    `json:"sleepNumber"`
			AlertID              int    `json:"alertId"`
			LastLink             string `json:"lastLink"`
			Pressure             int    `json:"pressure"`
		} `json:"leftSide"`
		RightSide struct {
			IsInBed              bool   `json:"isInBed"`
			AlertDetailedMessage string `json:"alertDetailedMessage"`
			SleepNumber          int    `json:"sleepNumber"`
			AlertID              int    `json:"alertId"`
			LastLink             string `json:"lastLink"`
			Pressure             int    `json:"pressure"`
		} `json:"rightSide"`
	} `json:"beds"`
	Error ServiceError `json:"Error"`
}

// BedFamilyStatus gets the settings for each bed as well as each side of a bed (if applicable)
func (s SleepIQ) BedFamilyStatus() (FamilyStatusDetails, error) {
	var response FamilyStatusDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/familyStatus?_k={{key}}", "{{key}}", s.loginKey, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed family status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed family status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// DETAILED STATUS
// ============================================================================

// BedDetailedInfo describes detailed settings for the given bed. This information is helpful
// for troubleshooting problems.
type BedDetailedInfo struct {
	BedID    string `json:"bedId"`
	Chambers struct {
		LeftChamberOccupancy       interface{} `json:"leftChamberOccupancy"`
		LeftChamberRefreshedState  interface{} `json:"leftChamberRefreshedState"`
		LeftChamberType            int         `json:"leftChamberType"`
		RightChamberOccupancy      interface{} `json:"rightChamberOccupancy"`
		RightChamberRefreshedState interface{} `json:"rightChamberRefreshedState"`
		RightChamberType           int         `json:"rightChamberType"`
	} `json:"chambers"`
	Foundation struct {
		FsCurrentPositionPresetLeft  string `json:"fsCurrentPositionPresetLeft"`
		FsCurrentPositionPresetRight string `json:"fsCurrentPositionPresetRight"`
		FsType                       string `json:"fsType"`
		Outlets                      []struct {
			OutletID int         `json:"outletId"`
			Setting  interface{} `json:"setting"`
		} `json:"outlets"`
	} `json:"foundation"`
	Pump struct {
		ActiveTask               int `json:"activeTask"`
		ChamberType              int `json:"chamberType"`
		LeftSideSleepNumber      int `json:"leftSideSleepNumber"`
		RightSideSleepNumber     int `json:"rightSideSleepNumber"`
		SleepNumberFavoriteLeft  int `json:"sleepNumberFavoriteLeft"`
		SleepNumberFavoriteRight int `json:"sleepNumberFavoriteRight"`
	} `json:"pump"`
	Smartoutlets []struct {
		Name     string `json:"name"`
		OutletID int    `json:"outletId"`
		Setting  int    `json:"setting"`
	} `json:"smartoutlets"`
	Error ServiceError `json:"Error"`
}

// BedDetailedStatus gets the settings for each bed as well as each side of a bed (if applicable)
func (s SleepIQ) BedDetailedStatus(bedID string) (BedDetailedInfo, error) {
	var response BedDetailedInfo

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/superStatus?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed detailed status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed detailed status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// NODES
// ============================================================================

// BedNodesDetails describes the nodes (whatever this means)
type BedNodesDetails struct {
	BedID string       `json:"bedId"`
	Nodes []int        `json:"nodes"`
	Error ServiceError `json:"Error"`
}

// BedNodes gets the nodes for the provided bed
func (s SleepIQ) BedNodes(bedID string) (BedNodesDetails, error) {
	var response BedNodesDetails

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/nodes?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed nodes - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed nodes - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// RESPONSIVE AIR
// ============================================================================

// ResponsiveAirSettings describes the settings for the responsive air system
type ResponsiveAirSettings struct {
	AdjustmentThreshold int          `json:"adjustmentThreshold"`
	InBedTimeout        int          `json:"inBedTimeout"`
	LeftSideEnabled     bool         `json:"leftSideEnabled"`
	OutOfBedTimeout     int          `json:"outOfBedTimeout"`
	PollFrequency       int          `json:"pollFrequency"`
	PrefSyncState       string       `json:"prefSyncState"`
	RightSideEnabled    bool         `json:"rightSideEnabled"`
	Error               ServiceError `json:"Error"`
}

// BedResponsiveAir gets the responsive air settings for the provided bed
func (s SleepIQ) BedResponsiveAir(bedID string) (ResponsiveAirSettings, error) {
	var response ResponsiveAirSettings

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/responsiveAir?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed responsive error settings - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed responsive error settings - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// FOOT WARMER STATUS
// ============================================================================

// FootWarmingStatus describes the status of the beds foot warmer
type FootWarmingStatus struct {
	FootWarmingStatusLeft  int          `json:"footWarmingStatusLeft"`
	FootWarmingStatusRight int          `json:"footWarmingStatusRight"`
	FootWarmingTimerLeft   int          `json:"footWarmingTimerLeft"`
	FootWarmingTimerRight  int          `json:"footWarmingTimerRight"`
	Error                  ServiceError `json:"Error"`
}

// BedFootWarmerStatus retrieves the foot warmer status for the bed
func (s SleepIQ) BedFootWarmerStatus(bedID string) (FootWarmingStatus, error) {
	var response FootWarmingStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/footwarming?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed foot warmer status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed foot warmer status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// SYSTEM STATUS
// ============================================================================

// BedSystemStatus describes the status of the bed computer board and lighting systems
type BedSystemStatus struct {
	BedType               int          `json:"fsBedType"`
	BoardFaults           int          `json:"fsBoardFaults"`
	BoardFeatures         int          `json:"fsBoardFeatures"`
	BoardHWRevisionCode   int          `json:"fsBoardHWRevisionCode"`
	BoardStatus           int          `json:"fsBoardStatus"`
	LeftUnderbedLightPWM  int          `json:"fsLeftUnderbedLightPWM"`
	RightUnderbedLightPWM int          `json:"fsRightUnderbedLightPWM"`
	Error                 ServiceError `json:"Error"`
}

// BedSystemStatus retrieves the board and lighting status of the bed
func (s SleepIQ) BedSystemStatus(bedID string) (BedSystemStatus, error) {
	var response BedSystemStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/system?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed system status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed system status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// PINCH STATUS
// ============================================================================

// BedPinchStatus describes the pinch status of the (whatever that means)
type BedPinchStatus struct {
	ContinuousPinchLeftFoot         bool         `json:"continuousPinchLeftFoot"`
	ContinuousPinchLeftHead         bool         `json:"continuousPinchLeftHead"`
	ContinuousPinchRightFoot        bool         `json:"continuousPinchRightFoot"`
	ContinuousPinchRightHead        bool         `json:"continuousPinchRightHead"`
	PinchEventsLeftFoot             int          `json:"pinchEventsLeftFoot"`
	PinchEventsLeftHead             int          `json:"pinchEventsLeftHead"`
	PinchEventsRightFoot            int          `json:"pinchEventsRightFoot"`
	PinchEventsRightHead            int          `json:"pinchEventsRightHead"`
	PinchSenseDisconnectedLeftFoot  bool         `json:"pinchSenseDisconnectedLeftFoot"`
	PinchSenseDisconnectedLeftHead  bool         `json:"pinchSenseDisconnectedLeftHead"`
	PinchSenseDisconnectedRightFoot bool         `json:"pinchSenseDisconnectedRightFoot"`
	PinchSenseDisconnectedRightHead bool         `json:"pinchSenseDisconnectedRightHead"`
	Error                           ServiceError `json:"Error"`
}

// BedPinchStatus retrieves the pinch status of the bed
func (s SleepIQ) BedPinchStatus(bedID string) (BedPinchStatus, error) {
	var response BedPinchStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/pinch?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed pinch status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed pinch status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// UNDERBED LIGHT STATUS
// ============================================================================

// UnderbedLightStatus describes the status of the underbed light
type UnderbedLightStatus struct {
	EnableAuto    bool         `json:"enableAuto"`
	PrefSyncState string       `json:"prefSyncState"`
	Error         ServiceError `json:"Error"`
}

// BedLightStatus retrieves the status of the underbed light
func (s SleepIQ) BedLightStatus(bedID string) (UnderbedLightStatus, error) {
	var response UnderbedLightStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/underbedLight?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed light status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed light status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// FOUNDATION STATUS
// ============================================================================

// BedFoundationStatus describes the status of the bed foundation, particularly
// in regards to it's position
type BedFoundationStatus struct {
	CurrentPositionPresetRight   string       `json:"fsCurrentPositionPresetRight"`
	NeedsHoming                  bool         `json:"fsNeedsHoming"`
	RightFootPosition            string       `json:"fsRightFootPosition"`
	LeftPositionTimerLSB         string       `json:"fsLeftPositionTimerLSB"`
	TimerPositionPresetLeft      string       `json:"fsTimerPositionPresetLeft"`
	CurrentPositionPresetLeft    string       `json:"fsCurrentPositionPresetLeft"`
	LeftPositionTimerMSB         string       `json:"fsLeftPositionTimerMSB"`
	RightFootActuatorMotorStatus string       `json:"fsRightFootActuatorMotorStatus"`
	CurrentPositionPreset        string       `json:"fsCurrentPositionPreset"`
	TimerPositionPresetRight     string       `json:"fsTimerPositionPresetRight"`
	Type                         string       `json:"fsType"`
	OutletsOn                    bool         `json:"fsOutletsOn"`
	LeftHeadPosition             string       `json:"fsLeftHeadPosition"`
	IsMoving                     bool         `json:"fsIsMoving"`
	RightHeadActuatorMotorStatus string       `json:"fsRightHeadActuatorMotorStatus"`
	StatusSummary                string       `json:"fsStatusSummary"`
	TimerPositionPreset          string       `json:"fsTimerPositionPreset"`
	LeftFootPosition             string       `json:"fsLeftFootPosition"`
	RightPositionTimerLSB        string       `json:"fsRightPositionTimerLSB"`
	TimedOutletsOn               bool         `json:"fsTimedOutletsOn"`
	RightHeadPosition            string       `json:"fsRightHeadPosition"`
	Configured                   bool         `json:"fsConfigured"`
	RightPositionTimerMSB        string       `json:"fsRightPositionTimerMSB"`
	LeftHeadActuatorMotorStatus  string       `json:"fsLeftHeadActuatorMotorStatus"`
	LeftFootActuatorMotorStatus  string       `json:"fsLeftFootActuatorMotorStatus"`
	Error                        ServiceError `json:"Error"`
}

// BedFoundationStatus retrieves the status of the bed foundation
func (s SleepIQ) BedFoundationStatus(bedID string) (BedFoundationStatus, error) {
	var response BedFoundationStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/status?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed foundation status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed foundation status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// UNDERBED LIGHT OUTLET STATUS
// ============================================================================

// UnderbedLightOutletStatus describes the status of the underbed light outlet
type UnderbedLightOutletStatus struct {
	BedID   string       `json:"bedId"`
	Outlet  int          `json:"outlet"`
	Setting int          `json:"setting"`
	Timer   interface{}  `json:"timer"`
	Error   ServiceError `json:"Error"`
}

// BedLightingOutletStatus retrieves the status of the underbed lighting outlet
func (s SleepIQ) BedLightingOutletStatus(bedID string, outletID int) (UnderbedLightOutletStatus, error) {
	var response UnderbedLightOutletStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/outlet?outletId={{outletId}}&_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)
	url = strings.Replace(url, "{{outletId}}", strconv.Itoa(outletID), -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed lighting outlet status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed lighting outlet status - %s", err)
	}

	// Check for an error returned from the service
	if response.Error.Code > 0 {
		return response, fmt.Errorf("error #%d: %s", response.Error.Code, response.Error.Message)
	}

	return response, nil
}

// ============================================================================
// UNDERBED LIGHT SYSTEM STATUS
// ============================================================================

// UnderbedLightSystemStatus describes the status of the underbed lighting system
type UnderbedLightSystemStatus struct {
	BedType               int          `json:"fsBedType"`
	BoardFaults           int          `json:"fsBoardFaults"`
	BoardFeatures         int          `json:"fsBoardFeatures"`
	BoardHWRevisionCode   int          `json:"fsBoardHWRevisionCode"`
	BoardStatus           int          `json:"fsBoardStatus"`
	LeftUnderbedLightPWM  int          `json:"fsLeftUnderbedLightPWM"`
	RightUnderbedLightPWM int          `json:"fsRightUnderbedLightPWM"`
	Error                 ServiceError `json:"Error"`
}

// BedLightingSystemStatus retrieves the status of the underbed lighting system
func (s SleepIQ) BedLightingSystemStatus(bedID string) (UnderbedLightSystemStatus, error) {
	var response UnderbedLightSystemStatus

	// Bail if there is not an active logged-in session
	if !s.isLoggedIn {
		return response, errors.New("user is not logged-in. Please login and try again")
	}

	// Make request
	url := strings.Replace("https://prod-api.sleepiq.sleepnumber.com/rest/bed/{{bedId}}/foundation/system?_k={{key}}", "{{key}}", s.loginKey, -1)
	url = strings.Replace(url, "{{bedId}}", bedID, -1)

	responseBytes, err := httpGet(url, s.cookies, getHeaders())
	if err != nil {
		return response, fmt.Errorf("unable to retrieve bed lighting system status - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read bed lighting system status - %s", err)
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

// getHeaders gets the headers required by the sleepiq service API
func getHeaders() map[string]string {
	var headers map[string]string
	headers = make(map[string]string)

	headers["Accept"] = "application/json"

	return headers
}
