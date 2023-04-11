package pkce

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

// CodeResponse represents the code received by the local server's callback handler.
type CodeResponse struct {
	Code  string
	State string
}

// bindLocalServer initializes a LocalServer that will listen on a given TCP port.
func bindLocalServer(addr string) (*localServer, error) {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}

	return &localServer{
		listener:   listener,
		resultChan: make(chan CodeResponse, 1),
	}, nil
}

type localServer struct {
	CallbackPath     string
	WriteSuccessHTML func(w io.Writer)

	resultChan chan CodeResponse
	listener   net.Listener
}

func (s *localServer) Close() error {
	return s.listener.Close()
}

func (s *localServer) Serve() error {
	return http.Serve(s.listener, s) //nolint: gosec
}

func (s *localServer) WaitForCode() (CodeResponse, error) {
	return <-s.resultChan, nil
}

// ServeHTTP implements http.Handler.
func (s *localServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.CallbackPath != "" && r.URL.Path != s.CallbackPath {
		w.WriteHeader(404)
		return
	}
	defer func() {
		_ = s.Close()
	}()

	params := r.URL.Query()
	s.resultChan <- CodeResponse{
		Code:  params.Get("code"),
		State: params.Get("state"),
	}

	w.Header().Add("content-type", "text/html")
	if s.WriteSuccessHTML != nil {
		s.WriteSuccessHTML(w)
	} else {
		defaultSuccessHTML(w)
	}
}

func defaultSuccessHTML(w io.Writer) {
	fmt.Fprintf(w, `
	<html>
		<head>
			<link href="https://fonts.googleapis.com/css?family=Roboto:300&display=swap" rel="stylesheet">  
			<title>Aserto</title>
		</head>
		<body style="background: #000; color: #e7e7e7; font-family: 'Roboto', -apple-system, BlinkMacSystemFont, sans-serif; font-weight: 300;">
			<center style="margin: 100">
				<img src="https://aserto.com/images/Aserto-logo-color-120px.png" alt="aserto" width="120" />
				<h2>You've logged in successfully.</h2>
				<h3>You can close this window and return to the aserto CLI.</h3>
			</center>
		</body>
	</html>`)
}
