# Command notify_by_webex

notify_by_webex is a versatile command-line utility designed to simplify sending messages with file attachments directly to your Webex rooms via the Webex Messages API.

## Features

- **File Upload:** Send local files as attachments to a Webex room. For more details, refer to the [Message Attachments](https://developer.webex.com/docs/basics#message-attachments) section on the Webex for Developers site.
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
   git clone https://github.com/hgrimm/notify_by_webex.git
   cd notify_by_webex
   ```

2. **Build the Executable**

   Make sure you have [Go installed](https://golang.org/dl/).

   ```bash
   go mod init github.com/hgrimm/notify_by_webex
   go mod tidy
   go build -o notify_by_webex
   ```

   This command compiles the source code and creates an executable named `notify_by_webex` in the current directory.

## Usage

To list all associated rooms (displaying only the room title and ID sorted by title), run:

./notify_by_webex -T <ACCESS_TOKEN> -L

Note: The access token (-T) is required when using the -L flag.

To Send a Message with an Attachment

Run the executable with the required flags:

./notify_by_webex -T <ACCESS_TOKEN> -R <ROOM_ID> -f <FILE_PATH> [-t <MESSAGE_TEXT>]

## Flags

- `-T`: **(Required)** Your Webex access token.
- `-R`: **(Required for sending messages)** The Webex room ID where the message will be sent.
- `-f`: **(Required for sending messages)** The local file path of the file to upload.
- `-t`: **(Optional)** A plain-text message to accompany the file attachment.
- `-L`: **(Optional)** List all associated rooms (overrides message sending mode).

Example – Sending a Message:

./notify_by_webex -T "YOUR_ACCESS_TOKEN" -R "Y2lzY2....." -f "/path/to/your/file.png" -t "Here is the attached file."

This command sends the specified file to the specified Webex room with an accompanying text message.

Example – Listing Rooms:

./notify_by_webex -T "YOUR_ACCESS_TOKEN" -L

This command retrieves and displays a list of Webex rooms with their titles and IDs in a formatted table.

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

This tool leverages the [Webex Messages API](https://developer.webex.com/docs/api/v1/messages) for sending messages and uploading attachments. Special thanks to olekukonko/tablewriter for the table formatting library used in listing rooms.