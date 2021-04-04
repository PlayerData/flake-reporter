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
			err = updateTestSummary(handler.client, handler.ctx, project, suite.Name, branch, test)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
