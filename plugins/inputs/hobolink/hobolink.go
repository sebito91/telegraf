package hobolink

import (
	"fmt"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// dataPath is the default enpoint for metrics collection
const dataPath = "/data/json"

// Description returns the plugin description.
func (h *HOBOlink) Description() string {
	return "Read stats from the HOBOlink API for a given user account"
}

// SampleConfig returns sample configuration for this plugin.
func (h *HOBOlink) SampleConfig() string {
	return sampleConfig
}

// HOBOlink is a plugin to read stats from a list of IoT RX Sensors via the HOBOlink API
// NOTE: this agent is meant to collect from a single source user account.
type HOBOlink struct {
	User          string
	Password      string
	Server        string
	SerialNumbers []string
	HTTPTimeout   internal.Duration

	client *http.Client
}

// NewHOBOlink returns a new instance of the HOBOlink collector
func NewHOBOlink() *HOBOlink {
	return &HOBOlink{
		HTTPTimeout: internal.Duration{Duration: time.Second * 5},
		Server:      fmt.Sprintf("https://webservice.hobolink.com/restv2%s", dataPath),
	}
}

func init() {
	inputs.Add("hobolink", func() telegraf.Input {
		return NewHOBOlink()
	})
}

// Gather reads the stats from the HOBOlink API and writes it to the Accumulator.
func (h *HOBOlink) Gather(acc telegraf.Accumulator) error {
	if h.client == nil {
		client, err := h.createHTTPClient()
		if err != nil {
			return err
		}

		h.client = client
	}

	return nil
}

// createHTTPClient is a helper to generate the HTTP Client connection for our endpoint
// TODO: check if TLS support available
func (h *HOBOlink) createHTTPClient() (*http.Client, error) {
	tr := &http.Transport{
		ResponseHeaderTimeout: h.HTTPTimeout.Duration,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   h.HTTPTimeout.Duration,
	}

	return client, nil
}

const sampleConfig = `
  ## specify the URL of the HOBOlink API 
  server = "https://webservice.hobolink.com/restv2"

  ## specify the username of the account to query from
  user = ""

  ## specify the password of the account to query from
  password = ""

  ## specify the list of serial numbers to query within the specified account
  ## NOTE: if this list is empty, this plugin will query all devices 
  ##       as found for the given user; if this list is non-empty, _only_ the
  ##       provided list of serial numbers will be queried
  serial_numbers = [""]

  ## Timeout for HTTP requests to the HOBOlink API URL 
  http_timeout = "5s"
`
