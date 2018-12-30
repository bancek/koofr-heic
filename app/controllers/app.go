package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/bancek/koofr-heic/app/models/koofrheictojpg"

	koofrclient "github.com/koofr/go-koofrclient"
	"github.com/revel/revel"
	"golang.org/x/oauth2"

	"github.com/bancek/koofr-heic/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	me := map[string]interface{}{}

	if koofr, ok := c.koofr(); ok {
		info, err := koofr.UserInfo()

		if err != nil {
			revel.ERROR.Println(err)
		} else {
			me["name"] = info.FirstName + " " + info.LastName
		}
	}

	authUrl := models.KoofrOAuthConfig.AuthCodeURL("")

	return c.Render(me, authUrl)
}

func (c App) Auth(code string) revel.Result {
	token, err := models.KoofrOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Index)
	}

	user := c.user()
	user.OAuth2Token = token
	return c.Redirect(App.Index)
}

type ConvertResult struct {
	mountId         string
	path            string
	convertMovToMp4 bool
	koofr           *koofrclient.KoofrClient
}

func (r *ConvertResult) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")

	writer := resp.GetWriter()
	flusher := writer.(http.Flusher)

	for i := 0; i < 4096; i++ {
		writer.Write([]byte(" "))
	}
	flusher.Flush()

	writer.Write([]byte("<pre><code>\n"))

	logger := func(line string) {
		writer.Write([]byte(line))
		for i := 0; i < 1024; i++ {
			writer.Write([]byte(" "))
		}
		writer.Write([]byte("\n"))
		flusher.Flush()

		if revel.DevMode {
			fmt.Println(line)
		}

		flusher.Flush()
	}

	err := koofrheictojpg.Convert(r.koofr, r.mountId, r.path, r.convertMovToMp4, logger)

	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	writer.Write([]byte("</pre></code>\n"))

	return
}

func (c App) Convert(mountId string, path string, webUrl string, convertMov string) revel.Result {
	if webUrl != "" {
		u, err := url.Parse(webUrl)
		if err == nil {
			result := regexp.MustCompile("^/app/storage/([^/]+)$").FindAllStringSubmatch(u.Path, -1)
			if len(result) == 1 && len(result[0]) == 2 {
				mountId = result[0][1]
				path = u.Query().Get("path")
			}
		}
	}
	if mountId == "" || path == "" {
		return c.RenderText("Missing mountId or path")
	}
	convertMovToMp4 := convertMov == "on"
	if koofr, ok := c.koofr(); ok {
		return &ConvertResult{
			mountId:         mountId,
			path:            path,
			koofr:           koofr,
			convertMovToMp4: convertMovToMp4,
		}
	} else {
		return c.Redirect(App.Index)
	}
}

func (c App) user() *models.User {
	return c.Args["user"].(*models.User)
}

func (c App) koofr() (*koofrclient.KoofrClient, bool) {
	user := c.user()

	if user.OAuth2Token == nil {
		return nil, false
	}

	onUpdateToken := func(t *oauth2.Token) {
		user.OAuth2Token = t
	}

	koofr := models.GetKoofrClient(user.OAuth2Token, onUpdateToken)

	return koofr, true
}

func setuser(c *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := c.Session["uid"]; ok {
		user = models.GetUser(c.Session["uid"])
	}
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = user.Id
	}
	c.Args["user"] = user
	return nil
}
