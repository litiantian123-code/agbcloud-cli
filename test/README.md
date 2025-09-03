# 测试文档

本项目采用分层测试策略，将测试分为单元测试和集成测试两个层次。

## 测试结构

```
test/
├── unit/           # 单元测试 - 不依赖外部环境
│   └── client_test.go
├── integration/    # 集成测试 - 需要外部API环境
│   ├── api_integration_test.go      # 带构建标签的集成测试
│   └── api_integration_ide_test.go  # IDE友好的集成测试
├── README.md       # 本文档
└── IDE_TESTING_GUIDE.md  # IDE测试运行指南
```

## 测试类型

### 单元测试 (Unit Tests)

- **位置**: `test/unit/`
- **特点**: 不依赖外部环境，使用mock服务器
- **覆盖范围**: 
  - 客户端创建和配置
  - API调用逻辑（使用httptest）
  - 错误处理
  - 认证机制
  - 配置管理

### 集成测试 (Integration Tests)

- **位置**: `test/integration/`
- **特点**: 测试与真实API的交互
- **文件类型**:
  - `oauth_integration_test.go` - 带构建标签 `//go:build integration`
- **覆盖范围**:
  - 完整的API请求流程
  - 网络错误处理
  - OAuth认证流程
  - 配置集成

## 运行测试

> **IDE用户注意**: 如果在IDE中运行集成测试遇到构建标签问题，请参考 [IDE测试运行指南](IDE_TESTING_GUIDE.md)

### 使用Make命令

```bash
# 运行单元测试（默认）
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 运行所有测试
make test-all

# 运行测试并生成覆盖率报告
make test-coverage

# 详细模式运行测试
make test-verbose
```

### 使用测试脚本

```bash
# 运行单元测试（默认）
./scripts/test.sh

# 运行单元测试
./scripts/test.sh --unit-only

# 运行集成测试
./scripts/test.sh --integration-only

# 运行所有测试
./scripts/test.sh --all

# 详细模式
./scripts/test.sh --verbose

# 查看帮助
./scripts/test.sh --help
```

### 直接使用Go命令

```bash
# 运行单元测试
go test ./test/unit/...

# 运行集成测试（带构建标签）
go test -tags=integration ./test/integration/...

# 运行IDE版本的集成测试
go test ./test/integration/... -run ".*IDE.*"

# 运行所有测试（包括internal包）
go test ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./test/unit/... ./internal/...
go tool cover -html=coverage.out -o coverage.html
```

## 环境变量

### 集成测试环境变量

| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `AGB_API_KEY` | AgbCloud API密钥 | - | 否 |
| `SKIP_INTEGRATION_TESTS` | 跳过集成测试 | false | 否 |

### 设置示例

```bash
# 设置API密钥（可选）
export AGB_API_KEY="your-api-key-here"

# 跳过集成测试
export SKIP_INTEGRATION_TESTS="true"

# 运行集成测试
make test-integration
```

## 测试原则

### 单元测试原则

1. **隔离性**: 每个测试独立运行，不依赖其他测试
2. **可重复性**: 测试结果应该是确定的和可重复的
3. **快速性**: 单元测试应该快速执行
4. **无外部依赖**: 使用mock和stub替代外部服务

### 集成测试原则

1. **真实环境**: 测试真实的API交互
2. **容错性**: 能够处理网络错误和API变更
3. **可配置**: 通过环境变量控制测试行为
4. **文档化**: 清楚说明测试的预期行为

## 测试覆盖率

项目目标是保持高测试覆盖率：

- **单元测试覆盖率**: 目标 > 80%
- **关键路径覆盖**: 100%
- **错误处理覆盖**: 100%

查看覆盖率报告：

```bash
make test-coverage
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

## 添加新测试

### 添加单元测试

1. 在 `test/unit/` 目录下创建测试文件
2. 使用 `httptest` 创建mock服务器
3. 测试所有成功和失败场景
4. 确保测试不依赖外部环境

### 添加集成测试

1. 在 `test/integration/` 目录下创建测试文件
2. 添加构建标签 `//go:build integration`
3. 使用环境变量控制测试行为
4. 处理网络错误和API不可用情况

## 持续集成

在CI/CD流水线中：

```yaml
# 示例GitHub Actions配置
- name: Run Unit Tests
  run: make test-unit

- name: Run Integration Tests
  run: make test-integration
  env:
    AGB_API_KEY: ${{ secrets.AGB_API_KEY }}
    SKIP_INTEGRATION_TESTS: "false"
```

## 故障排除

### 常见问题

1. **集成测试失败**
   - 检查网络连接
   - 验证API密钥设置
   - 查看API服务状态

2. **单元测试失败**
   - 检查mock服务器配置
   - 验证测试数据
   - 查看错误日志

3. **覆盖率不足**
   - 添加边界条件测试
   - 测试错误处理路径
   - 增加负面测试用例

### 调试技巧

```bash
# 运行特定测试
go test -v ./test/unit/ -run TestSpecificFunction

# 启用详细输出
go test -v ./test/unit/...

# 运行测试并显示覆盖率
go test -v -cover ./test/unit/...
``` 