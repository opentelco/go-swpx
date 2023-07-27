package stanza

import (
	"bytes"
	"html/template"
)

func FromTemplate(stanzaId string, content string, data interface{}) (string, error) {
	var err error

	tmp := template.New(stanzaId)
	temp, err := template.Must(tmp, err).Parse(content)
	if err != nil {
		return "", err
	}

	// store template execution in buffer

	var buf bytes.Buffer
	temp.Execute(&buf, data)
	return buf.String(), nil
}
