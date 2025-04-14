package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const maxFileSize = 100 * 1024 * 1024 // 100 MB


var verboseLevel int

func logDebug(level int, msg string, args ...interface{}) {
	if verboseLevel >= level {
		log.Printf(msg, args...)
	}
}

type Room struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type RoomsResponse struct {
	Items []Room `json:"items"`
}

func main() {
	token := flag.String("T", "", "Webex API Token (required)")
	roomID := flag.String("R", "", "Webex Room ID")
	recipient := flag.String("r", "", "Recipient's email address (alternative to Room ID)")
	text := flag.String("t", "", "Message text (optional)")
	markdown := flag.String("m", "", "Markdown message (optional)")
	filePath := flag.String("f", "", "Path to local file (optional)")
	fileURL := flag.String("F", "", "Public file URL to attach (optional)")
	cardFile := flag.String("A", "", "Path to Adaptive Card JSON file (optional)")
	listRoomsFlag := flag.Bool("L", false, "List available rooms")
	verbosity := flag.Int("v", 0, "Verbosity level (0 ... 2)")
	flag.Parse()
	verboseLevel = *verbosity

	if *token == "" {
		fmt.Fprintln(os.Stderr, "Fehler: Token (-T) ist erforderlich.")
		os.Exit(1)
	}

	if *listRoomsFlag {
		listRooms(*token)
		return
	}

	if *roomID == "" && *recipient == "" {
		fmt.Fprintln(os.Stderr, "Fehler: Entweder -R (Room ID) oder -r (Empfängeradresse) muss angegeben werden.")
		os.Exit(1)
	}

	endpointUrl := "https://webexapis.com/v1/messages"
	var req *http.Request
	var err error

	if *filePath != "" {
		// multipart/form-data für Dateiupload
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		// Datei öffnen
		fInfo, err := os.Stat(*filePath)
		if err != nil || fInfo.Size() > maxFileSize {
			log.Fatalf("Dateifehler: %v", err)
		}
		file, _ := os.Open(*filePath)
		defer file.Close()

		mimeType := mime.TypeByExtension(filepath.Ext(*filePath))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files"; filename="%s"`, filepath.Base(*filePath)))
		h.Set("Content-Type", mimeType)
		part, _ := writer.CreatePart(h)
		io.Copy(part, file)

		// Pflichtfelder
		if *roomID != "" {
			writer.WriteField("roomId", *roomID)
		} else {
			writer.WriteField("toPersonEmail", *recipient)
		}

		// Optional: Text, Markdown
		if *text != "" {
			writer.WriteField("text", *text)
		}
		if *markdown != "" {
			writer.WriteField("markdown", *markdown)
		}
		// Optional: öffentlicher Datei-Link
		if *fileURL != "" {
			writer.WriteField("files", *fileURL)
		}
		writer.Close()

		req, err = http.NewRequest("POST", endpointUrl, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else if *cardFile != "" {
		// JSON mit Adaptive Card direkt senden
		cardData, err := os.ReadFile(*cardFile)
		if err != nil {
			log.Fatalf("Fehler beim Lesen der Card-Datei: %v", err)
		}
		recipientField := "roomId"
		recipientValue := *roomID
		if *recipient != "" {
			recipientField = "toPersonEmail"
			recipientValue = *recipient
		}
		jsonPayload := fmt.Sprintf(`{
			"%s": "%s",
			"text": "%s",
			"attachments": [{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": %s
			}]
		}`, recipientField, recipientValue, *text, string(cardData))

		req, err = http.NewRequest("POST", endpointUrl, strings.NewReader(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
	} else {
		// Nur Textnachricht (kein File, keine Card)
		values := make(url.Values)
		if *roomID != "" {
			values.Set("roomId", *roomID)
		} else {
			values.Set("toPersonEmail", *recipient)
		}
		if *text != "" {
			values.Set("text", *text)
		}
		if *markdown != "" {
			values.Set("markdown", *markdown)
		}
		if *fileURL != "" {
			values.Set("files", *fileURL)
		}
		req, err = http.NewRequest("POST", endpointUrl, strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		log.Fatalf("Request-Erstellung fehlgeschlagen: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+*token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Fehler beim Senden: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}

func listRooms(token string) {
	req, err := http.NewRequest("GET", "https://webexapis.com/v1/rooms", nil)
	if err != nil {
		log.Fatalf("Fehler bei Request-Erstellung: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fehler beim Abrufen: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var r RoomsResponse
	if err := json.Unmarshal(body, &r); err != nil {
		log.Fatalf("Fehler beim Parsen: %v", err)
	}

	sort.Slice(r.Items, func(i, j int) bool {
		return r.Items[i].Title < r.Items[j].Title
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "ID"})
	for _, room := range r.Items {
		table.Append([]string{room.Title, room.ID})
	}
	table.Render()
}
