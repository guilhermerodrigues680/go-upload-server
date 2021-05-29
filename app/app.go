package app

import (
	"net"
	"net/http"
	"strconv"

	"github.com/guilhermerodrigues680/gouploadserver/handler"

	"github.com/sirupsen/logrus"
)

func Run(wd string, port int, keepOriginalUploadFileName bool, spaMode bool, logger *logrus.Entry) error {
	logger.Info("** Go Upload Server **")
	logger.Infof("Working directory: %s", wd)

	h := handler.NewServer(wd, keepOriginalUploadFileName, spaMode, logger.WithField("server", "handler"))

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: h,
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Errorf("Interface Addrs error: %s", err)
		return err
	}

	for _, a := range addrs {
		// if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() { // to ignore 127.0.0.1
		if ipnet, ok := a.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				logger.Infof("Listening on: http://%s%s", ipnet.IP, srv.Addr)
			}
		}
	}

	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("/etc/letsencrypt/live/app.guilhermeonline.com.br/fullchain.pem", "/etc/letsencrypt/live/app.guilhermeonline.com.br/privkey.pem")
	if err != nil {
		logger.Errorf("Server error: %s", err)
		return err
	}

	return nil
}
