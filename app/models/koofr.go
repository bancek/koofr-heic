package models

import (
	"context"

	koofrclient "github.com/koofr/go-koofrclient"
	"github.com/revel/revel"
	"golang.org/x/oauth2"
)

var KoofrBaseUrl = "https://app.koofr.net"
var KoofrOAuthConfig *oauth2.Config

type updateTokenSource struct {
	tokenSource   oauth2.TokenSource
	onUpdateToken func(t *oauth2.Token)
}

func (s *updateTokenSource) Token() (*oauth2.Token, error) {
	t, err := s.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	s.onUpdateToken(t)
	return t, nil
}

func InitKoofrOAuthConfig() {
	clientId := revel.Config.StringDefault("koofr.client_id", "")
	clientSecret := revel.Config.StringDefault("koofr.client_secret", "")
	redirectUrl := revel.Config.StringDefault("koofr.redirect_url", "")

	if clientId == "" || clientSecret == "" || redirectUrl == "" {
		panic("Missing koofr.client_id, koofr.client_secret, koofr.redirect_url")
	}

	KoofrOAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"public"},
		RedirectURL:  redirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  KoofrBaseUrl + "/oauth2/auth",
			TokenURL: KoofrBaseUrl + "/oauth2/token",
		},
	}
}

func GetKoofrClient(token *oauth2.Token, onUpdateToken func(t *oauth2.Token)) *koofrclient.KoofrClient {
	ctx := context.Background()
	tokenSource := KoofrOAuthConfig.TokenSource(ctx, token)
	tokenSource = &updateTokenSource{
		tokenSource:   tokenSource,
		onUpdateToken: onUpdateToken,
	}
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client := koofrclient.NewKoofrClient(KoofrBaseUrl, false)
	client.HTTPClient.Client = httpClient

	return client
}

// func KoofrUpload(koofr *koofrclient.KoofrClient, filePath string, name string) (shortUrl string, err error) {
// 	mounts, err := koofr.Mounts()
// 	if err != nil {
// 		return "", err
// 	}

// 	primaryMountId := ""

// 	for _, mount := range mounts {
// 		if mount.IsPrimary {
// 			primaryMountId = mount.Id
// 			break
// 		}
// 	}

// 	if primaryMountId == "" {
// 		return "", fmt.Errorf("Primary mount id not found")
// 	}

// 	reader, err := os.Open(filePath)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer reader.Close()

// 	newName, err := koofr.FilesPut(primaryMountId, "/", "YouTube to Koofr/"+name, reader)
// 	if err != nil {
// 		return "", err
// 	}

// 	remotePath := "/YouTube to Koofr/" + newName

// 	shortUrl, err = KoofrCreateShortUrl(koofr, primaryMountId, remotePath)
// 	if err != nil {
// 		return "", err
// 	}

// 	return shortUrl, nil
// }
