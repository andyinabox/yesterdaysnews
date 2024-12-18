package youtubedownloader

import (
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/shellargs"
)

type Request struct {
	// Format
	Format string `args:"--format"`

	// Subtitles
	WriteAutoSubs bool   `args:"--write-auto-subs"`
	SubFormat     string `args:"--sub-format"`

	// Verbosity & simulation
	Quiet          bool `args:"--quiet"`
	NoSimulate     bool `args:"--no-simulate"`
	DumpJSON       bool `args:"--dump-json"`
	DumpSingleJSON bool `args:"--dump-single-json"`

	// File output
	Output string `args:"--output"`
}

func (r *Request) Args() *shellargs.Args {
	a := shellargs.New()
	t := reflect.TypeOf(*r)
	v := reflect.ValueOf(*r)

	// iterate throug struct fields
	for i := 0; i < t.NumField(); i++ {

		// get field and tag
		field := t.Field(i)
		fieldKind := field.Type.Kind()
		tag := field.Tag.Get("args")

		// break up tag contents
		tagParts := strings.Split(tag, ",")
		if len(tagParts) == 0 {
			log.Errorf("missing tag values: %q", tag)
		}

		argName := tagParts[0]

		// log.Debug("parse field", "field", field, "fieldKind", fieldKind, "tag", tag, "argName", argName)

		// build shellargs based on kind and tag
		switch fieldKind {
		case reflect.Bool:
			val := v.FieldByName(field.Name).Bool()
			if val {
				a.Add(argName)
			}
			continue
		case reflect.String:
			val := v.FieldByName(field.Name).String()
			if val != "" {
				a.AddKeyedSingleQuoted(argName, val)
			}
			continue
		default:
			log.Errorf("invalid type: %s", fieldKind)
		}
	}

	return a
}

func (r *Request) String() string {
	return r.Args().String()
}
