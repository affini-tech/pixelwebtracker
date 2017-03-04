# pixelwebtracker
A simple pixel cookie web tracker

## output
```
{"CookieIdentifier" : "z22L6REjrtejrethrtrternlONCY1p","DateTime" : "2017-03-04 17:23:56.287241154 +0000 UTC","Referer" : "http://some.website.com/page.html","UserAgent" : "Mozilla/5.0 (Android 6.0; Mobile; rv:51.0) Gecko/51.0 Firefox/51.0","Host" : "","EscapedPath" : "/pixel.gif","Url" : "/pixel.gif","RemoteAddr" : "127.0.0.1:57562","TrueRemoteAddr" : "12.34.56.78","RequestURI" : "/pixel.gif","Query" : ""}
```


## configuration of the reverse proxy to get the real IP address
```
server {
	listen 80;
	listen [::]:80;

	server_name tracker.xxxx.xyz;
	access_log /var/log/nginx/tracker.access.log;
	error_log /var/log/nginx/tracker.error.log;
	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_pass	http://localhost:8088/;
	}
}
```
