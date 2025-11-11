package proxy

import (
	"chimney3-go/core"
	"chimney3-go/privacy"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type ProxyServer interface {
	Serve()
	Close()
}

type proxyServer struct {
	Host     string
	Password string
	Which    string
	Exit     bool
	server   *http.Server
	listener core.MySSLListener
}

func (p *proxyServer) Serve() error {
	key := privacy.MakeCompressKey(p.Password)
	II := privacy.NewMethodWithName(p.Which)
	l, err := core.ListenSSL(p.Host, key, II)
	if err != nil {
		return err
	}
	p.listener = l

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	p.server = server

	return server.Serve(l)
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(" fatal error on udp server: ", err)
		}
	}()
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go transfer(dest_conn, client_conn, &wg)
	go transfer(client_conn, dest_conn, &wg)
	wg.Wait()
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (p *proxyServer) Close() error {
	// Try to gracefully stop the HTTP server and close listener.
	if p.server != nil {
		_ = p.server.Close()
	}
	if p.listener != nil {
		_ = p.listener.Close()
	}
	p.Exit = true
	return nil
}
