// +build !lite

package middleman

import (
	"net/http"

	"github.com/rs/cors"
)

// Cross Origin Resource Sharing (https://www.w3.org/TR/cors/)
func CORS(hosts []string, heir http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: hosts,
	})
	return c.Handler(heir)
}
