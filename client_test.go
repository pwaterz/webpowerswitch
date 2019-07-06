package webpowerswitch

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const outletsResponse = `[
	{
			"name": "HLG651",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "WaterPump",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "HLG652",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "Aerator",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "ElectricSky180",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "Fan1",
			"critical": true,
			"transient_state": true,
			"cycle_delay": null,
			"physical_state": true,
			"locked": false,
			"state": true
	},
	{
			"name": "Outlet 7",
			"critical": false,
			"transient_state": false,
			"cycle_delay": null,
			"physical_state": false,
			"locked": false,
			"state": false
	},
	{
			"name": "Outlet 8",
			"critical": false,
			"transient_state": false,
			"cycle_delay": null,
			"physical_state": false,
			"locked": false,
			"state": false
	}
]`

func TestHTTPGetOutlets(t *testing.T) {
	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(outletsResponse))
	})

	httpClient, tearDown := testingHTTPClient(h)
	defer tearDown()
	client, err := NewClient("http://www.test.com")
	client.httpClient = httpClient
	assert.Nil(t, err)
	outlets, err := client.GetOutlets()
	assert.Nil(t, err)
	assert.Equal(t, len(outlets), 8)
	assert.Equal(t, outlets[0].Critical, true)
	assert.Equal(t, outlets[0].Name, "HLG651")
	assert.Equal(t, outlets[0].TransientState, true)
	assert.Nil(t, outlets[0].CycleDelay)
	assert.Equal(t, outlets[0].PhysicalState, true)
	assert.Equal(t, outlets[0].Locked, false)
	assert.Equal(t, outlets[0].State, true)

}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}
