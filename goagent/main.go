package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/phuslu/goproxy/httpproxy"
	"net/http"
	"time"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	addr := ":1080"
	ln, err := httpproxy.Listen("tcp4", addr)
	if err != nil {
		glog.Fatalf("Listen(\"tcp\", %s) failed: %s", addr, err)
	}
	h := httpproxy.Handler{
		Listener: ln,
		Net:      &httpproxy.SimpleNetwork{},
		RequestFilters: []httpproxy.RequestFilter{
			&httpproxy.StripRequestFilter{},
			&httpproxy.DirectRequestFilter{},
		},
		ResponseFilters: []httpproxy.ResponseFilter{
			&httpproxy.AlwaysRawResponseFilter{
				Sites: []string{"www.baidu.com"},
			},
			&httpproxy.ImageResponseFilter{},
			&httpproxy.RawResponseFilter{},
		},
	}
	s := &http.Server{
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	glog.Infof("ListenAndServe on %s\n", h.Listener.Addr().String())
	glog.Exitln(s.Serve(h.Listener))
}
