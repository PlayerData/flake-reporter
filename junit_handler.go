package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/joshdk/go-junit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JUnitHandler struct {
	client *firestore.Client
	ctx    context.Context
}

func updateTestStats(client *firestore.Client, ctx context.Context, project string, branch string, suite junit.Suite, test junit.Test) (err error) {
	path := fmt.Sprintf("projects/%s/branches/%s/suites/%s/tests/%s", project, branch, suite.Name, test.Name)
	doc := client.Doc(path)

	stats := TestStats{SuccessCount: 0, FailureCount: 0}

	docRef, err := doc.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}

	if docRef.Exists() {
		err = docRef.DataTo(&stats)
		if err != nil {
			return err
		}
	}

	if test.Status == junit.StatusPassed {
		stats.SuccessCount += 1
	} else {
		stats.FailureCount += 1
	}

	_, err = doc.Set(ctx, stats)

	return err
}

func populateFormValues(r *http.Request, project *string, branch *string, junitXml *bytes.Buffer) (err error) {
	r.ParseMultipartForm(32 << 20)

	*project = r.FormValue("project")
	*branch = r.FormValue("branch")

	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}

	defer file.Close()

	io.Copy(junitXml, file)

	return nil
}

func (handler *JUnitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var project string
	var branch string
	var junitXml bytes.Buffer
	err := populateFormValues(r, &project, &branch, &junitXml)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "%s", err)
	}

	suites, err := junit.Ingest(junitXml.Bytes())
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
