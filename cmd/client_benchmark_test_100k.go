package main_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	serverURL         = "http://localhost:8888" // URL сервера для тестирования
	requestsPerSecond = 100_000                 // Целевая нагрузка: 100,000 запросов в секунду
)

// BenchmarkLoadTest тестирует производительность сервера при нагрузке 100,000 RPS с случайными pageID
func BenchmarkLoad2Test(b *testing.B) {
	// Создаем HTTP-клиент с таймаутом 10 секунд
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Тестируем эндпоинт /increment
	b.Run("Increment", func(b *testing.B) {
		// Создаем тикер для ограничения скорости до 100,000 RPS
		rateLimiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer rateLimiter.Stop() // Освобождаем тикер после завершения

		var wg sync.WaitGroup // WaitGroup для синхронизации горутин
		b.ResetTimer()        // Сбрасываем таймер бенчмарка перед началом

		// Запускаем b.N итераций (определяется Go автоматически)
		for i := 0; i < b.N; i++ {
			wg.Add(1) // Увеличиваем счетчик ожидающих горутин
			go func() {
				defer wg.Done()               // Уменьшаем счетчик по завершении горутины
				<-rateLimiter.C               // Ожидаем тик тикера для ограничения скорости
				pageID := uuid.New().String() // Генерируем случайный UUID для pageID
				// Создаем POST-запрос для эндпоинта /increment/{pageID}
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/increment/%s", serverURL, pageID), nil)
				if err != nil {
					b.Errorf("Failed to create request: %v", err) // Логируем ошибку создания запроса
					return
				}
				// Выполняем запрос
				resp, err := client.Do(req)
				if err != nil {
					b.Errorf("Request failed: %v", err) // Логируем ошибку выполнения запроса
					return
				}
				defer resp.Body.Close() // Закрываем тело ответа
				// Проверяем статус ответа
				if resp.StatusCode != http.StatusOK {
					b.Errorf("Unexpected status code: %d", resp.StatusCode)
				}
			}()
		}
		wg.Wait() // Ожидаем завершения всех горутин
	})

	// Тестируем эндпоинт /count
	b.Run("GetCount", func(b *testing.B) {
		// Создаем тикер для ограничения скорости до 100,000 RPS
		rateLimiter := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer rateLimiter.Stop() // Освобождаем тикер после завершения

		var wg sync.WaitGroup // WaitGroup для синхронизации горутин
		b.ResetTimer()        // Сбрасываем таймер бенчмарка перед началом

		// Запускаем b.N итераций
		for i := 0; i < b.N; i++ {
			wg.Add(1) // Увеличиваем счетчик ожидающих горутин
			go func() {
				defer wg.Done()               // Уменьшаем счетчик по завершении горутины
				<-rateLimiter.C               // Ожидаем тик тикера для ограничения скорости
				pageID := uuid.New().String() // Генерируем случайный UUID для pageID
				// Выполняем GET-запрос для эндпоинта /count/{pageID}
				resp, err := client.Get(fmt.Sprintf("%s/count/%s", serverURL, pageID))
				if err != nil {
					b.Errorf("Request failed: %v", err) // Логируем ошибку выполнения запроса
					return
				}
				defer resp.Body.Close() // Закрываем тело ответа
				// Проверяем статус ответа
				if resp.StatusCode != http.StatusOK {
					b.Errorf("Unexpected status code: %d", resp.StatusCode)
				}
			}()
		}
		wg.Wait() // Ожидаем завершения всех горутин
	})
}
