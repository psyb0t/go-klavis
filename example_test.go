package klavis_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-klavis"
)

// This is the README example. It compiles without calling the live API.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := klavis.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization",
				"Bearer "+os.Getenv("KLAVIS_API_KEY"))

			return nil
		},
	)

	client, err := klavis.NewClientWithResponses(
		"https://api.klavis.ai",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAllMcpServersWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
