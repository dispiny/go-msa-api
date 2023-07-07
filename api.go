package main

import (
	"encoding/json"
	// "fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type StatusResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func main() {
	// 라우트 등록
	http.HandleFunc("/healthcheck", healthCheck)
	http.HandleFunc("/v1/status", status)

	// 서버 시작
	log.Fatal(http.ListenAndServe(":8080", logRequest(http.DefaultServeMux)))
}

// Health Check 핸들러
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Status Check 핸들러
services := map[string]string{
	"api":      ":8080",
	"status":   ":8081",
	"order":    ":8082",
	"token":    ":8083",
	"frontend": ":8084",
}

func status(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1"
	service_status := "bad"
	service := r.URL.Query().Get("service")
	
	if service == "all" {
				response, err := http.Get(url + services[service] + "/healthcheck")
		if err != nil {
			response := StatusResponse{
				Status:  service_status,
				Service: service,
			}
	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
	
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else {
			service_status = "good"
			response := StatusResponse{
				Status:  service_status,
				Service: service,
			}
	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
	
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		defer response.Body.Close()
	} else {
		// HTTP GET 요청 보내기
		response, err := http.Get(url + services[service] + "/healthcheck")
		if err != nil {
			response := StatusResponse{
				Status:  service_status,
				Service: service,
			}
	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
	
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else {
			service_status = "good"
			response := StatusResponse{
				Status:  service_status,
				Service: service,
			}
	
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
	
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		defer response.Body.Close()
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 로깅을 위한 logrus.Entry 생성
		logger := log.New(os.Stdout, "", 0)

		// responseWriter를 감싸는 구조체 생성
		rw := responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			size:           0,
		}

		// 원래의 핸들러 실행
		next.ServeHTTP(&rw, r)

		// 처리 시간 측정
		duration := time.Since(start)

		// 로그 작성
		logger.Printf("Method: %s, URI: %s, Status Code: %d, Content Length: %d, Duration: %f seconds\n",
			r.Method, r.RequestURI, rw.statusCode, rw.size, duration.Seconds())
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += int64(n)
	return n, err
}
