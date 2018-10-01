package controllers

import (
	"fmt"
	"net/http"

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
	mountId string
	path    string
	koofr   *koofrclient.KoofrClient
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

	err := koofrheictojpg.Convert(r.koofr, r.mountId, r.path, logger)

	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	writer.Write([]byte("</pre></code>\n"))

	return
}

func (c App) Convert(mountId string, path string) revel.Result {
	if koofr, ok := c.koofr(); ok {
		return &ConvertResult{
			mountId: mountId,
			path:    path,
			koofr:   koofr,
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
