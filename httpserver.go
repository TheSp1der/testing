package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func startServer() error {
	cert := "-----BEGIN CERTIFICATE-----\n"
	cert += "MIIDJzCCAg+gAwIBAgIJAKGiNftFxNRyMA0GCSqGSIb3DQEBBQUAMCoxKDAmBgNV\n"
	cert += "BAMMH25pY29sYWFzLXBjLnRlc3QtY2hhbWJlci0xMy5sYW4wHhcNMTgxMDExMDI0\n"
	cert += "MDU1WhcNMjgxMDA4MDI0MDU1WjAqMSgwJgYDVQQDDB9uaWNvbGFhcy1wYy50ZXN0\n"
	cert += "LWNoYW1iZXItMTMubGFuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\n"
	cert += "qpkl8L9iR0ZnXxtTtxUD/eF75uyOFunRPmEtvnFQMIuKMbX5TptXIJBdLuNpm+M3\n"
	cert += "jUw3qkYSrBm/WkdA/FZSPAiaUwernQNYQLdRw/ltfq38L1WMdvPwrSMgTErYwBGR\n"
	cert += "YmNbWpM+sdkiZFn5x2qECN6kDVrO+8C7atmv8T6Dyniw8+SM0F6Gq7l9OazBXcNM\n"
	cert += "h9CXHLaHEQ6LTXa+E3VsZyKAqJfkH34r49AhhGhZwPzhDpMIhSTxtJTOZAaxqNFP\n"
	cert += "f6Stga9X7jQvTMcL7N1V3SmnFd3RDLx0joyJNd4VJ1nXLbqJC3CyOSUcda48sj6N\n"
	cert += "ZyPSPOHGiV9FPcYNUim3KwIDAQABo1AwTjAdBgNVHQ4EFgQUEC2IfS/XQhbL/IpX\n"
	cert += "MvY/nWOeqfcwHwYDVR0jBBgwFoAUEC2IfS/XQhbL/IpXMvY/nWOeqfcwDAYDVR0T\n"
	cert += "BAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAc4FbWXtkxcZjAl5Ru0JRiFD8por4\n"
	cert += "rUl+974gacDQRK26/lvl5/urbrP7KqZ2fTBk66ypsIGVHiMT/1TWUx6S5gE2BEOU\n"
	cert += "uZlYPe32pqkAqM3En1bC0r6FG/nDX53FgSzyZhGeM5ekWGmhMKDj/pyJ3QdwRR1L\n"
	cert += "lfMlt5nuk1Q0tysLkr/X9D6YecFv3Ce6jHcAcUaz1RKsfUPG4izNw2zFa8J5XjEJ\n"
	cert += "NsnmM6mTzqyslrqwJosH9fyMceIvEir+0zpNGgZrvyloZz5/1g46MjPVWFYgtKZw\n"
	cert += "vtv4t7q69PYQTRudSp0f0pc1ccxS2r/o6iRMQE0gZECNlyFsHXq5gX3j3A==\n"
	cert += "-----END CERTIFICATE-----"

	key := "-----BEGIN RSA PRIVATE KEY-----\n"
	key += "MIIEogIBAAKCAQEAqpkl8L9iR0ZnXxtTtxUD/eF75uyOFunRPmEtvnFQMIuKMbX5\n"
	key += "TptXIJBdLuNpm+M3jUw3qkYSrBm/WkdA/FZSPAiaUwernQNYQLdRw/ltfq38L1WM\n"
	key += "dvPwrSMgTErYwBGRYmNbWpM+sdkiZFn5x2qECN6kDVrO+8C7atmv8T6Dyniw8+SM\n"
	key += "0F6Gq7l9OazBXcNMh9CXHLaHEQ6LTXa+E3VsZyKAqJfkH34r49AhhGhZwPzhDpMI\n"
	key += "hSTxtJTOZAaxqNFPf6Stga9X7jQvTMcL7N1V3SmnFd3RDLx0joyJNd4VJ1nXLbqJ\n"
	key += "C3CyOSUcda48sj6NZyPSPOHGiV9FPcYNUim3KwIDAQABAoIBADHWr/jXUJTWApkM\n"
	key += "WLah0xq2ZwYdkZ0sDc8VgNGkNPMZsPO43+6Q/zEqO67ZDR9XkAEdhR2ffxD8LKTp\n"
	key += "MBkH9tpHAR7EnOQv9/ZgF+kS02Qw2/3QFksiFOvf2S2wqAXkm/6MXEHnxmcasitz\n"
	key += "Bb+2ZIBa2r50CwgNVDNxCS+HPeVGRs8F2DwLNht1YDyLKKYgr1oBuFtmVu8PKyFY\n"
	key += "8BuXbRLOSZ3ji1o5PJPqp9AdZuDPQx3ZXlcSdOlUIEUVJ9caEMFR72O11Drsvg91\n"
	key += "xr+EsdYjYlazhpzL87qSdBgOht33hxmVbJK3bEUQBUFn9WUVRph+PzNsQLB8Qexl\n"
	key += "APWPX6ECgYEA3R3Mr9DvxcLpvrru4+mdwrPPZbfLkxYJzyCmdhaeSUWxkERxUREN\n"
	key += "mIbwEvQwnYyq5wXhBUwenMWhPgqDtuYnfSqBVGs+qCoShinP2hzDQqdVOeXOdNbx\n"
	key += "dY4syd5M04tAi3gHQ7uz01aSDO2M4Nvnkr+Ngc+oB7GfS5a7JC9mIhMCgYEAxYMT\n"
	key += "439fiRA8ZoPwCypEighnh3zzDScWdl8xwzvgkXbqsGwarzf9HsePDzHiA6ENAqgf\n"
	key += "u2NfivtK0/a1RMdKfR1Ks9q1kw/g00CZCJzzvF5GIGkhBN3pXis2SeKJYb80o29c\n"
	key += "tIdw7/CcLObp/gohO3QLVzk4TMdfB8pNgbDP+YkCgYAo7dIsnS002w5vWqTLlTu5\n"
	key += "hZUXS/0nvcWVDIMjiq5D+92RScn76n8sw5V+vKqfDyG3X7Q2Sc/EzyQ4mrOk0Fdw\n"
	key += "6MRFvxA7CoahRO4PfpF6LgUtkWc043CQhP+vYjGwWq9Y4Z/ensj7jqO8NuCD4tCr\n"
	key += "rj9gTvLYcb19vWnomcl69wKBgFKXfmCaacu581f3AhDZKvIBk7FPaZ9tYfI72mZG\n"
	key += "iqCpdngxrHLq2bjeQA9dj6Ju3S7oOOS2KETI0kCSoLhTEe4BqrRM17LYZ+5Oy++T\n"
	key += "GkUBsxdofrs3RJfxP/FjfolWWF+jeMOxA2QCXHxWTzDA8aaX3wopTkak9DMgwIpj\n"
	key += "8oPpAoGANqYUFBDxcC0TyZcJffTmWAo4fcXivaEt/U00QIp+SLZRR/2XSASRoX3s\n"
	key += "6QSVrw+1S4ZozzRA7s9wosWzUOxk4bykTbd7oVP0rxqdpPvIPCvIRvlHM6+M15m6\n"
	key += "bDNASPeykySbreyKVY4aFOVkN0F9UXgKGIjGeviadFqz12oKu8k=\n"
	key += "-----END RSA PRIVATE KEY-----"

	tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		if req.Method == "GET" {
			w.Write([]byte("This is an example server.\n"))
			w.Write([]byte("Method: " + req.Method + "\n"))
			w.Write([]byte("UserAgent: " + req.UserAgent() + "\n"))
		}
	})

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{tlsCert},
	}

	srv := &http.Server{
		Addr:              ":8443",
		Handler:           mux,
		TLSConfig:         cfg,
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		// time to read request headers
		ReadHeaderTimeout: time.Duration(2 * time.Second),
		// time from accept to full request body read
		ReadTimeout:       time.Duration(15 * time.Second),
		// non tls: time from request header read to the end of the response
		// tls:     time from accept to end of response
		WriteTimeout:      time.Duration(10 * time.Second),
		// time a Keep-Alive connection will be kept idle
		IdleTimeout:       time.Duration(120 * time.Second),
	}

	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}

	return nil
}
