package assist

import (
	"encoding/json"
	"golang.org/x/text/language"
	"net"
	"net/http"
	"net/url"
)

const (
	KeyRequestId  = "rid"
	KeyCtx        = "ctx"
)

func HttpRequestBodyResolve(w http.ResponseWriter, r *http.Request, i interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		SendJsonError(w, http.StatusBadRequest, "unable to parse request body: "+err.Error())
		return false
	}
	return true
}

func HttpRequestLanguage(r *http.Request) *language.Tag {
	header := r.Header.Get("Accept-Language")
	tags, _, _ := language.ParseAcceptLanguage(header)
	if len(tags) > 0 {
		return &tags[0]
	}
	return nil
}

func HttpRequestIP(r *http.Request) net.IP {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return net.ParseIP(ip)
}

func IsUrlOk(url *url.URL) bool {
	return url != nil && url.IsAbs()
}
