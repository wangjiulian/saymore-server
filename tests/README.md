# 测试用例说明

本目录包含了 SayMore 服务器项目的测试用例。测试用例分为多个层次，包括模型测试、控制器测试和集成测试。

## 目录结构

```
tests/
├── models/           # 模型测试
├── controllers/      # 控制器测试
├── repository/       # 仓库模拟实现
└── integration/      # 集成测试
```

## 运行测试

### 运行所有测试

```bash
go test ./tests/...
```

### 运行特定目录的测试

```bash
go test ./tests/models/     # 运行模型测试
go test ./tests/controllers/ # 运行控制器测试
```

### 运行特定测试文件

```bash
go test ./tests/models/student_test.go
```

### 运行特定测试函数

```bash
go test ./tests/controllers -run TestGetStudentProfile
```

## 测试覆盖率

生成测试覆盖率报告：

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 注意事项

1. 集成测试需要配置测试数据库，默认情况下会被跳过。要运行集成测试，请移除测试函数中的 `t.Skip()` 语句并配置测试数据库。

2. 控制器测试默认不会与实际数据库交互，它们主要测试路由和参数处理逻辑。要进行完整测试，请使用 `tests/repository` 中的模拟仓库。

3. 模型测试主要验证模型结构和字段定义的正确性。

## 添加新测试

添加新测试时，请遵循以下规则：

1. 为每个新模型创建对应的模型测试文件
2. 为每个控制器添加控制器测试文件
3. 为重要的用户流程添加集成测试
4. 使用 `assert` 包进行断言
5. 为复杂的测试场景准备测试数据

## 模拟数据

`tests/repository/mock_repository.go` 提供了模拟数据库实现，可用于控制器和集成测试。使用方法：

```go
// 初始化模拟仓库
mockRepo := repository.NewMockRepositories()
mockRepo.InitMockDB()

// 现在可以使用 repository.Repos 进行测试
``` 