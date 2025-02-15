# Follow Service

Follow Service 是一个基于Go语言开发的微服务，用于处理用户关注关系管理的相关功能。该服务是陌生人聊天应用的核心组件之一，提供关注、取消关注以及查询关注关系等功能。

## 功能特性

- 关注/取消关注用户
- 获取关注列表
- 获取粉丝列表
- 获取互关用户列表
- 提供gRPC接口供其他服务调用
- JWT认证支持
- MongoDB数据持久化

## 技术栈

- Go 1.22+
- Gin Web Framework
- gRPC
- MongoDB
- Protocol Buffers
- JWT Authentication

## 依赖服务

- User Service: 用户信息服务
- Post Service: 帖子内容服务
- MongoDB: 数据存储

## 快速开始

### 前置条件

- Go 1.22 或更高版本
- MongoDB 4.0 或更高版本
- Protocol Buffers 编译器

### 安装

1. 克隆项目

```bash
git clone <repository-url>
cd follow-service
```

2. 安装依赖

```bash
go mod download
```

3. 配置服务
创建 `config/config.yaml` 文件并配置以下内容：
```yaml
server:
  port: 8080

grpcServer:
  port: 50051

mongodb:
  uri: "mongodb://localhost:27017"
  database: "followdb"
  collection: "follows"
  replicaSet: "rs0"

userService:
  host: "localhost:50052"

postService:
  host: "localhost:50053"
```

4. 启动服务
```bash
go run main.go
```

## API 文档

### HTTP接口

#### 关注用户

```
POST /api/v1/follow/user
Authorization: Bearer <token>
```

#### 取消关注
```
DELETE /api/v1/follow/user?targetUserId=<user-id>
Authorization: Bearer <token>
```

#### 获取关注列表
```
GET /api/v1/follow/my-follows?limit=10&offset=0
Authorization: Bearer <token>
```

#### 获取粉丝列表
```
GET /api/v1/follow/my-fans?limit=10&offset=0
Authorization: Bearer <token>
```

#### 获取互关列表
```
GET /api/v1/follow/mutual?limit=10&offset=0
Authorization: Bearer <token>
```

### gRPC接口

服务定义详见 `proto/follow.proto`：
- GetFollowCount: 获取用户的关注数和粉丝数
- GetFollowingUserIds: 获取用户关注的所有用户ID

## 项目结构

```
.
├── config/         # 配置文件
├── handlers/       # HTTP和gRPC处理器
├── middleware/     # 中间件
├── models/        # 数据模型
├── proto/         # Protocol Buffers定义
├── main.go        # 程序入口
└── README.md      # 项目文档
```

## 开发

### 生成Protocol Buffers代码

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

### 运行测试

```bash
go test ./...
```

## 部署

服务支持Docker部署，详细部署文档请参考部署指南。

## 贡献

欢迎提交Issue和Pull Request。
