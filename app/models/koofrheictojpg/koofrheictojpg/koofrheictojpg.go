package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bancek/koofr-heic/app/models/koofrheictojpg"
	koofrclient "github.com/koofr/go-koofrclient"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var baseUrl string
	var email string
	var password string
	var mountId string
	var path string
	var convertMovToMp4 bool

	flag.StringVar(&baseUrl, "baseUrl", "https://app.koofr.net", "Koofr base URL")
	flag.StringVar(&email, "email", "", "Koofr email")
	flag.StringVar(&password, "password", "", "Koofr password")
	flag.StringVar(&mountId, "mountId", "", "Koofr mount id")
	flag.StringVar(&path, "path", "", "Koofr path")
	flag.BoolVar(&convertMovToMp4, "convertMovToMp4", false, "Convert MOV to MP4")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	if email == "" || password == "" || mountId == "" || path == "" {
		flag.Usage()
	}

	koofr := koofrclient.NewKoofrClient(baseUrl, false)
	koofr.HTTPClient.Headers.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(email+":"+password)))

	logger := func(s string) {
		log.Println(s)
	}

	err := koofrheictojpg.Convert(koofr, mountId, path, convertMovToMp4, logger)
	if err != nil {
		log.Fatal(err)
	}
}
