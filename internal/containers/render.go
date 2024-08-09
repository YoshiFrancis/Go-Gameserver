package containers

import (
	"bytes"
	"html/template"
)

func RenderTemplate[T any](template *template.Template, dataStruct T) []byte {
	var templateBuffer bytes.Buffer
	template.Execute(&templateBuffer, dataStruct)
	return templateBuffer.Bytes()
}
