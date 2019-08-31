package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/alphagov/iap/internal"
	"github.com/sirupsen/logrus"
)

type googleUserData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified_email"`
	AvatarURL string `json:"picture"`
	Domain    string `json:"hd"`
}

func authenticatedMiddleware(ctx internal.Context, hanlder http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uniqueCode := fmt.Sprintf("iap-%d", time.Now().Unix())
		data, err := obtainUserDataFromGoogle(r)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Debug("unable to userData")
			http.Redirect(w, r, routeDefault, http.StatusTemporaryRedirect)
			return
		}
		ctx.Logger.WithField("email", data.Email).Debug("authenticated user")

		hanlder(w, r)
	}
}

func obtainUserDataFromGoogle(r *http.Request) (googleUserData, error) {
	// TODO: instead of DDoS'ing google, we should consider to store the data somewhere and verify the accessToken...
	cookie, err := r.Cookie(tokenCookieName)
	fmt.Println(r.Cookies())
	if err != nil {
		return googleUserData{}, fmt.Errorf("unable to obtain accessToken: %s", err)
	}
	response, err := http.Get(oauthGoogleURLAPI + cookie.Value)
	if err != nil {
		return googleUserData{}, fmt.Errorf("failed getting user info: %s", err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return googleUserData{}, fmt.Errorf("failed read response: %s", err)
	}
	var data googleUserData
	return data, json.Unmarshal(contents, &data)
}
