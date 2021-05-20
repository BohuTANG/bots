package main

import (
	"os"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
	log "github.com/sirupsen/logrus"
)

const (
	path = "/webhooks"
)

func main() {
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		log.Fatal("Missing input 'GITHUB_TOKEN' env.")
	}
	hook, _ := github.New(github.Options.Secret(token))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			log.Infof("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			log.Infof("%+v", pullRequest)
		}
	})
	http.ListenAndServe(":3000", nil)
}
