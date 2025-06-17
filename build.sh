#!/bin/zsh

set -e

echo "🔨 Building MGit iOS Bridge..."

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo "${GREEN}✅ $1${NC}"; }
print_warning() { echo "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo "${RED}❌ $1${NC}"; }
print_info() { echo "${BLUE}ℹ️  $1${NC}"; }

# Function to properly set up gomobile
setup_gomobile() {
    print_info "Setting up gomobile properly..."
    
    # Install the binaries
    print_info "Installing gomobile and gobind..."
    go install golang.org/x/mobile/cmd/gomobile@latest
    go install golang.org/x/mobile/cmd/gobind@latest
    
    # Initialize gomobile
    print_info "Initializing gomobile..."
    gomobile init
    
    # Check if gomobile version works
    if gomobile version &> /dev/null; then
        print_status "gomobile is working: $(gomobile version)"
        return 0
    else
        print_warning "gomobile version failed, applying your fix..."
        # This is the key step you discovered!
        print_info "Adding mobile package to module dependencies..."
        go get golang.org/x/mobile/cmd/gomobile
        
        # Test again
        if gomobile version &> /dev/null; then
            print_status "gomobile is now working: $(gomobile version)"
            return 0
        else
            print_error "gomobile still not working after fix"
            return 1
        fi
    fi
}

# Check directory
if [[ ! -f "go.mod" ]]; then
    print_error "go.mod not found. Make sure you're in the mgit-ios-bridge directory."
    exit 1
fi

# Test Go build first
print_status "Testing Go build..."
if go build .; then
    print_status "Go build successful"
else
    print_error "Go build failed - fix Go issues before building framework"
    exit 1
fi

# Check and setup gomobile
print_status "Checking gomobile setup..."

if ! command -v gomobile &> /dev/null; then
    print_warning "gomobile not found. Installing..."
    if ! setup_gomobile; then
        exit 1
    fi
elif ! gomobile version &> /dev/null; then
    print_warning "gomobile not working properly. Applying fix..."
    # Apply your discovered fix
    print_info "Adding mobile package to module dependencies..."
    go get golang.org/x/mobile/cmd/gomobile
    
    if gomobile version &> /dev/null; then
        print_status "gomobile fixed: $(gomobile version)"
    else
        print_warning "Fix didn't work, trying full setup..."
        if ! setup_gomobile; then
            exit 1
        fi
    fi
else
    print_status "gomobile is ready: $(gomobile version)"
fi

# Build framework
print_status "Building iOS framework..."
FRAMEWORK_NAME="MGitBridge.xcframework"

# Remove old framework
rm -rf "$FRAMEWORK_NAME" 2>/dev/null || true

# Build framework
if gomobile bind -target ios -o "$FRAMEWORK_NAME" .; then
    print_status "iOS framework built successfully"
else
    print_error "iOS framework build failed"
    exit 1
fi

# Verify framework
if [[ -d "$FRAMEWORK_NAME" ]]; then
    print_status "Framework created successfully"
    echo "  📁 Location: $(pwd)/$FRAMEWORK_NAME"
    echo "  📊 Size: $(du -sh $FRAMEWORK_NAME | cut -f1)"
else
    print_error "Framework was not created"
    exit 1
fi

# Copy to React Native
RN_MODULE_PATH="../react-native-mgit/ios/frameworks"
if [[ -d "$RN_MODULE_PATH" ]]; then
    print_status "Copying to React Native module..."
    rm -rf "$RN_MODULE_PATH/$FRAMEWORK_NAME"
    cp -r "$FRAMEWORK_NAME" "$RN_MODULE_PATH/"
    print_status "Framework copied to $RN_MODULE_PATH"
else
    print_warning "React Native module path not found: $RN_MODULE_PATH"
fi

print_status "Build completed successfully! 🎉"