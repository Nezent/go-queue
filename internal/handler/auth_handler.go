package handler

import (
	"net/http"
	"time"

	"github.com/Nezent/go-queue/common"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	expire := time.Unix(0, 0).UTC()
	cookies := []string{"access_token", "refresh_token"}
	for _, name := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  expire,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			// Domain: "example.com", // optional
		})
	}
	common.RespondJSON(w, http.StatusOK, common.SuccessResponse("Logged out successfully", nil))
}
