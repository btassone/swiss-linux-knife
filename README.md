# Swiss Linux Knife

A collection of GUI tools for Linux system management and configuration.

## Project Structure

```
swiss-linux-knife/
├── cmd/
│   └── swiss-linux-knife/     # Application entry point
│       └── main.go
├── internal/                   # Private application code
│   ├── gui/                   # GUI components
│   │   └── shellconfig.go     # Shell configuration GUI
│   ├── shellconfig/           # Shell configuration logic
│   │   └── config.go          # Config parsing and management
│   └── tools/                 # Tool registry
│       └── registry.go        # Available tools registration
├── pkg/                       # Public libraries (future use)
├── bin/                       # Build output directory
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── Makefile                   # Build automation
└── README.md                  # This file
```

## Features

### Shell Config Manager
- Visual editor for .bashrc/.zshrc configuration
- Environment variable management
- PATH editor with drag-and-drop ordering
- Alias management
- Oh My Zsh theme and plugin configuration
- Custom function editor
- Shell history viewer

## Building

```bash
# Build the application
make build

# Run the application
make run

# Clean build artifacts
make clean

# Run tests
make test

# Format code
make fmt
```

## Installation

```bash
# Clone the repository
git clone https://github.com/btassone/swiss-linux-knife.git
cd swiss-linux-knife

# Install dependencies
make deps

# Build the application
make build

# Run the application
./bin/swiss-linux-knife
```

## Development

The project follows Go best practices with a clear separation of concerns:

- `cmd/`: Contains the application entry points
- `internal/`: Private application code that cannot be imported by other projects
- `pkg/`: Public libraries that can be imported (future use)

### Adding New Tools

1. Create a new GUI component in `internal/gui/`
2. Register the tool in `internal/tools/registry.go`
3. The tool will automatically appear in the main window

## Requirements

- Go 1.24 or higher
- Fyne framework dependencies
- Linux operating system (for full functionality)

## License

[Add your license here]