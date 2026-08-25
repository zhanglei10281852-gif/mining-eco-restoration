# BENZHI_README

这是一个 Go 后端服务，用于该后端服务用于矿山生态修复项目、地块、监测样本、治理任务和验收协同。

## 项目说明

- 项目：zhanglei10281852-gif/mining-eco-restoration
- 项目用途：该后端服务用于矿山生态修复项目、地块、监测样本、治理任务和验收协同。系统采用 Go 1.22、SQLite 关系数据库和版本化迁移，提供登录、可撤销会话、RBAC、幂等键、乐观锁状态流、审计和后台提醒 worker。
- Go 工具链：`golang:1.22`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
cd '/app' && GOTOOLCHAIN=local go run .
cd '/app' && GOTOOLCHAIN=local go run ./cmd/server

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-352-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-352-arm64 linux/arm64
docker run -it benzhi-task-352-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-352-arm64:latest
```

## 题目验证命令

1. 预期退出码 0：`go test ./internal/notification -run '^TestPublishPropagatesSinkFailure$' -count=1`
2. 预期退出码 0：`go test ./...`
3. 预期退出码 0：`GOTOOLCHAIN=local go build -buildvcs=false ./... && GOTOOLCHAIN=local go vet ./...`
