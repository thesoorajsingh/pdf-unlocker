# PDF Unlocker with Go

A simple HTTP service built in Go that removes password protection from PDF files using the pdfcpu library.

## Features

- HTTP API for unlocking password-protected PDFs
- Support for user password decryption
- File validation to ensure only PDF files are processed
- Simple REST API interface
- Lightweight and fast processing

## Prerequisites

- Go 1.19 or higher
- Git

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd go-pdf-unlock
```

2. Initialize Go module and install dependencies:

```bash
go mod init go-pdf-unlock
go mod tidy
```

3. Run the service:

```bash
go run main.go
```

The service will start on port 8080.

## Usage

### API Endpoint

**POST** `/unlock`

Unlocks a password-protected PDF file.

#### Request

- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file`: The password-protected PDF file (required)
  - `password`: The password for the PDF file (required)

#### Response

- **Success (200)**: Returns the unlocked PDF file
  - **Content-Type**: `application/pdf`
  - **Content-Disposition**: `attachment; filename="unlocked.pdf"`
- **Error (4xx/5xx)**: Returns error message in plain text

#### Example using curl

```bash
curl -X POST \
  -F "file=@protected.pdf" \
  -F "password=mypassword" \
  -o unlocked.pdf \
  http://localhost:8080/unlock
```

#### Example using JavaScript (fetch)

```javascript
const formData = new FormData();
formData.append("file", pdfFile); // File object from input
formData.append("password", "mypassword");

fetch("http://localhost:8080/unlock", {
  method: "POST",
  body: formData,
})
  .then((response) => {
    if (response.ok) {
      return response.blob();
    }
    throw new Error("Failed to unlock PDF");
  })
  .then((blob) => {
    // Handle the unlocked PDF blob
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "unlocked.pdf";
    a.click();
  });
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- **400 Bad Request**:
  - Invalid request method (not POST)
  - No file provided
  - Uploaded file is not a PDF
  - Password not provided
- **405 Method Not Allowed**: Request method other than POST
- **500 Internal Server Error**:
  - Failed to read file content
  - Failed to unlock PDF (incorrect password or corrupted file)

## Dependencies

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library

## Building for Production

To build a binary for production:

```bash
go build -o pdf-unlock-service main.go
./pdf-unlock-service
```

## Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o pdf-unlock-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/pdf-unlock-service .
EXPOSE 8080
CMD ["./pdf-unlock-service"]
```

Build and run with Docker:

```bash
docker build -t pdf-unlock-service .
docker run -p 8080:8080 pdf-unlock-service
```

## Security Considerations

- This service processes uploaded files in memory
- Passwords are transmitted in form data (consider HTTPS in production)
- No authentication or rate limiting is implemented
- Files are not persisted on the server

## License

MIT License
