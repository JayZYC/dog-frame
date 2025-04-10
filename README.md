# Dog-Frame 项目

Dog-Frame 是一个基于 Gin 框架的 Go 语言 Web 应用骨架，提供了完善的项目结构和丰富的功能组件，帮助开发者快速构建高性能、可维护的 Web 应用。

## 特性

- **完善的项目结构**：采用领域驱动设计（DDD）思想，清晰的分层架构
- **丰富的中间件**：
  - CORS 跨域支持
  - XSS 防护
  - CSRF 防护
  - 请求速率限制
  - IP 黑白名单
  - 请求追踪（Trace）
  - 访问日志记录
  - Prometheus 指标收集
- **全面的监控系统**：
  - 健康检查接口
  - 性能指标监控
  - 请求统计分析
- **优化的日志系统**：
  - 结构化日志
  - 日志轮转
  - 调用链追踪
- **数据库支持**：
  - 基于 GORM 的数据库操作
  - 读写分离支持
  - 数据库连接池管理
- **缓存支持**：
  - Redis 缓存集成
  - 缓存操作封装

## 项目结构

```
dog-frame/
├── api/                  # API 层，处理 HTTP 请求
│   ├── controller/       # 控制器，处理请求逻辑
│   ├── request/          # 请求模型定义
│   ├── response/         # 响应模型定义
│   └── router/           # 路由定义
├── common/               # 公共组件
│   ├── app/              # 应用相关工具
│   ├── enum/             # 枚举定义
│   ├── errcode/          # 错误码定义
│   ├── logger/           # 日志组件
│   ├── metrics/          # 监控指标
│   ├── middleware/       # 中间件
│   └── util/             # 工具函数
├── config/               # 配置管理
├── dal/                  # 数据访问层
│   ├── cache/            # 缓存操作
│   ├── dao/              # 数据库访问对象
│   └── model/            # 数据库模型
├── logic/                # 业务逻辑层
│   ├── appservice/       # 应用服务
│   ├── do/               # 领域对象
│   └── domainservice/    # 领域服务
└── main.go               # 应用入口
```

## 快速开始

### 环境要求

- Go 1.18+
- MySQL 5.7+
- Redis 6.0+

### 安装

1. 克隆项目

```bash
git clone https://github.com/jayzyc/dog-frame.git
cd dog-frame
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置数据库和 Redis

编辑 `config/application.dev.yaml` 文件，配置数据库和 Redis 连接信息：

```yaml
database:
  master:
    type: mysql
    dsn: root:password@tcp(localhost:3306)/dog-frame?charset=utf8&parseTime=True&loc=Asia%2FShanghai
    maxopen: 100
    maxidle: 10
    maxlifetime: 300000000000
  slave:
    type: mysql
    dsn: root:password@tcp(localhost:3306)/dog-frame?charset=utf8&parseTime=True&loc=Asia%2FShanghai
    maxopen: 100
    maxidle: 10
    maxlifetime: 300000000000
redis:
  addr: localhost:6379
  password:
  pool_size: 10
  db: 0
```

4. 运行应用

```bash
go run main.go
```

应用将在 `http://localhost:8080` 启动。

### 中间件注册

在 `api/router/router.go` 中注册中间件：

```go
func Register(e *gin.Engine) {
	// 初始化路由
	e.Use(
		middleware.StartTrace,       // 请求追踪
		middleware.LogAccess,         // 访问日志
		middleware.GinPanicRecovery,  // 异常恢复
		middleware.CORS(),            // CORS 跨域
		middleware.XSS(),             // XSS 防护
		middleware.CSRF(),            // CSRF 防护
		middleware.RateLimit(),       // 速率限制
		middleware.IPFilter(),        // IP 黑白名单
		middleware.PrometheusMiddleware(), // Prometheus 监控
	)

	r := e.Group("")

	// 注册路由
	registerRoutes(r)
}
```

## 中间件使用指南

### CORS 中间件

CORS 中间件用于处理跨域请求，支持自定义配置：

```go
// 使用默认配置
router.Use(middleware.CORS())

// 使用自定义配置
config := middleware.CORSConfig{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           86400,
}
router.Use(middleware.CORSWithConfig(config))
```

### XSS 防护中间件

XSS 防护中间件用于过滤请求中的恶意脚本：

```go
// 使用默认配置
router.Use(middleware.XSS())

// 使用自定义配置
config := middleware.XSSConfig{
    FilterRequestBody:  true,
    FilterResponseBody: false,
    FilterURLParams:    true,
    FilterHeaders:      true,
    FilterHeaderNames:  []string{"Cookie", "Authorization"},
}
router.Use(middleware.XSSWithConfig(config))
```

