package server

import (
	"fmt"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/GeertJohan/go.rice"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	stdliblog "log" // this has to be log for double header reporter
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// debugLogger is from https://rocketeer.be/blog/2018/01/multiple-response-writeheader-calls/
type debugLogger struct{}

func (d debugLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, "multiple response.WriteHeader") {
		debug.PrintStack()
	}
	return os.Stderr.Write(p)
}

func Start(cfg *Config, srv api.SessionServiceServer) error {
	var opts []grpc.ServerOption

	if cfg.TLS.Enabled {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLS.Cert, cfg.TLS.Key)
		if err != nil {
			return err
		}
		tls := grpc.Creds(creds)
		opts = append(opts, tls)
		log.WithField("certFile", cfg.TLS.Cert).WithField("keyFile", cfg.TLS.Key).Info("Enabling transport layer security")
	}

	grpcServer := grpc.NewServer(opts...)
	api.RegisterSessionServiceServer(grpcServer, srv)

	if cfg.GRPC.Enabled {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			return err
		}

		go func() { log.Fatal(grpcServer.Serve(lis)) }()
		log.WithField("port", cfg.GRPC.Port).Info("gRPC API server started")
	}

	if cfg.UI.Enabled {
		wrappedGrpc := grpcweb.WrapServer(grpcServer)
		handler := hstsHandler(
			grpcTrafficSplitter(
				http.FileServer(rice.MustFindBox("../../client/build").HTTPBox()),
				wrappedGrpc,
			),
		)

		// Now use the logger with your http.Server:
		logger := stdliblog.New(debugLogger{}, "", 0)

		srv := &http.Server{
			// These interfere with websocket streams, disable for now
			// ReadTimeout: 5 * time.Second,
			// WriteTimeout: 10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       120 * time.Second,
			Addr:              fmt.Sprintf(":%d", cfg.UI.Port),
			Handler:           handler,
			ErrorLog:          logger,
		}
		if cfg.UI.TLS {
			log.WithField("certFile", cfg.UI.Cert).WithField("keyFile", cfg.UI.Key).Info("Serving UI over HTTPS")
			go func() { log.Fatal(srv.ListenAndServeTLS(cfg.UI.Cert, cfg.UI.Key)) }()
		} else {
			go func() { log.Fatal(srv.ListenAndServe()) }()
		}
		log.WithField("port", cfg.UI.Port).Info("UI web server started")
	}

	return nil
}

// hstsHandler wraps an http.HandlerFunc such that it sets the HSTS header.
func hstsHandler(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		fn(w, r)
	})
}

func grpcTrafficSplitter(fallback http.Handler, wrappedGrpc *grpcweb.WrappedGrpcServer) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			// Fall back to other servers.
			fallback.ServeHTTP(resp, req)
		}
	})
}
