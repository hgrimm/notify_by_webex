# notify_by_webex

notify_by_webex is a versatile command-line utility designed to simplify sending messages with file attachments directly to your Webex rooms via the Webex Messages API.

## Features

- **File Upload:** Send local files as attachments to a Webex room.
- **Command-Line Flags:** Provide the Webex authorization token, room ID, message text, and file path through convenient flags.
- **Validation:** Ensures the file is one of the supported types and does not exceed a maximum size of 100 MB.
- **Logging:** Outputs the `Content-Disposition` and `Content-Type` headers for verification and debugging.

## Supported File Types

The tool supports the following file types:
- Microsoft Word: `.doc`, `.docx`
- Microsoft Excel: `.xls`, `.xlsx`
- Microsoft PowerPoint: `.ppt`, `.pptx`
- Adobe PDF: `.pdf`
- Images: `.jpg`, `.jpeg`, `.bmp`, `.gif`, `.png`

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-username/notify_by_webex.git
   cd notify_by_webex
