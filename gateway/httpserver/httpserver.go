package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tolling/common"
)

type Server struct {
	Addr    string
	targets map[string]string
	logger  common.Logger
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func NewServer(targets map[string]string, addr string, logger common.Logger) *Server {
	return &Server{
		Addr:    addr,
		logger:  logger,
		targets: targets,
	}
}

func (s *Server) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) Start() error {
	// register endpoint
	http.HandleFunc("/api/", s.handleError(s.handleRequest))

	// start server
	s.logger.New().Infof("starting API-Gateway on port %s", s.Addr)
	return http.ListenAndServe(s.Addr, nil)
}

func (s *Server) getOrigin(requestURL string) (string, error) {

	parts := strings.Split(requestURL, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid request URL: %s", requestURL)
	}
	targetKey := parts[2]

	targetAddr, ok := s.targets[targetKey]
	if !ok {
		return "", fmt.Errorf("not registered target %s", targetKey)
	}

	remainingPath := strings.Join(parts[3:], "/")
	targetURL := targetAddr
	if remainingPath != "" {
		targetURL = targetURL + "/" + remainingPath
	}

	return targetURL, nil
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) error {
	targetURL, err := s.getOrigin(r.URL.Path)
	if err != nil {
		return err
	}

	client := &http.Client{}

	targetReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return err
	}

	targetReq.Header = r.Header
	targetReq.URL.RawQuery = r.URL.RawQuery

	resp, err := client.Do(targetReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("rrror copying response body:", err)
		return err
	}

	return nil
}

func (s *Server) handleError(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			fmt.Println("writing error")
			s.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		fmt.Println("writing oK")

	}
}
