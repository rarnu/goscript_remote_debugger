package util

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
)

func GenerateSourceMap(filename string, src string) []byte {
	result := api.Transform(src, api.TransformOptions{
		Sourcemap:         api.SourceMapInline,
		SourcesContent:    api.SourcesContentInclude,
		Sourcefile:        filename,
		MinifyWhitespace:  false,
		MinifyIdentifiers: false,
		MinifySyntax:      false,
	})

	if len(result.Errors) > 0 {
		fmt.Println(result.Errors)
	}

	return result.Code
}
