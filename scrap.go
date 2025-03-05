package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

// scrapeUserData scrapes the Daily.dev profile page using the provided URL.
func scrapeUserData(profileURL string) (*User, error) {
	c := colly.NewCollector()
	var data map[string]interface{}
	var userData User
	var scrapeErr error

	var wg sync.WaitGroup
	wg.Add(1) // Add one goroutine to wait for

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		defer wg.Done() // Ensure Done() is always called

		err := json.Unmarshal([]byte(e.Text), &data)
		if err != nil {
			scrapeErr = fmt.Errorf("JSON parse error: %v", err)
			return
		}

		props, ok := data["props"].(map[string]interface{})
		if !ok {
			scrapeErr = fmt.Errorf("Error extracting props")
			return
		}

		pageProps, ok := props["pageProps"].(map[string]interface{})
		if !ok {
			scrapeErr = fmt.Errorf("Error extracting pageProps")
			return
		}

		uData, ok := pageProps["user"].(map[string]interface{})
		if !ok {
			scrapeErr = fmt.Errorf("Error extracting user data")
			return
		}

		userJSON, err := json.Marshal(uData)
		if err != nil {
			scrapeErr = fmt.Errorf("Error marshaling user data: %v", err)
			return
		}

		id, ok := uData["id"].(string)
		if !ok {
			scrapeErr = fmt.Errorf("Error extracting user id")
			return
		}
		userData.CardId = id

		fmt.Println(id)

		if err := json.Unmarshal(userJSON, &userData); err != nil {
			scrapeErr = fmt.Errorf("Error unmarshaling into User struct: %v", err)
			return
		}
	})

	// Handle request errors
	c.OnError(func(_ *colly.Response, err error) {
		scrapeErr = fmt.Errorf("Request failed: %v", err)
		wg.Done() // Avoid deadlock if request fails
	})

	// Start scraping
	if err := c.Visit(profileURL); err != nil {
		return nil, fmt.Errorf("Failed to visit URL: %v", err)
	}

	// Wait for completion
	wg.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}
	return &userData, nil
}
