package main_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	serverURL         = "http://localhost:8888"
	pageID            = "page1"
	requestsPerSecond = 10_000
)

// BenchmarkLoadTest тестирует производительность сервера при нагрузке 1000 RPS
func BenchmarkLoadTest(b *testing.B) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Тестируем эндпоинт /increment
	b.Run("Increment", func(b *testing.B) {
		// Запускаем 1000 запросов в секунду
		rateLimiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer rateLimiter.Stop()

		var wg sync.WaitGroup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-rateLimiter.C // Ограничиваем до 1000 RPS
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/increment/%s", serverURL, pageID), nil)
				if err != nil {
					b.Errorf("Failed to create request: %v", err)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					b.Errorf("Request failed: %v", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					b.Errorf("Unexpected status code: %d", resp.StatusCode)
				}
			}()
		}
		wg.Wait()
	})

	// Тестируем эндпоинт /count
	b.Run("GetCount", func(b *testing.B) {
		rateLimiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer rateLimiter.Stop()

		var wg sync.WaitGroup
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-rateLimiter.C
				resp, err := client.Get(fmt.Sprintf("%s/count/%s", serverURL, pageID))
				if err != nil {
					b.Errorf("Request failed: %v", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					b.Errorf("Unexpected status code: %d", resp.StatusCode)
				}
			}()
		}
		wg.Wait()
	})
}
