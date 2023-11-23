# Storage
### 存储包

### Install
```shell
go get -u github.com/jiebutech/storage
```

### Quick Start


1. 注入驱动
```go
import _ "github.com/jiebutech/storage/driver/s3"
```

2. 实例化客户端
```go

cfg := config.StorageConfig{
    AccessKey:     "AccessKey",
    AccessSecret:  "AccessSecret",
    DefaultBucket: "test",
    Domain:        "http://www.tests3.com",
}
cli := storage.New(cfg)

url, err := cli.Upload(context.Background(), "./1", strings.NewReader("123"))

println(url)

```