### CSRF 防护中间件

CSRF 防护中间件用于防止跨站请求伪造攻击：

```go
// 使用默认配置
router.Use(middleware.CSRF())

// 使用自定义配置
config := middleware.CSRFConfig{
    Secret:         "your-csrf-secret",
    CookieName:     "csrf_token",
    HeaderName:     "X-CSRF-Token",
    FormFieldName:  "_csrf",
    CookieMaxAge:   86400,
    CookieSecure:   true,
    CookieHTTPOnly: true,
    CookiePath:     "/",
    CookieDomain:   "example.com",
    ExcludedPaths:  []string{"/api/health"},
}
router.Use(middleware.CSRFWithConfig(config))
```

### 请求速率限制中间件

速率限制中间件用于限制 API 请求频率：

```go
// 使用默认配置
router.Use(middleware.RateLimit())

// 使用自定义配置
config := middleware.RateLimiterConfig{
    Rate:            10,  // 每秒请求数
    Burst:           20,  // 突发请求数
    ExpirationTime:  3600, // 过期时间（秒）
    CleanupInterval: 60,   // 清理间隔（秒）
    KeyFunc: func(c *gin.Context) string {
        return c.ClientIP() // 自定义限流键
    },
    ExcludedPaths: []string{"/api/health"},
}
router.Use(middleware.RateLimitWithConfig(config))
```

### IP 黑白名单中间件

IP 黑白名单中间件用于控制 IP 访问权限：

```go
// 使用默认配置（黑名单模式）
router.Use(middleware.IPFilter())

// 使用自定义配置
config := middleware.IPFilterConfig{
    Mode:          middleware.IPFilterModeWhitelist, // 白名单模式
    IPs:           []string{"192.168.1.1", "10.0.0.1"},
    CIDRs:         []string{"192.168.0.0/24"},
    ExcludedPaths: []string{"/api/health"},
}
router.Use(middleware.IPFilterWithConfig(config))

// 动态添加 IP
filter := middleware.GetIPFilter()
filter.AddIP("192.168.1.2")
filter.AddCIDR("10.0.0.0/24")
```

### Prometheus 监控中间件

Prometheus 监控中间件用于收集应用性能指标：

```go
// 使用 Prometheus 中间件
router.Use(middleware.PrometheusMiddleware())

// 暴露 Prometheus 指标
router.GET("/metrics", controller.PrometheusHandler())
```

## 日志系统

Dog-Frame 使用 zap 作为日志库，提供了结构化日志记录功能：

```go
// 记录信息日志
logger.Info(ctx, "用户登录成功", "user_id", userId, "ip", c.ClientIP())

// 记录警告日志
logger.Warn(ctx, "登录尝试失败", "user_id", userId, "reason", "密码错误")

// 记录错误日志
logger.Error(ctx, "数据库连接失败", "err", err)

// 记录调试日志
logger.Debug(ctx, "API 请求参数", "params", params)
```

## 错误处理

Dog-Frame 提供了统一的错误处理机制：

```go
// 返回成功响应
app.NewResponse(c).Success(data)

// 返回错误响应
app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))

// 自定义错误
var ErrUserNotFound = errcode.NewError(10100, "用户不存在")
app.NewResponse(c).Error(ErrUserNotFound)
```

## 监控系统

### 健康检查

访问 `/health` 和 `/ready` 端点可以检查应用的健康状态和就绪状态。

### Prometheus 指标

访问 `/metrics` 端点可以查看 Prometheus 指标，包括：

- `http_requests_total` - HTTP 请求总数
- `http_request_duration_seconds` - HTTP 请求处理时间
- `http_request_size_bytes` - HTTP 请求大小
- `http_response_size_bytes` - HTTP 响应大小
- `http_request_errors_total` - HTTP 请求错误总数

## 配置管理

Dog-Frame 使用 Viper 管理配置，支持不同环境的配置文件：

- `config/application.dev.yaml` - 开发环境配置
- `config/application.test.yaml` - 测试环境配置
- `config/application.prod.yaml` - 生产环境配置

通过环境变量 `ENV` 指定使用的配置文件：

```bash
ENV=dev go run main.go  # 使用开发环境配置
ENV=prod go run main.go # 使用生产环境配置
```

## 贡献指南

欢迎贡献代码或提出建议！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 联系方式

如有任何问题或建议，请通过以下方式联系我们：

- 邮箱：584405019@qq.com
- GitHub Issues：https://github.com/jayzyc/dog-frame/issues
