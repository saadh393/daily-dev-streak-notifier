package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

// Download image from URL
func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Crop the image
func cropImage(img image.Image, x, y, width, height int) image.Image {
	rect := image.Rect(0, 0, width, height)
	cropped := image.NewRGBA(rect)
	draw.Draw(cropped, rect, img, image.Point{x, y}, draw.Src)
	return cropped
}

// Convert image to PNG buffer
func saveImageToBuffer(img image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Send the cropped image to OCR.space
func extractTextFromImage(imgBuffer *bytes.Buffer) (string, error) {
	// OCR.space API endpoint
	apiURL := "https://api.ocr.space/parse/image"

	// Create a multipart form request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "cropped.png")
	if err != nil {
		return "", err
	}
	part.Write(imgBuffer.Bytes())
	writer.WriteField("language", "eng")      // English language OCR
	writer.WriteField("apikey", "helloworld") // Free anonymous API key
	writer.Close()

	// Send POST request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}

	// Extract text from JSON response
	if parsedResults, found := result["ParsedResults"].([]interface{}); found && len(parsedResults) > 0 {
		text := parsedResults[0].(map[string]interface{})["ParsedText"].(string)
		return text, nil
	}

	return "", fmt.Errorf("OCR failed: %s", string(respBody))
}

func extractFirstNumber(input string) string {
	// Split into lines
	lines := strings.Split(input, "\n")

	// Regular expression to find the first number
	re := regexp.MustCompile(`\d+`)

	// Find number in the first line
	if len(lines) > 0 {
		match := re.FindString(lines[0])
		return match
	}

	return ""
}

func ocr(userData User) string {
	// Step 1: Download the Image
	imageURL := fmt.Sprintf("https://api.daily.dev/devcards/v2/%s.png?type=default&r=cik", userData.CardId)

	img, err := downloadImage(imageURL)
	if err != nil {
		log.Fatal("Error downloading image:", err)
	}

	// Step 2: Crop the Image
	croppedImg := cropImage(img, 200, 384, 200, 100)

	// Step 3: Convert Image to Buffer
	imgBuffer, err := saveImageToBuffer(croppedImg)
	if err != nil {
		log.Fatal("Error converting image to buffer:", err)
	}

	// Save img buffer as file
	// // file, err := os.Create("cropped_image.png")
	// // if err != nil {
	// // 	fmt.Printf("Error creating file")
	// // }
	// // defer file.Close()

	// _, err = file.Write(imgBuffer.Bytes())
	// if err != nil {
	// 	log.Fatal("Error saving image file:", err)
	// }

	// Step 4: Extract text using OCR.space
	text, err := extractTextFromImage(imgBuffer)
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	number := extractFirstNumber(text)

	return number
}
