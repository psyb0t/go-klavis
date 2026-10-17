package klavis_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-klavis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey        = "EXAMPLE-DO-NOT-USE"
	testAuthorization = "Bearer " + testAPIKey
	authHeader        = "Authorization"
	jsonHeader        = "Content-Type"
	jsonType          = "application/json"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *klavis.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := klavis.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(authHeader, testAuthorization)

			return nil
		},
	)

	client, err := klavis.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func TestGetAllMcpServers(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mcp-server/servers", r.URL.Path)
		assert.Equal(t, testAuthorization, r.Header.Get(authHeader))

		w.Header().Set(jsonHeader, jsonType)
		_, err := w.Write([]byte(`{"servers":[{
			"id":"00000000-0000-4000-8000-000000000001",
			"name":"Example","description":null,"tools":null
		}]}`))
		assert.NoError(t, err)
	})

	res, err := client.GetAllMcpServersWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Servers, 1)
	assert.Equal(t, "Example", res.JSON200.Servers[0].Name)
	assert.Nil(t, res.JSON200.Servers[0].Description)
	assert.Nil(t, res.JSON200.Servers[0].Tools)
}

func TestCreateServerInstance(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/mcp-server/instance/create", r.URL.Path)
		assert.Equal(t, testAuthorization, r.Header.Get(authHeader))
		assert.Equal(t, jsonType, r.Header.Get(jsonHeader))

		var body klavis.CreateServerRequest

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, klavis.McpServerName("Notion"), body.ServerName)
		assert.Equal(t, "user-example", body.UserId)
		assert.Nil(t, body.PlatformName)

		w.Header().Set(jsonHeader, jsonType)
		_, err := w.Write([]byte(`{
			"serverUrl":"https://mcp.example.com/instance-example",
			"instanceId":"instance-example","oauthUrl":null
		}`))
		assert.NoError(t, err)
	})

	res, err := client.CreateServerInstanceWithResponse(
		context.Background(),
		klavis.CreateServerRequest{
			ServerName: klavis.McpServerName("Notion"),
			UserId:     "user-example",
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	assert.Equal(t, "instance-example", res.JSON200.InstanceId)
	assert.Equal(t,
		"https://mcp.example.com/instance-example", res.JSON200.ServerUrl)
	assert.Nil(t, res.JSON200.OauthUrl)
}

func TestGetAllMcpServersErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		status int
		body   string
	}{
		{
			name:   "unauthorized",
			status: http.StatusUnauthorized,
			body:   `{"detail":"Invalid API key"}`,
		},
		{
			name:   "rate limited",
			status: http.StatusTooManyRequests,
			body:   `{"detail":"Too many requests"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, testAuthorization,
						r.Header.Get(authHeader))
					w.Header().Set(jsonHeader, jsonType)
					w.WriteHeader(tc.status)
					_, err := w.Write([]byte(tc.body))
					assert.NoError(t, err)
				})

			res, err := client.GetAllMcpServersWithResponse(
				context.Background(),
			)
			require.NoError(t, err)
			assert.Equal(t, tc.status, res.StatusCode())
			assert.JSONEq(t, tc.body, string(res.Body))
			assert.Nil(t, res.JSON200)
		})
	}
}

func TestGetAllMcpServersMalformedJSON(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(jsonHeader, jsonType)
		_, err := w.Write([]byte(`{"servers":`))
		assert.NoError(t, err)
	})

	_, err := client.GetAllMcpServersWithResponse(context.Background())

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
}
