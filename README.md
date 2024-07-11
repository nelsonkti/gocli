```stylint
                      _  _ 
   __ _   ___    ___ | |(_)
  / _` | / _ \  / __|| || |
 | (_| || (_) || (__ | || |
  \__, | \___/  \___||_||_|
  |___/                                     
```

## gocli 脚手架工具
gocli 脚手架工具
这是一款自动生成项目模板的工具，它能根据用户的选择和需求，创建项目的目录结构、文件、依赖项等。使用这个脚手架工具可以节省手动创建和配置项目的时间和精力，从而让开发者更专注于业务逻辑的开发

### 安装 gocli 工具

使用以下命令安装 gocli 工具：

```sh
go get -u github.com/nelsonkti/gocli
```

### 配置环境变量

对于 Windows 用户，如果希望全局使用 `gocli` 命令，请将其可执行文件路径 (`${GOPATH}/src/bin/gocli.exe`) 添加到系统的环境变量中。

### 使用示例

#### 创建 Model

使用 `gocli` 创建 Model：

```sh
gocli make:model -d="root:root+@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local" -n=测试 -f=code_user_from_ma,v3_order -p=order 
```

#### 创建 Repository

使用 `gocli` 创建 Repository：

```sh
gocli make:repository -n=测试 -f=code_user_from_ma,v3_order -p=order 
```

#### 创建 Service

使用 `gocli` 创建 Service：

```sh
gocli make:service -n=测试 -f=code_user_from_ma,v3_order -p=order 
```

#### 一键创建 Model、Repository 和 Service

使用 `gocli` 一键创建 Model、Repository 和 Service：

```sh
gocli make:mrs -d="root:root+@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local" -f=code_user_from_ma,v3_order -p=order/test 
```

#### 创建 RPC

使用 `gocli` 创建 RPC：

```sh
gocli make:rpc -p=./proto/demo
```

### 环境要求

- Go 版本需 >= 1.20