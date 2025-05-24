# Swiss Linux Knife - Complete Setup Guide

## Overview
Swiss Linux Knife is a GUI application built with Go and Fyne that provides visual tools for managing Linux system configurations, including shell config management and URL shortening capabilities.

## Prerequisites

### System Requirements
- Ubuntu/Debian-based Linux distribution (Ubuntu 20.04+, Debian 10+, Linux Mint, etc.)
- X11 display server (Wayland users may need XWayland)
- Go 1.24.0 or later
- Git

### Required Packages

## Step-by-Step Installation Guide

### Step 1: Install X11 Dependencies

The Fyne framework requires several X11 and graphics libraries to render the GUI properly.

```bash
# Update package list
sudo apt update

# Install essential X11 development libraries
sudo apt install -y \
    xorg-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxext-dev \
    libxfixes-dev \
    libxrender-dev \
    libgl1-mesa-dev \
    libglu1-mesa-dev

# Install additional graphics libraries
sudo apt install -y \
    libglfw3-dev \
    libgles2-mesa-dev \
    libegl1-mesa-dev

# Install font and image processing libraries
sudo apt install -y \
    libfreetype6-dev \
    libfontconfig1-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev
```

### Step 2: Install Go (if not already installed)

```bash
# Download and install Go 1.24.0
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Remove any existing Go installation
sudo rm -rf /usr/local/go

# Extract and install Go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add Go to PATH (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# Reload shell configuration
source ~/.bashrc

# Verify installation
go version
```

### Step 3: Clone and Build Swiss Linux Knife

```bash
# Clone the repository
cd ~/code/personal/tools  # or your preferred directory
git clone <repository-url> swiss-linux-knife
cd swiss-linux-knife

# Download dependencies
go mod download

# Build the application
go build -o swiss-linux-knife

# Make it executable (if needed)
chmod +x swiss-linux-knife
```

### Step 4: Running the Application

```bash
# Run from the build directory
./swiss-linux-knife

# Or install globally
sudo cp swiss-linux-knife /usr/local/bin/

# Then run from anywhere
swiss-linux-knife
```

## Troubleshooting Common Issues

### Issue 1: "cannot open display" Error

If you encounter display errors, ensure X11 is properly configured:

```bash
# Check if DISPLAY variable is set
echo $DISPLAY

# If empty, set it manually
export DISPLAY=:0

# For SSH connections, enable X11 forwarding
ssh -X username@hostname
```

### Issue 2: Missing OpenGL Support

```bash
# Install Mesa drivers
sudo apt install -y mesa-utils

# Test OpenGL
glxinfo | grep "OpenGL version"
```

### Issue 3: Font Rendering Issues

```bash
# Install additional font packages
sudo apt install -y \
    fonts-liberation \
    fonts-dejavu-core \
    fontconfig
```

### Issue 4: Wayland Compatibility

For Wayland users:

```bash
# Install XWayland
sudo apt install -y xwayland

# Run with X11 backend
GDK_BACKEND=x11 ./swiss-linux-knife
```

### Issue 5: Build Errors

If you encounter build errors:

```bash
# Clean module cache
go clean -modcache

# Update dependencies
go mod tidy

# Rebuild with verbose output
go build -v -o swiss-linux-knife
```

## Features Overview

### 1. Shell Config Manager
- Visual editor for .bashrc/.zshrc files
- Manage environment variables
- Configure aliases
- Oh My Zsh theme and plugin management
- Custom function editor

### 2. URL Shortener
- Create and manage short URLs (feature in development)

## System Integration

### Creating a Desktop Entry

Create a desktop launcher for easy access:

```bash
# Create desktop entry
cat > ~/.local/share/applications/swiss-linux-knife.desktop << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Swiss Linux Knife
Comment=Linux system configuration tools
Icon=utilities-terminal
Exec=/usr/local/bin/swiss-linux-knife
Terminal=false
Categories=System;Settings;
EOF

# Update desktop database
update-desktop-database ~/.local/share/applications/
```

### Setting up as a System Service (Optional)

For background URL shortener service:

```bash
# Create systemd service file
sudo cat > /etc/systemd/system/swiss-knife-url.service << EOF
[Unit]
Description=Swiss Linux Knife URL Shortener Service
After=network.target

[Service]
Type=simple
User=$USER
ExecStart=/usr/local/bin/swiss-linux-knife --url-service
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl enable swiss-knife-url.service
sudo systemctl start swiss-knife-url.service
```

## Performance Optimization

### For Better Performance

1. **Enable Hardware Acceleration**:
```bash
# Check if hardware acceleration is available
glxinfo | grep "direct rendering"
```

2. **Adjust Fyne Settings**:
```bash
# Set environment variables for better performance
export FYNE_SCALE=1.0
export FYNE_THEME=light
```

## Security Considerations

1. The application modifies shell configuration files - always backup before use:
```bash
cp ~/.zshrc ~/.zshrc.backup
cp ~/.bashrc ~/.bashrc.backup
```

2. Review any changes before saving
3. The application requires read/write access to your home directory

## Additional Resources

- [Fyne Documentation](https://docs.fyne.io/)
- [Go Documentation](https://golang.org/doc/)
- [X11 Documentation](https://www.x.org/wiki/)

## Support

For issues specific to:
- **X11/Display problems**: Check your distribution's documentation
- **Go build issues**: Refer to Go's official documentation
- **Application bugs**: Submit issues to the project repository

---

Last updated: 2025-05-24