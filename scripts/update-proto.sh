#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

PROTO_MODULE="github.com/legendaryneobatman/shop-proto-repo"

echo -e "${YELLOW}🔍 Checking for updates...${NC}"

# Используем go list для получения последней версии
# Go сам знает как правильно сортировать семантические версии
LATEST_VERSION=$(go list -m -versions ${PROTO_MODULE} 2>/dev/null | awk '{print $NF}')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}❌ Could not fetch latest version${NC}"
    echo -e "${YELLOW}Make sure the module is accessible${NC}"
    exit 1
fi

# Текущая версия
CURRENT_VERSION=$(go list -m -f '{{.Version}}' ${PROTO_MODULE} 2>/dev/null || echo "not installed")

echo -e "${YELLOW}Current: ${CURRENT_VERSION}${NC}"
echo -e "${GREEN}Latest:  ${LATEST_VERSION}${NC}"

if [ "$CURRENT_VERSION" = "$LATEST_VERSION" ]; then
    echo -e "${GREEN}✅ Already up to date!${NC}"
    exit 0
fi

# Обновляем
echo -e "${GREEN}📦 Updating...${NC}"
go get ${PROTO_MODULE}@${LATEST_VERSION}
go mod tidy

# Опционально: автокоммит
if [ "$1" = "--commit" ] || [ "$1" = "-c" ]; then
    echo -e "${GREEN}📝 Committing changes...${NC}"
    git add go.mod go.sum
    git commit -m "chore: update proto contracts to ${LATEST_VERSION}"
    echo -e "${GREEN}✨ Committed!${NC}"
else
    echo -e "${YELLOW}💡 Run with --commit flag to auto-commit changes${NC}"
fi

echo -e "${GREEN}✨ Updated from ${CURRENT_VERSION} to ${LATEST_VERSION}${NC}"