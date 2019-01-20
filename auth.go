package sleepiq

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// LOGIN
// ============================================================================

// LoginResponse describes the properties returned from a login request to the SleepIQ service
type LoginResponse struct {
	UserID            string       `json:"userId"`
	Key               string       `json:"key"`
	RegistrationState int          `json:"registrationState"`
	EdpLoginStatus    int          `json:"edpLoginStatus"`
	EdpLoginMessage   string       `json:"edpLoginMessage"`
	Error             ServiceError `json:"Error"`
}

type loginCredentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// InsightsLoginResponse describes the properties returned from a login request to the Insights service
type InsightsLoginResponse struct {
	Token     string       `json:"token"`
	SleeperID string       `json:"sleeperId"`
	Error     ServiceError `json:"Error"`
}

// Login authenticates the user against the SleepIQ service
func (s *SleepIQ) Login(username string, password string) (LoginResponse, error) {
	var response LoginResponse
	s.isLoggedIn = false
	s.loginKey = ""

	// Create the loginCredentials object
	creds := loginCredentials{
		Username: username,
		Password: password,
	}

	credBytes, err := json.Marshal(creds)
	if err != nil {
		return response, fmt.Errorf("login could not execute - %s", err)
	}

	// Login request
	responseBytes, cookies, err := httpPut("https://prod-api.sleepiq.sleepnumber.com/rest/login", credBytes, s.cookies)
	if err != nil {
		return response, fmt.Errorf("login failed - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read login response - %s", err)
	}

	// Check if the response contained an error
	if response.Error.Code > 0 {
		return response, fmt.Errorf("Login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
	}

	// Update the sleepiq object
	s.loginKey = response.Key
	s.isLoggedIn = true
	s.cookies = cookies

	return response, nil
}

// ============================================================================
// INSIGHTS LOGIN
// ============================================================================

// InsightsLogin authenticates the user against the SleepIQ Insights
// service. This is a separate service that provides additional
// analysis on a sleeper's sleep behavior
func (s *SleepIQ) InsightsLogin(username string, password string) (InsightsLoginResponse, error) {
	var response InsightsLoginResponse
	s.isInsightsLoggedIn = false
	s.insightsToken = ""

	// Create the loginCredentials object
	creds := loginCredentials{
		Username: username,
		Password: password,
	}

	credBytes, err := json.Marshal(creds)
	if err != nil {
		return response, fmt.Errorf("login could not execute - %s", err)
	}

	// Login request
	responseBytes, err := httpPost("https://sleepiqapi.azure-api.net/prod/accesstoken", credBytes)
	if err != nil {
		return response, fmt.Errorf("login failed - %s", err)
	}

	// Marshal the response to a loginResponse object
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return response, fmt.Errorf("could not read login response - %s", err)
	}

	// Check if the response contained an error
	if response.Error.Code > 0 {
		return response, fmt.Errorf("Login failed - Error #%d: %s", response.Error.Code, response.Error.Message)
	}

	// Update the sleepiq object
	s.insightsToken = response.Token
	s.isInsightsLoggedIn = true

	return response, nil
}
