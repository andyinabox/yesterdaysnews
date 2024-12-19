package server

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/charmbracelet/log"
)

type JSVars struct {
	ObjectStoreUrl string `json:"objectStoreUrl"`
}

func (j *JSVars) ToJS() template.JS {
	b, err := json.MarshalIndent(j, "", "  ")

	if err != nil {
		log.Errorf("cannot marshal JSVars: %v", j)
	}

	return template.JS(fmt.Sprintf("window.YN_ENV = %s", string(b)))
}
