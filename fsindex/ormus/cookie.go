package ormus

import (
	"net/http"
	"net/url"
)

func getCookie(cname string, client *http.Request) *http.Cookie {
	var result *http.Cookie
	if xid, e := client.Cookie(cname); e == nil {
		result = xid
	}
	return result
}
func getCookieValue(cname string, client *http.Request) string {
	cookie := getCookie(cname, client)
	cookieValue := ""
	if cookie != nil {
		if sessid, x := url.QueryUnescape(cookie.Value); x == nil {
			cookieValue = sessid
		}
	}
	return cookieValue
}
func cookieValue(cookie *http.Cookie) string {
	cookieValue := ""
	if cookie != nil {
		if sessid, x := url.QueryUnescape(cookie.Value); x == nil {
			cookieValue = sessid
		}
	}
	return cookieValue
}
