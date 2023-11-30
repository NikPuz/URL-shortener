package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/jxskiss/base62"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Response struct {
	Count   int `json:"count"`
	Entries []struct {
		List string `json:"Link"`
	} `json:"entries"`
}

func main() {
	response := collectLinks()

	testRecord(response)

	testReading(response)
}

func testRecord(response *Response) {
	fmt.Printf("\ntestRecord\n\n")

	OkCount := atomic.Uint64{}
	wg := sync.WaitGroup{}
	timeStart := time.Now()

	for i := 1; i <= 1000 && i < len(response.Entries); i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			resp, err := http.DefaultClient.Get(fmt.Sprintf("http://localhost:8080/a/?url=%s", response.Entries[i].List))
			if err != nil {
				panic(err.Error())
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				OkCount.Add(1)
			} else {
				fmt.Printf("%d\t|\t%s\t|\t\"%s\"\n", i, resp.Status, response.Entries[i])
			}

			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	fmt.Printf("1000 Запросов за %.3f секунд\n", time.Duration(time.Now().UnixNano()-timeStart.UnixNano()).Seconds())
	fmt.Printf("200 OK - %d\n", OkCount.Load())
}

func testReading(response *Response) {
	fmt.Printf("\ntestReading\n\n")

	shortURLs := make([]string, 0, len(response.Entries))

	for _, v := range response.Entries {
		sha := sha256.New()
		sha.Write([]byte(v.List))
		hash := sha.Sum(nil)

		shortURLs = append(shortURLs, string(base62.Encode(hash)[:8]))
	}

	foundCount := atomic.Uint64{}
	wg := sync.WaitGroup{}
	timeStart := time.Now()

	for i := 1; i <= 1000 && i < len(shortURLs); i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			client := new(http.Client)
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
			resp, err := client.Get(fmt.Sprintf("http://localhost:8080/s/%s", shortURLs[i]))
			if err != nil {
				fmt.Printf("%d\t|\t%s\t|\t\"%s\"\n", i, err.Error(), shortURLs[i])
				wg.Done()
				return
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusFound {
				foundCount.Add(1)
			} else {
				fmt.Printf("%d\t|\t%s\t|\t\"%s\"\n", i, resp.Status, shortURLs[i])
			}

			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	fmt.Printf("1000 Запросов за %.3f секунд\n", time.Duration(time.Now().UnixNano()-timeStart.UnixNano()).Seconds())
	fmt.Printf("302 Found - %d\n", foundCount.Load())
}

func collectLinks() *Response {
	resp, err := http.DefaultClient.Get("https://api.publicapis.org/entries")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	return &response
}
