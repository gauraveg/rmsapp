package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gauraveg/rmsapp/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"io"
	"net/http"
)

func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                                                                                                                                                 // Allow all origins, adjust as necessary for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},                                                                                                                  // HTTP methods allowed by CORS
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma", "x-started-at", "x-api-key"}, // Allowed headers for CORS
		ExposedHeaders:   []string{"Link"},                                                                                                                                                              // Headers that are exposed to the client
		AllowCredentials: true,                                                                                                                                                                          // Allow credentials such as cookies or authorization headers
	})
}

func CommonMiddleware(loggers *logger.ZapLogger) chi.Middlewares {
	return chi.Chain(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		},
		corsOptions().Handler,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func(loggers *logger.ZapLogger) {
					err := loggers.Logger.Sync()
					if err != nil {
						loggers.ErrorWithContext(r.Context(), err.Error())
						return
					}
				}(loggers)

				requestId := uuid.New().String()
				if r.Context().Value("requestId") != "" {
					r = r.WithContext(context.WithValue(r.Context(), "requestId", requestId))
				}
				r = r.WithContext(context.WithValue(r.Context(), "logContext", loggers))

				//Logging the request body and adding the payload to context
				if body, err := io.ReadAll(r.Body); err == nil {
					payload := string(body)
					r = r.WithContext(context.WithValue(r.Context(), "payload", payload))
					var reqBody map[string]interface{}
					err := json.Unmarshal([]byte(payload), &reqBody)
					if err != nil {
						loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to unmarhsal request body", "error": err.Error()})
						return
					}
					delete(reqBody, "password")
					loggers.InfoWithContext(r.Context(), map[string]string{"requestBody": fmt.Sprintf("%v", reqBody)})
				}

				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					err := recover()
					if err != nil {
						jsonBody, _ := json.Marshal(map[string]string{
							"error": "Internal server error",
						})
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Request Panic err", "error": err.(string)})
						_, err := w.Write(jsonBody)
						if err != nil {
							loggers.ErrorWithContext(r.Context(), map[string]string{"message": "Failed to send response from middleware", "error": err.Error()})
						}
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	)
}
