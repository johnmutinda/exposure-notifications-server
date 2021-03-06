// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwar
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authorizedapps

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/exposure-notifications-server/internal/admin"
	"github.com/google/exposure-notifications-server/internal/authorizedapp/database"
	"github.com/google/exposure-notifications-server/internal/authorizedapp/model"
	"github.com/google/exposure-notifications-server/internal/serverenv"
)

type viewController struct {
	config *admin.Config
	env    *serverenv.ServerEnv
}

func NewView(c *admin.Config, env *serverenv.ServerEnv) *viewController {
	return &viewController{config: c, env: env}
}

func (h *viewController) Execute(c *gin.Context) {
	ctx := c.Request.Context()
	m := admin.TemplateMap{}

	appID, _ := c.GetQuery("apn")
	authorizedApp := model.NewAuthorizedApp()

	if appID == "" {
		m.AddJumbotron("Authorized Applications", "Create New Authorized Application")
		m["new"] = true
	} else {
		aadb := database.New(h.env.Database())
		var err error
		authorizedApp, err = aadb.GetAuthorizedApp(ctx, h.env.SecretManager(), appID)
		if err != nil {
			admin.ErrorPage(c, err.Error())
			return
		}
		m.AddJumbotron("Authorized Applications", fmt.Sprintf("Edit: `%v`", authorizedApp.AppPackageName))
	}
	m["app"] = authorizedApp
	m["previousKey"] = base64.StdEncoding.EncodeToString([]byte(authorizedApp.AppPackageName))
	c.HTML(http.StatusOK, "authorizedapp", m)
}
