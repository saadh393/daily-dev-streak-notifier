package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// User holds the scraped data from Daily.dev
type User struct {
	Name       string `json:"name"`
	Reputation int    `json:"reputation"`
}

// CacheData holds the user data, profile URL, and last update timestamp.
type CacheData struct {
	User       User      `json:"user"`
	Timestamp  time.Time `json:"timestamp"`
	ProfileURL string    `json:"profile_url"`
}

// getCacheFilePath returns the path of the JSON cache file.
func getCacheFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".dailydev_data.json")
}

// loadCache reads and unmarshals cache data from the file.
func loadCache() (*CacheData, error) {
	filePath := getCacheFilePath()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

// saveCache saves the CacheData to the JSON file.
func saveCache(cache *CacheData) error {
	filePath := getCacheFilePath()
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

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

// promptProfileURL prompts the user for their Daily.dev profile URL.
func promptProfileURL() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Daily.dev profile URL (e.g., https://app.daily.dev/username): ")
	url, _ := reader.ReadString('\n')
	return strings.TrimSpace(url)
}

// installOnStartup appends the executable’s path to the shell profile (e.g., .bashrc or .zshrc)
// so that the program runs every time a new terminal session starts.
func installOnStartup() {
	shell := os.Getenv("SHELL")
	var profilePath string
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Choose the appropriate profile file.
	if strings.Contains(shell, "zsh") {
		profilePath = filepath.Join(usr.HomeDir, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		profilePath = filepath.Join(usr.HomeDir, ".bashrc")
	} else {
		profilePath = filepath.Join(usr.HomeDir, ".profile")
	}

	// Get current executable path.
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	// Check if the profile file already contains our command.
	data, err := ioutil.ReadFile(profilePath)
	if err == nil {
		if strings.Contains(string(data), exePath) {
			// Already installed.
			return
		}
	}

	// Append the command to the shell profile.
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening shell profile:", err)
		return
	}
	defer f.Close()
	line := fmt.Sprintf("\n# Daily.dev Reputation Utility\n%s\n", exePath)
	if _, err := f.WriteString(line); err != nil {
		fmt.Println("Error writing to shell profile:", err)
		return
	}
	fmt.Println("Installed to run on terminal startup. Please restart your terminal.")
}

// displayReputation prints the welcome message and reputation in a fancy format.
func displayReputation(user *User) {
	fmt.Println("========================================")
	fmt.Printf(" 🚀 Welcome back, %s! 🚀\n", user.Name)
	fmt.Println("========================================")
	fmt.Printf(" ⭐ Your current reputation: %d ⭐\n", user.Reputation)
	fmt.Println("========================================")
}

func main() {
	// Install the binary to run on every terminal startup.
	installOnStartup()

	// Load the cached data if it exists.
	cache, err := loadCache()
	if err != nil || cache.ProfileURL == "" {
		// First run: prompt for profile URL.
		profileURL := promptProfileURL()

		cache = &CacheData{
			ProfileURL: profileURL,
		}
	}

	now := time.Now()
	// If data is older than 24 hours, scrape fresh data.
	if cache.Timestamp.IsZero() || now.Sub(cache.Timestamp) > 24*time.Hour {
		userData, err := scrapeUserData(cache.ProfileURL)
		if err != nil {
			fmt.Println("Error scraping user data:", err)
			return
		}
		fmt.Printf("Profile URL: %s\n", cache.ProfileURL)
		cache.User = *userData
		cache.Timestamp = now
		if err := saveCache(cache); err != nil {
			fmt.Println("Error saving cache:", err)
		}
	}

	// Display the reputation in the terminal.
	displayReputation(&cache.User)
}
