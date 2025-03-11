# Go-cli-app
A command-line application written in Go that reads URLs from a CSV file, downloads their content concurrently, and saves the results to disk.

## Features

- Reads URLs from a CSV file.
- Concurrently downloads content using a worker pool.
- Saves downloaded content as files.
- Logs metrics such as processed URLs, success/failure rates, and download duration.
- Fully tested with unit tests covering edge cases.

## Installation

### Prerequisites

- Go 1.18 or later installed
- Git installed

### Clone the Repository

```sh
git clone https://github.com/2231sshubham/Go-cli-app.git
cd Go-cli-app
```

### Build

```sh
go build -o go-app main.go
```

## Usage

Run the application by providing a CSV file containing URLs:

```sh
./go-app urls.csv output_directory
```

## Running Tests

Run unit tests with:

```sh
go test ./tests/...
```

## License

This project is licensed under the MIT License.

---

Contributions and issues are welcome! 🚀
