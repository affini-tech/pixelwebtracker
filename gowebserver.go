package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

const cookieIdentifier = "trackbox000"

var oneGif = []byte("GIF87a\x01\x00\x01\x00\x80\x00\x00\xff\xff\xff\xff\xff\xff,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")

func getRequestData(r *http.Request, cid string) string {
	referer := r.Referer()
	ua := r.UserAgent()
	addr := r.RemoteAddr
	uri := r.RequestURI
	url := r.URL

	// Poor (& lazy) man JSON marshaling
	outputJSON := fmt.Sprintf("{\"CookieIdentifier\" : \"%s\"", cid)
	outputJSON = fmt.Sprintf("%s,\"DateTime\" : \"%s\"", outputJSON, time.Now().String())
	outputJSON = fmt.Sprintf("%s,\"Referer\" : \"%s\"", outputJSON, referer)
	outputJSON = fmt.Sprintf("%s,\"UserAgent\" : \"%s\"", outputJSON, ua)
	outputJSON = fmt.Sprintf("%s,\"Host\" : \"%s\"", outputJSON, url.Host)
	outputJSON = fmt.Sprintf("%s,\"EscapedPath\" : \"%s\"", outputJSON, url.EscapedPath())
	outputJSON = fmt.Sprintf("%s,\"Url\" : \"%s\"", outputJSON, url.String())
	outputJSON = fmt.Sprintf("%s,\"RemoteAddr\" : \"%s\"", outputJSON, addr)
	outputJSON = fmt.Sprintf("%s,\"RequestURI\" : \"%s\"", outputJSON, uri)
	outputJSON = fmt.Sprintf("%s,\"Query\" : \"%s\"", outputJSON, url.RawQuery)
	outputJSON = fmt.Sprintf("%s}", outputJSON)
	return outputJSON
}

func testReqDataHandler(w http.ResponseWriter, r *http.Request) {
	output := getRequestData(r, "unknown")
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(output))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("hello <a href='/testreqdata?abc=fgdgfd'>reqdata</a> <img src=\"/pixel.gif\"/>"))
}

func randStr(strSize int) string {
	var dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func buildCookie(r *http.Request) (string, *http.Cookie) {
	c, _ := r.Cookie(cookieIdentifier)
	out := randStr(32)
	if c != nil {
		out = c.Value
		//fmt.Printf("[track] Actual cookie : %s\n", out)
	}
	//fmt.Printf("[track] Gen cookie : %s\n", out)
	theCookie := &http.Cookie{
		Name:    cookieIdentifier,
		Expires: time.Now().AddDate(1, 0, 0),
		Value:   out,
	}
	// Poor (& lazy) man JSON marshaling
	output := fmt.Sprintf("{\"trackId\" : \"%s\"}", out)
	return output, theCookie
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	output, myCookie := buildCookie(r)
	data := getRequestData(r, myCookie.Value)
	//fmt.Printf("[track] Request Data : %s\n", data)
	fmt.Printf("%s\n", data)
	w.Header().Add("Content-Type", "application/json")
	http.SetCookie(w, myCookie)
	w.Write([]byte(output))
}

func pixelHandler(w http.ResponseWriter, r *http.Request) {
	_, myCookie := buildCookie(r)
	http.SetCookie(w, myCookie)
	data := getRequestData(r, myCookie.Value)
	//fmt.Printf("[track] Request Data : %s\n", data)
	fmt.Printf("%s\n", data)
	/* Make sure that the GIF does not get cached */
	// private = don't cache in proxy, only browser can cache
	// no-cache= always revalidate when in cache
	// no-cache=Set-Cookie : you can store the content in cache, but not this header
	// proxy-revalidate: Make sure some possibly-broken proxies don't interfere
	w.Header().Add("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "image/gif")
	w.Write(oneGif)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/track", trackHandler)
	http.HandleFunc("/pixel.gif", pixelHandler)
	http.ListenAndServe(":8088", nil)
}
