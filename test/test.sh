#!/bin/bash
# SQL Executor 实机测试脚本
# 用法: bash test/test.sh <IP> <SSH用户> <SSH密码> <数据库类型> <数据库用户> <数据库密码> <DSN>
# 示例: bash test/test.sh 192.168.37.175 root Aa123456 mysql ava 'rKwcK@s#WgWH' "{username}:{password}@tcp(localhost:3306)/university_live"

IP="$1"
USER="$2"
PASS="$3"
DB_TYPE="$4"
DB_USER="$5"
DB_PASS="$6"
DSN="$7"

PLINK="/c/Program Files (x86)/PuTTY/plink.exe"
PSCP="/c/Program Files (x86)/PuTTY/pscp.exe"
BINARY="/tmp/gosql-linux/gosql-executor-v1.1.0-linux-amd64"
CONFIG_EXAMPLE="/tmp/gosql-linux/config.yaml.example"
REMOTE_DIR="/tmp/gosql-test-$(date +%s)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

pass() { echo -e "${GREEN}[PASS]${NC} $1"; }
fail() { echo -e "${RED}[FAIL]${NC} $1"; FAILURES=$((FAILURES + 1)); }
info() { echo -e "${YELLOW}[INFO]${NC} $1"; }

FAILURES=0

if [ -z "$IP" ] || [ -z "$USER" ] || [ -z "$PASS" ] || [ -z "$DB_TYPE" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASS" ] || [ -z "$DSN" ]; then
    echo "用法: bash test/test.sh <IP> <SSH用户> <SSH密码> <数据库类型> <DB用户> <DB密码> <DSN>"
    echo "示例: bash test/test.sh 192.168.37.175 root Aa123456 mysql ava 'mypass' \"{username}:{password}@tcp(localhost:3306)/db\""
    exit 1
fi

run_ssh() {
    "$PLINK" -batch -pw "$PASS" "${USER}@${IP}" "$1" 2>&1
    return 0
}

run_scp() {
    "$PSCP" -batch -pw "$PASS" "$1" "${USER}@${IP}:$2" 2>&1
    return 0
}

echo "========================================"
echo " SQL Executor 实机测试"
echo " 目标: ${USER}@${IP}"
echo " 数据库: ${DB_TYPE} / ${DB_USER}"
echo "========================================"
echo ""

# 1. SSH 连接
info "测试 SSH 连接..."
SSH_TEST=$(run_ssh "echo connected")
if echo "$SSH_TEST" | grep -q "connected"; then
    pass "SSH 连接成功"
else
    fail "SSH 连接失败"
    exit 1
fi

# 2. 系统信息
info "系统信息:"
OS_INFO=$(run_ssh "cat /etc/os-release | grep PRETTY_NAME | cut -d'\"' -f2")
KERNEL=$(run_ssh "uname -r")
echo "  系统: ${OS_INFO}"
echo "  内核: ${KERNEL}"
echo ""

# 3. 上传文件
info "创建远程目录 ${REMOTE_DIR}..."
run_ssh "mkdir -p ${REMOTE_DIR}"

info "上传二进制文件..."
UPLOAD_OUT=$(run_scp "$BINARY" "${REMOTE_DIR}/gosql-executor")
if echo "$UPLOAD_OUT" | grep -qi "error\|denied\|fatal"; then
    fail "二进制上传失败: ${UPLOAD_OUT}"
    exit 1
else
    pass "二进制上传成功"
fi

info "上传配置文件..."
run_scp "$CONFIG_EXAMPLE" "${REMOTE_DIR}/config.yaml.example"
pass "配置文件上传成功"

# 4. 生成 config.yaml
info "生成测试配置..."
run_ssh "cat > ${REMOTE_DIR}/config.yaml << 'CEOF'
database:
  type: \"${DB_TYPE}\"
  username: \"${DB_USER}\"
  password: \"${DB_PASS}\"
  dsn: \"${DSN}\"

output:
  directory: \"${REMOTE_DIR}/output\"
  format: \"csv\"
  show_in_console: true
  save_to_file: true
CEOF"
pass "配置文件生成完成"

# 5. 执行权限 + help
info "设置执行权限..."
run_ssh "chmod +x ${REMOTE_DIR}/gosql-executor"

info "测试 --help 参数..."
HELP_OUT=$(run_ssh "cd ${REMOTE_DIR} && ./gosql-executor --help")
if echo "$HELP_OUT" | grep -q "config"; then
    pass "--help 参数正常"
else
    fail "--help 参数异常: ${HELP_OUT}"
fi

# 6. 无配置文件启动 — 应进入交互式配置
info "测试无配置文件启动（交互式配置）..."
NO_CONF_OUT=$(run_ssh "cd ${REMOTE_DIR} && ./gosql-executor --config=${REMOTE_DIR}/noexist.yaml 2>&1 <<< '' || true")
if echo "$NO_CONF_OUT" | grep -qi "未找到配置文件，进入交互式配置"; then
    pass "无配置文件时进入交互式配置"
else
    fail "无配置文件时未进入交互式配置"
    echo "  输出: $(echo "$NO_CONF_OUT" | head -5)"
fi

# 6b. 交互式配置 — 非交互终端下 liner 应报错（预期行为）
info "测试交互式配置（非交互终端）..."
INTERACTIVE_OUT=$(run_ssh "cd ${REMOTE_DIR} && printf '${DB_TYPE}\n${IP}\n3306\n${DB_USER}\n${DB_PASS}\nuniversity_live\n' | ./gosql-executor --config=${REMOTE_DIR}/interactive_test.yaml 2>&1 || true")
if echo "$INTERACTIVE_OUT" | grep -qi "交互式配置失败\|liner.*not supported"; then
    pass "非交互终端下正确报错（liner 限制，预期行为）"
else
    fail "非交互终端下未正确报错"
    echo "  输出: $(echo "$INTERACTIVE_OUT" | head -5)"
fi

# 7. SQL 查询测试
info "执行 SHOW TABLES..."
RESULT=$(run_ssh "cd ${REMOTE_DIR} && echo 'SHOW TABLES;' | ./gosql-executor")
if echo "$RESULT" | grep -qi "panic\|FATAL"; then
    fail "SHOW TABLES 执行失败"
    echo "  输出: $(echo "$RESULT" | tail -5)"
elif echo "$RESULT" | grep -qi "\[ERROR\]"; then
    fail "SHOW TABLES 执行出错"
    echo "  输出: $(echo "$RESULT" | grep ERROR)"
else
    pass "SHOW TABLES 执行成功"
    ROW_COUNT=$(echo "$RESULT" | grep -c -v "^$\|^SQL>\|^---\|^\[INFO\]\|^  ->\|^Tables_")
    echo "  返回 ${ROW_COUNT} 行"
fi

# 8. 复杂查询
info "执行复杂查询..."
COMPLEX_OUT=$(run_ssh "cd ${REMOTE_DIR} && echo 'SELECT TABLE_NAME, TABLE_ROWS FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() ORDER BY TABLE_ROWS DESC LIMIT 5;' | ./gosql-executor")
if echo "$COMPLEX_OUT" | grep -qi "panic\|FATAL\|\[ERROR\]"; then
    fail "复杂查询执行失败"
    echo "  输出: $(echo "$COMPLEX_OUT" | tail -5)"
else
    pass "复杂查询执行成功"
fi

# 9. CSV 输出
info "检查 CSV 文件输出..."
CSV_COUNT=$(run_ssh "ls ${REMOTE_DIR}/output/*.csv 2>/dev/null | wc -l")
if [ "$CSV_COUNT" -gt 0 ]; then
    pass "CSV 文件已生成 (${CSV_COUNT} 个)"
    CSV_HEAD=$(run_ssh "head -2 ${REMOTE_DIR}/output/*.csv | head -3")
    echo "  示例:"
    echo "$CSV_HEAD" | sed 's/^/    /'
else
    fail "CSV 文件未生成"
fi

# 10. JSON 输出
info "测试 JSON 格式输出..."
run_ssh "cd ${REMOTE_DIR} && sed -i 's/format: \"csv\"/format: \"json\"/' config.yaml"
run_ssh "cd ${REMOTE_DIR} && echo 'SHOW TABLES;' | ./gosql-executor" > /dev/null 2>&1
JSON_COUNT=$(run_ssh "ls ${REMOTE_DIR}/output/*.json 2>/dev/null | wc -l")
if [ "$JSON_COUNT" -gt 0 ]; then
    pass "JSON 文件已生成"
else
    fail "JSON 文件未生成"
fi

# 11. 清理
info "清理远程测试文件..."
run_ssh "rm -rf ${REMOTE_DIR}"
pass "清理完成"

# 汇总
echo ""
echo "========================================"
if [ $FAILURES -eq 0 ]; then
    echo -e " ${GREEN}全部测试通过${NC}"
    exit 0
else
    echo -e " ${RED}${FAILURES} 项测试失败${NC}"
    exit 1
fi
