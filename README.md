# 矿山生态保护修复协同平台

该后端服务用于矿山生态修复项目、地块、监测样本、治理任务和验收协同。系统采用 Go 1.22、SQLite 关系数据库和版本化迁移，提供登录、可撤销会话、RBAC、幂等键、乐观锁状态流、审计和后台提醒 worker。

## 本地运行

```text
GOTOOLCHAIN=local go run .
```

默认账号：`admin@eco.local/admin123`、`inspector@eco.local/inspect123`、`operator@eco.local/operate123`。生产环境应通过环境变量覆盖数据库和会话配置。

主要 API：`/healthz`、`/readyz`、`POST /api/auth/login`、`POST /api/auth/logout`、`GET|POST /api/projects`、`PUT /api/projects?plot=1`、`GET|POST|PATCH /api/tasks`、`GET|POST /api/samples`、`POST /api/inspections`。

## 数据与恢复

启动时按文件名顺序应用 `internal/migrations/sql` 中的迁移并记录 `schema_migrations`。所有核心写入使用真实 SQL；治理任务创建和验收在同一事务中更新关联实体、事件和审计，服务重启后任务状态从数据库恢复。

## 验证

```text
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go test -race ./... -count=1
GOTOOLCHAIN=local go vet ./...
GOTOOLCHAIN=local go build ./...
```
