package openapihelper

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
}

func UrlDecodeParam(obj interface{}, r *http.Request) error {
	return decoder.Decode(obj, r.URL.Query())
}
