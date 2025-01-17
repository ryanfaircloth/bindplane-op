// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/observiq/bindplane-op/internal/server"
	"github.com/observiq/bindplane-op/internal/server/sessions"
)

// CheckSession checks to see if the attached cookie session is authenticated
// and if so sets authenticated to true on the context.  If not authenticated it
// goes to the next handler.
func CheckSession(server server.BindPlane) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := server.Store().UserSessions().Get(c.Request, sessions.CookieName)
		if err != nil {
			// Clear the cookie, this can happen when sessions-secrets change
			// and we see a cookie with the previous secret is read.
			session.Options.MaxAge = 0

			saveErr := session.Save(c.Request, c.Writer)
			if saveErr != nil {
				c.AbortWithError(http.StatusInternalServerError, saveErr)
				return
			}

			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Check the authenticated value in the session storage - if unset or false go to next handler
		if session.Values["authenticated"] == nil || session.Values["authenticated"] == false {
			c.Next()
			return
		}

		c.Keys["authenticated"] = true
		// Extend the cookies life by 15 minutes since the user is active and making requests.
		session.Options.MaxAge = 15 * 60
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
