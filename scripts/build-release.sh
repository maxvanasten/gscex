#!/bin/bash

# Build and Release Script for gscex
# Usage: ./scripts/build-release.sh [version]
# Example: ./scripts/build-release.sh v1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get version from argument or use default
VERSION=${1:-"dev"}
BINARY_NAME="gscex"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Building gscex ${VERSION}${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Create dist directory
DIST_DIR="dist"
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

# Build platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

echo -e "${YELLOW}Building for multiple platforms...${NC}"
echo ""

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=$(echo ${PLATFORM} | cut -d'/' -f1)
    GOARCH=$(echo ${PLATFORM} | cut -d'/' -f2)
    
    OUTPUT_NAME="${BINARY_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo -e "  Building ${GOOS}/${GOARCH}..."
    
    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags "-s -w -X main.version=${VERSION}" \
        -o "${DIST_DIR}/${OUTPUT_NAME}" \
        ./cmd/gscex/
    
    if [ $? -eq 0 ]; then
        echo -e "    ${GREEN}✓${NC} ${OUTPUT_NAME}"
    else
        echo -e "    ${RED}✗${NC} Failed to build ${OUTPUT_NAME}"
        exit 1
    fi
done

echo ""
echo -e "${YELLOW}Creating archives...${NC}"
echo ""

cd "${DIST_DIR}"

# Create tar.gz for Unix systems and zip for Windows
for file in ${BINARY_NAME}-*; do
    if [[ "$file" == *.exe ]]; then
        # Windows - create zip
        zip "${file%.exe}.zip" "$file"
        rm "$file"
        echo -e "  ${GREEN}✓${NC} ${file%.exe}.zip"
    else
        # Unix - create tar.gz
        tar -czf "${file}.tar.gz" "$file"
        rm "$file"
        echo -e "  ${GREEN}✓${NC} ${file}.tar.gz"
    fi
done

cd ..

echo ""
echo -e "${YELLOW}Generating checksums...${NC}"
cd "${DIST_DIR}"
sha256sum *.tar.gz *.zip > checksums.txt 2>/dev/null || shasum -a 256 *.tar.gz *.zip > checksums.txt
cd ..
echo -e "  ${GREEN}✓${NC} checksums.txt"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Build Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Artifacts in ${DIST_DIR}/${NC}:"
ls -lh "${DIST_DIR}"
echo ""

# Check if we should create a GitHub release
if [ "$VERSION" != "dev" ] && command -v gh &> /dev/null; then
    echo -e "${YELLOW}GitHub CLI detected. Create release?${NC}"
    read -p "Create GitHub release ${VERSION}? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Creating GitHub release...${NC}"
        
        # Create release with notes
        gh release create "${VERSION}" \
            --title "gscex ${VERSION}" \
            --notes "Release ${VERSION}" \
            ${DIST_DIR}/*
        
        if [ $? -eq 0 ]; then
            echo ""
            echo -e "${GREEN}✓ GitHub release ${VERSION} created successfully!${NC}"
            echo ""
            echo -e "${BLUE}Release URL:${NC}"
            gh release view "${VERSION}" --json url -q .url
        else
            echo ""
            echo -e "${RED}✗ Failed to create GitHub release${NC}"
            exit 1
        fi
    fi
else
    if [ "$VERSION" != "dev" ]; then
        echo -e "${YELLOW}To create a GitHub release, install GitHub CLI:${NC}"
        echo -e "  https://cli.github.com/"
        echo ""
        echo -e "${BLUE}Manual release instructions:${NC}"
        echo -e "1. Go to: https://github.com/yourusername/gscex/releases/new"
        echo -e "2. Tag version: ${VERSION}"
        echo -e "3. Upload files from ${DIST_DIR}/"
    fi
fi

echo ""
echo -e "${GREEN}Done!${NC}"
