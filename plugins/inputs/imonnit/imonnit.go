package imonnit

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// Description returns the plugin description.
func (s *Sensor) Description() string {
	return "Read stats from the iMonnit API for a given user account"
}

// SampleConfig returns sample configuration for this plugin.
func (s *Sensor) SampleConfig() string {
	return sampleConfig
}

// defaultServer is the default location to forward any imonnit queries
const defaultServer = "https://www.imonnit.com/json"

// NewiMonnitSensor returns a new instance of the iMonnit collector
func NewiMonnitSensor() *Sensor {
	return &Sensor{
		HttpTimeout: internal.Duration{Duration: time.Second * 5},
		Server:      defaultServer,
	}
}

func init() {
	inputs.Add("imonnit", func() telegraf.Input {
		return NewiMonnitSensor()
	})
}

// Gather reads the stats from the HOBOlink API and writes it to the Accumulator.
func (s *Sensor) Gather(acc telegraf.Accumulator) error {
	if s.c == nil {
		c, err := s.createHTTPClient()
		if err != nil {
			return err
		}

		s.c = c
	}

	sensors, err := s.SensorList()
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG -- sensors yo! %+v\n", sensors)

	return nil
}

// SensorList ...
func (s *Sensor) SensorList() (*SensorList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sensorlist/%s", s.Server, s.Token), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sensorList SensorList
	if err := json.NewDecoder(resp.Body).Decode(&sensorList); err != nil {
		return nil, err
	}

	return &sensorList, nil
}

// NetworkList ...
func (s *Sensor) NetworkList() (*NetworkList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/networklist/%s", s.Server, s.Token), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var networkList NetworkList
	if err := json.NewDecoder(resp.Body).Decode(&networkList); err != nil {
		return nil, err
	}
	return &networkList, nil
}

// createHTTPClient is a helper to generate the HTTP Client connection for our endpoint
// TODO: check if TLS support available
func (s *Sensor) createHTTPClient() (*http.Client, error) {
	tr := &http.Transport{
		ResponseHeaderTimeout: s.HttpTimeout.Duration,
		TLSClientConfig:       &tls.Config{},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   s.HttpTimeout.Duration,
	}

	return client, nil
}

// Sensor is the default struct for the iMonnit API
type Sensor struct {
	c      *http.Client
	Server string
	Token  string

	HttpTimeout internal.Duration
}

// NetworkList is used to retrieve and parse the list of networks
type NetworkList struct {
	Method string    `json:"Method"`
	Result []Network `json:"Result"`
}

// Network is use to outline the details of a particular network
type Network struct {
	NetworkID           int    `json:"NetworkID"`
	NetworkName         string `json:"NetworkName"`
	SendNotification    bool   `json:"SendNotifications"`
	ExternalAccessUntil string `json:"ExternalAccessUntil"`
}

// SensorList is used to retrieve and parse the list of sensors
type SensorList struct {
	Method string   `json:"Method"`
	Result []Sensor `json:"Result"`
}

// SensorDetail is used to outline the details of a given sensor and its measurements
type SensorDetail struct {
	SensorID                    int    `json:"SensorID"`
	ApplicationID               int    `json:"ApplicationID"`
	CSNetID                     int    `json:"CSNetID"`
	SensorName                  string `json:"SensorName"`
	LastCommunicationDate       string `json:"LastCommunicateDate"`
	NextCommunicationDate       string `json:"NextCommunicateDate"`
	LastDataMessageMessageGUID  string `json:"LastDataMessageMessageGUID"`
	PowerSourceID               int    `json:"PowerSourceID"`
	Status                      int    `json:"Status"`
	CanUpdate                   bool   `json:"CanUpdate"`
	CurrentReadng               string `json:"CurrentReading"`
	BatteryLevel                int    `json:"BatteryLevel"`
	SignalStrength              int    `json:"SignalStrength"`
	AlertsActive                bool   `json:"AlertsActvie"`
	CheckDigit                  string `json:"CheckDigit"`
	AccountID                   int    `json:"AccountID"`
	MonnitApplicationID         int    `json:"MonnitApplicationID"`
	ReportInterval              int    `json:"ReportInterval"`
	ActiveStateInterval         int    `json:"ActiveStateInterval"`
	InactivityAlert             int    `json:"InactivityAlert"`
	MeasurementsPerTransmission int    `json:"MeasurementsPerTransmission"`
	MinimumThreshold            int    `json:"MinimumThreshold"`
	MaximumThreshold            int    `json:"MaximumThreshold"`
	Hysteresis                  int    `json:"Hysteresis"`
	Tag                         string `json:"Tag"`
}

const sampleConfig = `
  ## NOTE: this plugin assumes a one-hour window for query for each of the
  ##       serial numbers outlined in the configuration

  ## specify the URL of the iMonnit API (e.g. https://www.imonnit.com/json) 
  ## NOTE: this plugin _only_ supports the JSON API
  server = ""

  ## specify the token of the account to query from
  token = ""

  ## specify the list of serial numbers to query within the specified account
  ## NOTE: if this list is empty, this plugin will query all devices 
  ##       as found for the given user; if this list is non-empty, _only_ the
  ##       provided list of serial numbers will be queried
  serial_numbers = [""]

  ## Timeout for HTTP requests to the HOBOlink API URL 
  # http_timeout = "5s"
`
