# gosuite

> gosuite provide a suite of tools to help you develop your server/web go project.

## 代码设计规范

### 文件命名规范

+ 文件名应使用小写字母，不应使用大写字母或特殊字符。
+ 如果文件名包含多个单词，应使用下划线（_）将它们连接起来。
+ 如果文件包含 Go 代码，则文件名应以 .go 为后缀。
+ 如果文件是 Go 包的测试文件，则文件名应以 _test.go 为后缀。

## packages

+ captcha - graphic captcha
+ converter - convert between different types
+ crypto - common used crypto
+ debugger - retrieve debug info of code [file, line, func]
+ httputil - http client util
+ logger - logger based on zerolog
+ mail - send email helper
+ mysqldb - mysql db helper
+ redisdb - redis db helper
+ s3db - s3 db helper
+ sms - send sms helper
+ verificationcode - server generated verification code management

## 参考资料

+ [Google Go 编程规范](https://gocn.github.io/styleguide/)
