package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/Bipolar-Penguin/bff-website/pkg/service"
)

const (
	authHeader = "Authorization"
)

var (
	errNotFound            error = errors.New("not found")
	errUnprocessableEntity error = errors.New("unprocessable entity")
)

type authRequest struct {
	UserID string `json:"user_id"`
}

type HTTPServer struct {
	port        int
	Logger      log.Logger
	Origins     string
	Router      *mux.Router
	ReadTimeout time.Duration
	service     *service.Service
	http.Server
}

type errorResponse struct {
	Message string `json:"error"`
}

func NewHttpServer(port int, logger log.Logger, service *service.Service) *HTTPServer {
	srv := &HTTPServer{
		port:    port,
		Logger:  logger,
		Router:  mux.NewRouter(),
		service: service,
	}
	srv.configureRouter()

	return srv
}

func appendJSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func (s *HTTPServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(authHeader)
		fmt.Println("got header", authToken)
		userID, err := s.service.Authenticate(authToken)
		if err != nil {
			s.abortWithError(rw, http.StatusUnauthorized, err)
			return
		}
		fmt.Println("GOT USERID", userID)

		if userID == "" {
			s.abortWithError(rw, http.StatusUnauthorized, errors.New("not authorized"))
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func (s *HTTPServer) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	var jsonBody []byte

	var err error

	jsonBody, err = json.Marshal(data)
	if err != nil {
		s.abortWithError(rw, http.StatusUnprocessableEntity, err)
	}

	rw.WriteHeader(code)

	rw.Write(jsonBody)
}

func (s *HTTPServer) abortWithError(rw http.ResponseWriter, code int, err error) {
	var res errorResponse

	rw.WriteHeader(code)

	res.Message = err.Error()
	json.NewEncoder(rw).Encode(res)
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *HTTPServer) getNotFoundHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.abortWithError(rw, http.StatusNotFound, errNotFound)
	}
}

// ServeHTTP is http.Handler implementation
func (s *HTTPServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(rw, r)
}

// ListenAndServe overrides http.Server method
func (s *HTTPServer) ListenAndServe() error {
	lport := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(lport, s)
}

func (s *HTTPServer) configureRouter() {
	s.Router.NotFoundHandler = s.getNotFoundHandler()
	s.Router.Use(appendJSONHeader)
	s.Router.Use(s.authenticateUser)

	user := s.Router.PathPrefix("/user").Subrouter()
	{
		user.HandleFunc("/", s.saveUser).Methods(http.MethodPost, http.MethodOptions)
		user.HandleFunc("/auth", s.authorizeUser).Methods(http.MethodPost, http.MethodOptions)
	}
}

func (s *HTTPServer) saveUser(rw http.ResponseWriter, r *http.Request) {
	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.abortWithError(rw, http.StatusUnprocessableEntity, errUnprocessableEntity)
		return
	}

	user, err := s.service.SaveUser(user)
	if err != nil {
		s.abortWithError(rw, http.StatusInternalServerError, err)
		return
	}

	s.respond(rw, r, http.StatusOK, user)
}

func (s *HTTPServer) authorizeUser(rw http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.abortWithError(rw, http.StatusInternalServerError, err)
		return
	}

	token, err := s.service.GenerateToken(req.UserID)
	if err != nil {
		s.abortWithError(rw, http.StatusInternalServerError, err)
		return
	}

	s.respond(rw, r, http.StatusOK, token)
}
