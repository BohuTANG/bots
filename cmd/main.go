package main

import (
	"encoding/json"
	"os"
	"strings"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
	"github.com/jimschubert/labeler"
	log "github.com/sirupsen/logrus"
)

const (
	path = "/webhooks"
)

/*
* export GITHUB_SECRET='don't tell'
* export GITHUB_TOKEN='don't tell'
* export GITHUB_REPO='datafuselabs/datafuse'
*/
func main() {
	// Check GITHUB_SECRET.
	secret, found := os.LookupEnv("GITHUB_SECRET")
	if !found {
		log.Fatal("Missing input 'GITHUB_SECRET' env.")
	}

	// Check GITHUB_TOKEN.
	_, found = os.LookupEnv("GITHUB_TOKEN")
	if !found {
		log.Fatal("Missing input 'GITHUB_TOKEN' env.")
	}

	hook, _ := github.New(github.Options.Secret(secret))

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
			pr := payload.(github.PullRequestPayload)
			log.Infof("Pull reqeust id: %+v", pr.Number)
			repo := os.Getenv("GITHUB_REPO")
			repoParts := strings.Split(repo, "/")
			owner := repoParts[0]
			repoName := repoParts[1]
			body, _:= json.Marshal(pr)
			data := string(body)
			l, err := labeler.New(owner, repoName, "pull_request", int(pr.Number), &data)
			if err != nil {
				log.Errorf("Could not construct a labeling: %v", err)
			}
			if l.Execute()!= nil {
				log.Errorf("Failed to execute labeling: %v", err)
			}

			log.Info("Done labeling.")
		}
	})
	http.ListenAndServe(":3000", nil)
}
