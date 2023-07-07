package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func main() {
	// 라우트 등록
	http.HandleFunc("/healthcheck", healthCheck)
	http.HandleFunc("/status", status)

	// 서버 시작
	log.Fatal(http.ListenAndServe(":8081", logRequest(http.DefaultServeMux)))
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

func status(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1"

	// HTTP GET 요청 보내기
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("GET 요청 실패: %s", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// 응답 바디(body) 읽기
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("응답 읽기 실패: %s", err), http.StatusInternalServerError)
		return
	}

	// 응답 출력
	fmt.Fprintln(w, string(body))
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
