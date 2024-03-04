package proxy

import (
	"awesome-proxy/config"
	"github.com/gorilla/handlers"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"os"
)

type responder struct{}

func (r responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.Errorf("Error while proxying request: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Run(config *config.Config) {
	responder := &responder{}
	for _, p := range config.Proxy {
		p := p
		go func() {
			proxyPass, err := url.Parse(p.ProxyPass)
			if err != nil {
				klog.Errorf("Can't parse URL: %s", p.ProxyPass)
				return
			}

			transport := &Transport{
				PathPrepend: p.Location,
			}

			proxy := NewUpgradeAwareHandler(proxyPass, transport, false, false, responder)
			proxy.UseRequestLocation = true
			proxy.UseLocationHost = true
			http.HandleFunc(p.Location, func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set("X-Proxy-Prefix", p.Location)
				proxy.ServeHTTP(w, r)
			})
		}()
	}

	http.ListenAndServe(config.Server.Listen, handlers.CustomLoggingHandler(os.Stdout, http.DefaultServeMux, writeCustomLog))
}
