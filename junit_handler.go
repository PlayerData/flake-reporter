package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/joshdk/go-junit"
)

type JUnitHandler struct {
	client *firestore.Client
	ctx    context.Context
}

func updateTestStats(client *firestore.Client, ctx context.Context, project string, branch string, suite junit.Suite, test junit.Test) (err error) {
	path := fmt.Sprintf("projects/%s/branches/%s/suites/%s/tests/%s", project, branch, suite.Name, test.Name)
	doc := client.Doc(path)

	stats := TestStats{SuccessCount: 0, FailureCount: 0}

	if test.Status == junit.StatusPassed {
		stats.SuccessCount += 1
	} else {
		stats.FailureCount += 1
	}

	_, err = doc.Set(ctx, stats)

	return err
}

func (handler *JUnitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var xml bytes.Buffer

	project := r.FormValue("project")
	branch := r.FormValue("branch")

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer file.Close()

	io.Copy(&xml, file)

	suites, err := junit.Ingest(xml.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%s", err)
		return
	}

	for _, suite := range suites {
		for _, test := range suite.Tests {
			err = updateTestStats(handler.client, handler.ctx, project, branch, suite, test)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
