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
   ```

2. **Build the Executable**

   Make sure you have [Go installed](https://golang.org/dl/).

   ```bash
   go build -o notify_by_webex
   ```

   This command compiles the source code and creates an executable named `notify_by_webex` in the current directory.

## Usage

Run the executable with the required flags:

```bash
./notify_by_webex -T <ACCESS_TOKEN> -R <ROOM_ID> -f <FILE_PATH> [-t <MESSAGE_TEXT>]
```

### Flags

- `-T`: **(Required)** Your Webex access token.
- `-R`: **(Required)** The Webex room ID where the message will be sent.
- `-f`: **(Required)** The local file path of the file to upload.
- `-t`: *(Optional)* A plain-text message to accompany the file attachment.

### Example

```bash
./notify_by_webex -T "YOUR_ACCESS_TOKEN" -R "Y2lzY2....." -f "/home/desktop/example.png" -t "Example message with attachment"
```

This command sends the `example.png` file to the specified Webex room with an accompanying text message.

## Logging

For debugging purposes, the tool logs the following HTTP header values used for the file part of the multipart upload:
- **Content-Disposition:** Provides the file name and form-data name.
- **Content-Type:** Specifies the MIME type of the file.

These logs are printed to the console prior to sending the HTTP request.

## Error Handling

The tool will exit with an error message in the following cases:
- Missing required flags.
- File not found or inaccessible.
- File exceeds 100 MB in size.
- File type is not supported.
- Errors during file reading or HTTP request setup.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests for new features, bug fixes, or enhancements.

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgements

This tool leverages the [Webex Messages API](https://developer.webex.com/docs/api/v1/messages) for sending messages and uploading attachments.

