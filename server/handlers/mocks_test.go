package handlers_test

import (
	"github.com/Thiht/smocker/server"
	"github.com/Thiht/smocker/server/config"
	"github.com/Thiht/smocker/server/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestMocks_GenericHandler(t *testing.T) {
	server, mocks := server.NewMockServer(config.Config{
		LogLevel:             "panic",
		ConfigListenPort:     8081,
		ConfigBasePath:       "/",
		MockServerListenPort: 8080,
		StaticFiles:          ".",
		Build: config.Build{
			AppName:      "smocker",
			BuildVersion: "dev",
			BuildDate:    time.Now().String(),
		},
	})
	session := mocks.NewSession("test")
	_, err := mocks.AddMock(session.ID, &types.Mock{
		Request: types.MockRequest{
			Method: types.StringMatcher{Matcher: "ShouldMatch", Value: "GET"},
			Path:   types.StringMatcher{Matcher: "ShouldMatch", Value: "/api/v1"},
		},
		Response: &types.MockResponse{
			Status: 200,
			Body:   "test",
		},
		Context: &types.MockContext{
			Times: 200,
		},
	})
	require.NoError(t, err)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	for i := 0; i < 1; i++ {
		resp, err := http.Get("http://localhost:8080/api/v1")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "test", string(body))
	}
}
