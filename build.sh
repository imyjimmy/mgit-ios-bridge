#!/bin/zsh

# MGit iOS Bridge Build Script
set -e

echo "🔨 Building MGit iOS Bridge..."

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() { echo "${GREEN}✅ $1${NC}"; }
print_warning() { echo "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo "${RED}❌ $1${NC}"; }

# Check directory
if [[ ! -f "go.mod" ]]; then
    print_error "go.mod not found. Make sure you're in the mgit-ios-bridge directory."
    exit 1
fi

# Step 1: Clean up dependencies
print_status "Cleaning up Go modules..."
go mod tidy

# Step 2: Test Go build
print_status "Testing Go build..."
if go build .; then
    print_status "Go build successful"
else
    print_error "Go build failed"
    exit 1
fi

# Step 3: Check and fix gomobile installation
print_status "Checking gomobile installation..."

if ! command -v gomobile &> /dev/null; then
    print_warning "gomobile not found. Installing..."
    go install golang.org/x/mobile/cmd/gomobile@latest
    go install golang.org/x/mobile/cmd/gobind@latest
fi

# Check if gomobile init is needed
if ! gomobile version &> /dev/null; then
    print_warning "gomobile not initialized. Running gomobile init..."
    if ! gomobile init; then
        print_error "gomobile init failed. Trying alternative approach..."
        
        # Alternative: reinstall everything
        print_warning "Reinstalling gomobile tools..."
        go install golang.org/x/mobile/cmd/gomobile@latest
        go install golang.org/x/mobile/cmd/gobind@latest
        
        # Try init again
        if ! gomobile init; then
            print_error "Failed to initialize gomobile. Please run manually:"
            echo "  go install golang.org/x/mobile/cmd/gomobile@latest"
            echo "  go install golang.org/x/mobile/cmd/gobind@latest"
            echo "  gomobile init"
            exit 1
        fi
    fi
fi

# Verify gomobile is working
print_status "Verifying gomobile..."
if gomobile version; then
    print_status "gomobile is ready"
else
    print_error "gomobile verification failed"
    exit 1
fi

# Step 4: Build iOS framework
print_status "Building iOS framework..."
if gomobile bind -target ios -o MGitBridge.xcframework .; then
    print_status "iOS framework built successfully"
else
    print_error "iOS framework build failed"
    exit 1
fi

# Step 5: Copy to React Native module
RN_MODULE_PATH="../react-native-mgit/ios/frameworks"
if [[ -d "$RN_MODULE_PATH" ]]; then
    print_status "Copying framework to React Native module..."
    rm -rf "$RN_MODULE_PATH/MGitBridge.xcframework"
    cp -r MGitBridge.xcframework "$RN_MODULE_PATH/"
    print_status "Framework copied to $RN_MODULE_PATH"
else
    print_warning "React Native module path not found: $RN_MODULE_PATH"
fi

# Step 6: Show results
if [[ -d "MGitBridge.xcframework" ]]; then
    print_status "Framework details:"
    echo "  📁 Location: $(pwd)/MGitBridge.xcframework"
    echo "  📊 Size: $(du -sh MGitBridge.xcframework | cut -f1)"
fi

print_status "Build completed successfully! 🎉"
