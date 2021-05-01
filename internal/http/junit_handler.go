package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/joshdk/go-junit"

	"playerdata.co.uk/flake-reporter/internal/models"
)

type JUnitHandler struct {
	Client *firestore.Client
	Ctx    context.Context
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

	fmt.Fprintf(w, "%s\n", "Ingesting test results:")

	for _, suite := range suites {
		for _, test := range suite.Tests {
			err = models.UpdateBranchSummary(handler.Client, handler.Ctx, project, suite.Name, branch, test)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s", err)
				return
			}
			fmt.Fprintf(w, "%s - %s\n", suite.Name, test.Name)
		}
	}

	w.WriteHeader(http.StatusOK)
}
