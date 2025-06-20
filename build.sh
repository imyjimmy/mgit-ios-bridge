#!/bin/zsh

set -e

echo "🔨 Building MGit Bridge..."

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo "${GREEN}✅ $1${NC}"; }
print_warning() { echo "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo "${RED}❌ $1${NC}"; }
print_info() { echo "${BLUE}ℹ️  $1${NC}"; }

# iOS build function
build_ios() {
    print_status "Building iOS framework..."
    FRAMEWORK_NAME="MGitBridge.xcframework"
    
    # Remove old framework
    rm -rf "$FRAMEWORK_NAME" 2>/dev/null || true
    
    # Build iOS framework
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
    
    # Copy to React Native iOS
    RN_MODULE_PATH="../react-native-mgit/ios/frameworks"
    if [[ -d "$RN_MODULE_PATH" ]]; then
        print_status "Copying to React Native iOS module..."
        rm -rf "$RN_MODULE_PATH/$FRAMEWORK_NAME"
        cp -r "$FRAMEWORK_NAME" "$RN_MODULE_PATH/"
        print_status "Framework copied to $RN_MODULE_PATH"
    else
        print_warning "React Native iOS module path not found: $RN_MODULE_PATH"
    fi
}

# Android build function  
build_android() {
    print_status "Building Android library..."
    LIBRARY_NAME="mgitbridge.aar"
    
    # Remove old library
    rm -rf "$LIBRARY_NAME" 2>/dev/null || true
    
    # Build Android library
    if gomobile bind -target android -o "$LIBRARY_NAME" .; then
        print_status "Android library built successfully"
    else
        print_error "Android library build failed"
        exit 1
    fi
    
    # Verify library
    if [[ -f "$LIBRARY_NAME" ]]; then
        print_status "Library created successfully"
        echo "  📁 Location: $(pwd)/$LIBRARY_NAME"
        echo "  📊 Size: $(du -sh $LIBRARY_NAME | cut -f1)"
    else
        print_error "Library was not created"
        exit 1
    fi
    
    # Copy to React Native Android (future)
    RN_ANDROID_PATH="../react-native-mgit/android/libs"
    if [[ -d "$RN_ANDROID_PATH" ]]; then
        print_status "Copying to React Native Android module..."
        cp "$LIBRARY_NAME" "$RN_ANDROID_PATH/"
        print_status "Library copied to $RN_ANDROID_PATH"
    else
        print_warning "React Native Android module path not found: $RN_ANDROID_PATH"
        print_info "Will create when React Native Android support is added"
    fi
}

# Parse command line arguments
TARGET="ios"  # Default target
while [[ $# -gt 0 ]]; do
    case $1 in
        --target)
            TARGET="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [--target <ios|android>]"
            echo ""
            echo "Options:"
            echo "  --target    Target platform (ios or android) [default: ios]"
            echo "  -h, --help  Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Validate target
case $TARGET in
    ios|android)
        print_info "Building for target: $TARGET"
        ;;
    *)
        print_error "Invalid target: $TARGET. Supported targets: ios, android"
        exit 1
        ;;
esac

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
    print_error "go.mod not found. Make sure you're in the mgit-bridge directory."
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

# Build based on target
case $TARGET in
    ios)
        build_ios
        ;;
    android)
        build_android
        ;;
esac

print_status "Build completed successfully for $TARGET! 🎉"