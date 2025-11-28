#!/bin/bash

# 清理敏感信息脚本
# 此脚本将从 Git 仓库中移除敏感文件

echo "🧹 开始清理敏感信息..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 警告
echo -e "${RED}⚠️  警告：此操作将修改 Git 历史记录！${NC}"
echo "   建议在独立分支上执行此操作。"
echo ""
read -p "确认继续？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo "步骤 1: 从 Git 索引中移除 .pem 文件（保留本地文件）..."
git rm --cached *.pem 2>/dev/null
git rm --cached **/*.pem 2>/dev/null

echo ""
echo "步骤 2: 提交更改..."
git add .gitignore
git commit -m "chore: remove sensitive files from git tracking

- Remove all .pem files from git tracking
- Update .gitignore to prevent future commits
- Add .env.example for configuration template
- Add SECURITY.md documentation"

echo ""
echo -e "${GREEN}✅ 清理完成！${NC}"
echo ""
echo "📝 后续步骤："
echo "1. 本地 .pem 文件已保留，但不再被 Git 跟踪"
echo "2. 确认 .gitignore 包含 *.pem 和 .env"
echo "3. 复制 .env.example 为 .env 并填入实际配置"
echo "4. 运行: ./scripts/security-check.sh 进行验证"
echo "5. 推送到远程仓库: git push origin main"
echo ""
echo -e "${YELLOW}⚠️  注意：如果需要清理 Git 历史记录，请使用以下命令：${NC}"
echo "   git filter-branch --force --index-filter \\"
echo "     'git rm --cached --ignore-unmatch *.pem' \\"
echo "     --prune-empty --tag-name-filter cat -- --all"
echo ""
echo "   或使用更安全的工具："
echo "   pip install git-filter-repo"
echo "   git filter-repo --invert-paths --path-regex '.*\.pem$'"
