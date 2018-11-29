package server

import (
    "fmt"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	"google.golang.org/grpc"
    "github.com/improbable-eng/grpc-web/go/grpcweb"
    "github.com/GeertJohan/go.rice"
	"net"
    "net/http"
    "log"
    "time"
)

func Start(cfg *Config, srv api.SessionServiceServer) error {
    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)
    api.RegisterSessionServiceServer(grpcServer, srv)

    if cfg.GRPC.Enabled {
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
        if err != nil {
            return err
        }

        go func() { log.Fatal(grpcServer.Serve(lis)) }()
    }

    if cfg.UI.Enabled {
        wrappedGrpc := grpcweb.WrapServer(grpcServer)
        srv := &http.Server{
            // These interfere with websocket streams, disable for now
            // ReadTimeout: 5 * time.Second,
            // WriteTimeout: 10 * time.Second,
            ReadHeaderTimeout: 5 * time.Second,
            IdleTimeout:       120 * time.Second,
            Addr:              fmt.Sprintf(":%d", cfg.UI.Port),
            Handler: hstsHandler(
                grpcTrafficSplitter(
                    http.FileServer(rice.MustFindBox("../../client/build").HTTPBox()),
                    wrappedGrpc,
                ),
            ),
        }
        go func() { log.Fatal(srv.ListenAndServe()) }()
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
        }

        // Fall back to other servers.
        fallback.ServeHTTP(resp, req)
	})
}