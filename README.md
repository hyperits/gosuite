# gosuite

> gosuite provides a suite of tools to help you develop your server/web go project.

## 组件简介

### db - 数据库客户端

提供统一的数据库客户端接口和实现，支持连接池配置。

| 子包 | 描述 |
|------|------|
| `db` | 数据库客户端公共接口定义（`Client`、`SQLClient`、`KVClient`） |
| `db/mysql` | MySQL 客户端，基于 GORM，支持连接池管理 |
| `db/postgres` | PostgreSQL 客户端，基于 GORM，支持 SSL 和时区配置 |
| `db/redis` | Redis 客户端，支持单机、哨兵、集群三种模式 |

### errors - 错误处理

统一的错误处理包，提供标准错误变量和错误包装函数。

- 预定义错误：`ErrNilConfig`、`ErrNotConnected`、`ErrTimeout` 等
- 错误包装：`Wrap`、`Wrapf` 添加上下文信息
- 操作错误：`OpError` 结构化错误类型

### kit - 工具包

| 子包 | 描述 |
|------|------|
| `kit/cmd` | 命令执行工具，封装 `os/exec` |
| `kit/conv` | 类型转换工具，包括对象转 JSON、对象转 Map |
| `kit/debug` | 运行时信息获取，如当前函数名、文件、行号 |

### logger - 日志

基于 [zerolog](https://github.com/rs/zerolog) 的日志组件，支持：

- 多输出目标（控制台 + 文件）
- 日志文件自动轮转（基于 lumberjack）
- 多日志级别（Debug/Info/Warn/Error/Fatal/Panic）
- 结构化日志和运行时信息注入

### net - 网络

| 子包 | 描述 |
|------|------|
| `net/httpx` | HTTP 客户端封装，支持 GET/POST/PUT/DELETE，函数式配置 |
| `net/mail` | 邮件发送接口定义（`Sender`），支持附件 |
| `net/sms` | 短信发送接口定义（`Sender`），支持模板参数 |

### providers - 第三方服务

具体的第三方服务实现：

| 子包 | 描述 |
|------|------|
| `providers/aliyun/sms` | 阿里云短信服务，实现 `sms.Sender` 接口 |
| `providers/smtp/mail` | SMTP 邮件服务，实现 `mail.Sender` 接口 |

### security - 安全

| 子包 | 描述 |
|------|------|
| `security/captcha` | 图形验证码，支持 Redis 存储，返回 Base64 图片 |
| `security/hash` | 密码哈希，使用 bcrypt 算法，支持自定义成本因子 |
| `security/verify` | 数字验证码，基于 Redis 存储，支持自定义长度和过期时间 |

### storage - 存储

| 子包 | 描述 |
|------|------|
| `storage/s3` | S3 兼容对象存储客户端，基于 MinIO SDK，支持上传/下载/删除/预签名 URL |

---

## Git 提交格式

+ `feat` 添加了新特性
+ `fix` 修复问题
+ `style` 无逻辑改动的代码风格调整
+ `perf` 性能/优化
+ `refactor` 重构
+ `revert` 回滚提交
+ `test` 测试
+ `docs` 文档
+ `chore` 依赖或者脚手架调整
+ `workflow` 工作流优化
+ `ci` 持续集成
+ `types` 类型定义
+ `wip` 开发中

---

## 参考资料

+ [Google Go 编程规范](https://gocn.github.io/styleguide/)
