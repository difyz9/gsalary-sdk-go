#!/bin/bash

# 安全检查脚本 - 在 git push 前检查敏感信息

echo "🔍 开始安全检查..."
echo ""

# 颜色定义
RED='\033[0:31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ISSUES_FOUND=0

# 1. 检查 .pem 文件
echo "📝 检查密钥文件..."
PEM_FILES=$(find . -name "*.pem" -type f 2>/dev/null | grep -v "node_modules" | grep -v ".git")
if [ -n "$PEM_FILES" ]; then
    echo -e "${RED}❌ 发现 .pem 文件:${NC}"
    echo "$PEM_FILES"
    echo ""
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
else
    echo -e "${GREEN}✅ 未发现 .pem 文件${NC}"
fi

# 2. 检查 .key 文件
echo "📝 检查 .key 文件..."
KEY_FILES=$(find . -name "*.key" -type f 2>/dev/null | grep -v "node_modules" | grep -v ".git")
if [ -n "$KEY_FILES" ]; then
    echo -e "${RED}❌ 发现 .key 文件:${NC}"
    echo "$KEY_FILES"
    echo ""
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
else
    echo -e "${GREEN}✅ 未发现 .key 文件${NC}"
fi

# 3. 检查 .env 文件
echo "📝 检查环境变量文件..."
ENV_FILES=$(find . -name ".env" -o -name ".env.local" -o -name ".env.production" 2>/dev/null | grep -v "node_modules")
if [ -n "$ENV_FILES" ]; then
    echo -e "${YELLOW}⚠️  发现 .env 文件（请确保已在 .gitignore 中）:${NC}"
    echo "$ENV_FILES"
    echo ""
fi

# 4. 检查硬编码的 AppID
echo "📝 检查硬编码的 AppID..."
HARDCODED_APPID=$(grep -r "AppID.*=.*\"[0-9a-f-]\{36\}\"" --include="*.go" --exclude-dir=".git" . 2>/dev/null)
if [ -n "$HARDCODED_APPID" ]; then
    echo -e "${RED}❌ 发现硬编码的 AppID:${NC}"
    echo "$HARDCODED_APPID"
    echo ""
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
else
    echo -e "${GREEN}✅ 未发现硬编码的 AppID${NC}"
fi

# 5. 检查 Git 暂存区
echo "📝 检查 Git 暂存区..."
STAGED_PEM=$(git diff --cached --name-only | grep "\.pem$")
if [ -n "$STAGED_PEM" ]; then
    echo -e "${RED}❌ Git 暂存区包含 .pem 文件:${NC}"
    echo "$STAGED_PEM"
    echo ""
    echo "运行以下命令取消暂存:"
    echo "git reset HEAD $STAGED_PEM"
    echo ""
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
else
    echo -e "${GREEN}✅ Git 暂存区无敏感文件${NC}"
fi

# 6. 检查 .gitignore
echo "📝 检查 .gitignore 配置..."
if [ -f .gitignore ]; then
    if grep -q "\.pem" .gitignore && grep -q "\.env" .gitignore; then
        echo -e "${GREEN}✅ .gitignore 配置正确${NC}"
    else
        echo -e "${YELLOW}⚠️  .gitignore 可能缺少必要的规则${NC}"
        echo "请确保包含: *.pem, .env, .env.local"
        echo ""
    fi
else
    echo -e "${RED}❌ .gitignore 文件不存在${NC}"
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
fi

# 总结
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ 安全检查通过！可以安全提交代码。${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
else
    echo -e "${RED}❌ 发现 $ISSUES_FOUND 个安全问题！${NC}"
    echo -e "${YELLOW}请修复以上问题后再提交代码。${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 1
fi
