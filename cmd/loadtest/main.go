package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	url := flag.String("url", "http://localhost:50470/api/v1/wallet", "endpoint")
	walletID := flag.String("wallet", "", "wallet uuid")
	op := flag.String("op", "DEPOSIT", "DEPOSIT or WITHDRAW")
	amount := flag.Int64("amount", 1, "amount")
	rps := flag.Int("rate", 1000, "requests per second")
	duration := flag.Duration("duration", 60*time.Second, "test duration")
	workers := flag.Int("workers", 200, "number of workers")
	flag.Parse()

	if *walletID == "" {
		log.Fatal("-wallet is required")
	}

	body := fmt.Sprintf(
		`{"wallet_id":%q,"operation_type":%q,"amount":%d}`,
		*walletID, *op, *amount,
	)

	// контекст только для лимитера/длительности теста
	testCtx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	limiter := rate.NewLimiter(rate.Limit(*rps), 1)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        2000,
			MaxIdleConnsPerHost: 2000,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	var ok, failed, serverErrors atomic.Int64
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for limiter.Wait(testCtx) == nil {
				req, err := http.NewRequestWithContext(
					context.Background(), // у запроса свой Timeout в клиенте
					http.MethodPost,
					*url,
					strings.NewReader(body),
				)
				if err != nil {
					failed.Add(1)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					failed.Add(1)
					continue
				}

				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()

				switch {
				case resp.StatusCode >= 500:
					serverErrors.Add(1)
				case resp.StatusCode == http.StatusOK:
					ok.Add(1)
				default:
					failed.Add(1)
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf(
		"ok=%d failed=%d 5xx=%d elapsed=%s actual_rps=%.0f\n",
		ok.Load(),
		failed.Load(),
		serverErrors.Load(),
		elapsed,
		float64(ok.Load())/elapsed.Seconds(),
	)

	if serverErrors.Load() > 0 {
		log.Printf("ВНИМАНИЕ: есть ответы 5xx: %d", serverErrors.Load())
	}
}
