#!/bin/bash

# Detect the OS and set the file extension accordingly
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    # Windows environment (using Git Bash or WSL)
    EXE_SUFFIX=".exe"
else
    # Unix-like environment (Linux or macOS)
    EXE_SUFFIX=""
fi

# Function to build a Go module
build_module() {
    local module_dir=$1
    local output_name=$2$EXE_SUFFIX

    echo "Building $module_dir module..."
    cd "$module_dir" || exit 1
    go build -o "$output_name"
    if [[ $? -ne 0 ]]; then
        echo "Error: Failed to build $output_name."
        exit 1
    fi
    echo "$output_name built successfully."
    cd ..
}

# Build the 'bank' module as 'bankGo'
build_module "bank" "bankGo"

# Build the 'mergesort' module as 'mergesortGo'
build_module "mergesort" "mergesortGo"

echo "Build process completed for all modules."
