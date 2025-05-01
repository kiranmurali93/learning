package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Ignore this error
func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(rid int) {
			defer wg.Done()
			url := fmt.Sprintf("http://localhost:3000/claim-job?jobId=1&recruiterId=%d", rid)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Recruiter %d error: %v\n", rid, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Recruiter %d: %s\n", rid, resp.Status)
		}(i)
	}

	wg.Wait()
}
