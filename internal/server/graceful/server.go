package graceful

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CertConfig struct {
	CertFile string
	KeyFile  string
}

type httpServerOption struct {
	certs           *CertConfig
	shutdownTimeout time.Duration
}

type HttpServerOptionFn func(*httpServerOption)

func WithCert(cert *CertConfig) HttpServerOptionFn {
	return func(opt *httpServerOption) {
		opt.certs = cert
	}
}
func WithShutdownTimeout(d time.Duration) HttpServerOptionFn {
	return func(opt *httpServerOption) {
		opt.shutdownTimeout = d
	}
}

func Run(ctx context.Context, s *http.Server, opts ...HttpServerOptionFn) error {
	opt := httpServerOption{
		shutdownTimeout: 5 * time.Second,
	}

	for _, fn := range opts {
		fn(&opt)
	}

	// setup graceful server
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go

	/*s := http.Server{
		Addr:              bind,
		Handler:           handler,
		ReadHeaderTimeout: 1 * time.Minute,
		WriteTimeout:      3 * time.Minute,
	}*/

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return errors.WithStack(err)
	}

	errc := make(chan error)
	go func() {
		defer close(errc)

		logrus.Infof("Listening and serving HTTP on '%s'", s.Addr)
		if opt.certs == nil {
			if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
				errc <- err
			}
		} else {
			if err := s.ServeTLS(l, opt.certs.CertFile, opt.certs.KeyFile); err != nil && err != http.ErrServerClosed {
				errc <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		logrus.Warn("shutting down gracefully, press Ctrl+C again to force")
	case err := <-errc:
		logrus.Errorf("listen: %s\n", err)
	}

	// nCtx for shutdown timeout only
	nCtx, cancel := context.WithTimeout(context.Background(), opt.shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(nCtx); err != nil {
		return errors.Wrapf(err, "Server forced to shutdown")
	}

	return nil
}
