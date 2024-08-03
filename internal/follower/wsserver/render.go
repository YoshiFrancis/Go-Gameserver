package wsserver

import (
	"bytes"
	"html/template"
)

func renderTemplate[T any](template *template.Template, dataStruct T) []byte {
	var templateBuffer bytes.Buffer
	template.Execute(&templateBuffer, dataStruct)
	return templateBuffer.Bytes()
}
