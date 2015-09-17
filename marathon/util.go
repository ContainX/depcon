package marathon

import (
	"github.com/gondor/depcon/pkg/envsubst"
	"io"
	"os"
)

func substFileTokens(in io.Reader, filename string, params map[string]string) (parsed string, missing bool) {
	parsed = envsubst.Substitute(in, true, func(s string) string {
		if params != nil && params[s] != "" {
			return params[s]
		}
		if os.Getenv(s) == "" {
			log.Warning("Cannot find a value for varable ${%s} which was defined in %s", s, filename)
			missing = true
		}
		return os.Getenv(s)
	})
	return parsed, missing
}
