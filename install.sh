#!/bin/bash

installfractal() {
    echo "Installing Fractal..."

    # Clone the repository containing docker-compose.yml
    git clone https://github.com/SkySingh04/fractal.git /tmp/fractal || {
        echo "Failed to clone the Fractal repository."
        exit 1
    }

    cd /tmp/fractal || exit

    echo "Running Docker Compose..."

    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        echo "Docker Compose not found. Please install Docker Compose and try again."
        exit 1
    fi

    # Start the app using Docker Compose
    docker-compose up -d

    # Create an alias for starting the app with Docker Compose
    alias_cmd="alias fractal='cd /tmp/fractal && docker-compose up -d'"
    
    # Add alias to shell configuration file based on the user's shell
    current_shell="$(basename "$SHELL")"
    if [[ "$current_shell" == "zsh" ]]; then
        if ! grep -q "alias fractal=" "$HOME/.zshrc"; then
            echo "$alias_cmd" >> "$HOME/.zshrc"
        fi
        source "$HOME/.zshrc"
    elif [[ "$current_shell" == "bash" ]]; then
        if ! grep -q "alias fractal=" "$HOME/.bashrc"; then
            echo "$alias_cmd" >> "$HOME/.bashrc"
        fi
        source "$HOME/.bashrc"
    else
        if ! grep -q "alias fractal=" "$HOME/.profile"; then
            echo "$alias_cmd" >> "$HOME/.profile"
        fi
        source "$HOME/.profile"
    fi

    echo "Alias 'fractal' added! Use 'fractal' to run the app."
}

installfractal