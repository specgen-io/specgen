package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/java/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func generateThrowClientException(w *writer.Writer, errorMessage string, wrapException string) {
	w.Line(`var errorMessage = %s;`, errorMessage)
	w.Line(`logger.error(errorMessage);`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw new ClientException(%s);`, params)
}
