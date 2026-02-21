package bundler

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

// TransformJS is a standalone utility for one-off JSX -> JS transforms.
// Used by the watcher for incremental single-file rebuilds in dev mode
// instead of full rebuild on every keystroke.
func TransformJSX(source string, filename string) (string, error) {
	result := api.Transform(source, api.TransformOptions{
		Loader:     api.LoaderTSX, // Handles both .tsx and .jsx
		JSX:        api.JSXAutomatic,
		Target:     api.ES2020,
		Format:     api.FormatESModule,
		Sourcefile: filename, // Shows correct filename in errors
	})

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("transform %s: %s", filename, result.Errors[0].Text)
	}

	return string(result.Code), nil
}
