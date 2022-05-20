# Go 开发套件（Go Development Kits）

## 简介

go mod + 私有仓库的方式管理公共开发套件。

## 布局

```plain
.
├── go.mod # go module 声明文件
├── go.sum # go module lock file
├── pkg # 所有都放在 pkg 下，如果不希望暴露的，可以在任意层级建立 internal 文件夹
│   ├── app_context # 按功能划分文件夹
│   │   ├── app_context.go
│   │   └── README.md
│   └── kube_utils
│       ├── client.go
│       └── README.md
├── README.md
├── tests # 单元测试
│   └── kube_utils
│       ├── client_test.go
│       └── kubeconfig
└── vendor # 第三方包，已加入 gitignore 请不要同步此包
```

## 配置方式

```bash
# 配置开启 GoMod
go env -w GO111MODULE="on"

# 配置 GoMod 私有仓库
go env -w GOPRIVATE="git@github.com"
```

公有仓库使用方式：

如果自己修改过公有仓库，可以推到远程仓库的新分支，如 v0.0.1，然后项目需要引用此公有仓库，则使用命令:

```bash
# 即加上分支号
go get github.com/bigquant/GoDevKits@v0.0.1
```

推荐 tag 使用方式：vx.y.z

- v 为固定字符
- x 为主版本，0 表示仍在测试中，1 表示第一次 release，更新 golang 版本时 +1
- y 为次版本，新增模块时 +1
- z 为 Fix 版本，修复 bug 时 +1
- 可附加 `-alpha-1`、`-beta-1` 表示测试中的，`-rc-1` 用于准备正式发布的分支，数值仍采用递增

如果在项目中引用了，依次执行：

```bash
go mod tidy

go mod vendor
```

此项目中 vendor 加入 git 管理，但每次更新 vendor 请在单独的 commit 中，并标注为更新 commit。

在您的项目中，需要将 vendor 纳入版本管理，方便管理版本和构建二进制程序。

## 示例

用 go mod 初始化一个 golang 工程，新建 test.go 如下：

```go
package kube_utils

import (
  "context"
  "testing"

  "github.com/bigquant/GoDevKits/pkg/kube_utils"
  metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateClient(t *testing.T) {
  opts := &kube_utils.KubeClientOpts{
    InCluster:      false,
    KubeConfigPath: "./kubeconfig",
    MasterURL:      "https://bigquant:6443",
    Namespace:      "bigquant",
    AllNamespaces:  false,
  }
  client := kube_utils.NewKubeClient(opts)
  kubeCli := client.GetClientSetOrDie()
  if res, err := kubeCli.CoreV1().Namespaces().List(
      context.Background(),
      metaV1.ListOptions{},
  ); err != nil {
      t.Errorf("Cannot list namespaces")
  } else {
      t.Logf("Get res: %v", res.Items[0])
  }
}

```

执行

```bash
go mod tidy

go mod vendor
```

可以看到 vendor 分区中已经加载了此项目中的代码
