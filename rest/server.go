package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/pmh-only/spoti2wall/utils"
)

const serverPort = 49823

func redirect2Auth(w http.ResponseWriter, r *http.Request) {
	authUrl := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "authorize",
		RawQuery: url.Values{
			"response_type": {"code"},
			"client_id":     {},
			"scope":         {"user-read-playback-state"},
			"redirect_uri":  {redirectUri},
		}.Encode(),
	}

	http.Redirect(w, r, authUrl.String(), http.StatusFound)
}

func callback(server *http.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "You can now close this window.\n")

		code := r.URL.Query().Get("code")
		AccessToken, RefreshToken = getAuthToken(code)

		utils.SaveRefreshToken(RefreshToken)

		color.Cyan("🔑 Login successful.\n")

		go server.Close()
		go KeepRefreshToken()
	}
}

func StartAuthServer() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: mux,
	}

	mux.HandleFunc("/authorize", redirect2Auth)
	mux.HandleFunc("/callback", callback(server))

	go func() {
		color.Green("🔑 Please open this url in your browser:")
		color.Green("http://localhost:%d/authorize", serverPort)

		utils.OpenBrowser(
			fmt.Sprintf("http://localhost:%d/authorize", serverPort))
	}()

	server.ListenAndServe()
}
