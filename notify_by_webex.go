package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

const maxFileSize = 100 * 1024 * 1024 // 100 MB

// allowedExtensions lists supported file types.
var allowedExtensions = map[string]bool{
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".pdf":  true,
	".jpg":  true,
	".jpeg": true,
	".bmp":  true,
	".gif":  true,
	".png":  true,
}

func main() {
	// Define command-line flags.
	tokenFlag := flag.String("T", "", "Webex access token")
	roomIDFlag := flag.String("R", "", "Room ID")
	textFlag := flag.String("t", "", "Message text (optional)")
	fileFlag := flag.String("f", "", "File path to upload")
	flag.Parse()

	// Check that required flags were provided.
	if *tokenFlag == "" || *roomIDFlag == "" || *fileFlag == "" {
		fmt.Println("Usage: -T <token> -R <roomId> -f <filename> [-t <text>]")
		os.Exit(1)
	}

	// Check file existence and its size.
	fileInfo, err := os.Stat(*fileFlag)
	if err != nil {
		fmt.Printf("Error accessing file: %v\n", err)
		os.Exit(1)
	}
	if fileInfo.Size() > maxFileSize {
		fmt.Println("Error: file exceeds maximum allowed size of 100 MB")
		os.Exit(1)
	}

	// Check that the file extension is supported.
	ext := strings.ToLower(filepath.Ext(*fileFlag))
	if _, ok := allowedExtensions[ext]; !ok {
		fmt.Println("Error: file type not supported")
		os.Exit(1)
	}

	// Open the file.
	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Prepare a buffer and a multipart writer.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Determine the MIME type for the file.
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Create a custom part for the file with proper headers.
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="files"; filename="%s"`, filepath.Base(*fileFlag)))
	h.Set("Content-Type", mimeType)

	// Log the headers.
	fmt.Printf("Content-Disposition: %s\n", h.Get("Content-Disposition"))
	fmt.Printf("Content-Type: %s\n", h.Get("Content-Type"))

	filePart, err := writer.CreatePart(h)
	if err != nil {
		fmt.Printf("Error creating form file part: %v\n", err)
		os.Exit(1)
	}

	// Copy the file content to the multipart file part.
	_, err = io.Copy(filePart, file)
	if err != nil {
		fmt.Printf("Error copying file content: %v\n", err)
		os.Exit(1)
	}

	// Add the "roomId" field.
	if err := writer.WriteField("roomId", *roomIDFlag); err != nil {
		fmt.Printf("Error adding roomId field: %v\n", err)
		os.Exit(1)
	}

	// Add the "text" field if provided.
	if *textFlag != "" {
		if err := writer.WriteField("text", *textFlag); err != nil {
			fmt.Printf("Error adding text field: %v\n", err)
			os.Exit(1)
		}
	}

	// Close the writer to finalize the multipart body.
	if err := writer.Close(); err != nil {
		fmt.Printf("Error closing multipart writer: %v\n", err)
		os.Exit(1)
	}

	// Create the HTTP POST request.
	req, err := http.NewRequest("POST", "https://webexapis.com/v1/messages", &requestBody)
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		os.Exit(1)
	}

	// Set required headers.
	req.Header.Set("Authorization", "Bearer "+*tokenFlag)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and print the response.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", respBody)
}
