##	使用 go-check 进行单元测试


gocheck 作为 golang 的一种测试框架，可以直接继承 `go test` 使用，允许之前基于 testing 框架的测试平滑迁移到 gocheck 框架而不会发生冲突，gocheck API 与 testing 也有很多相似之处：

-   gocheck 拥有丰富的断言 api，如 `Assert`、`Check`、`Skip` 等
-   gocheck 单元测试代码最基本的结构如例子所示，声明一个测试套件，在测试套件下写测试方法，即测试用例
-   直接使用 go test 在测试文件当前路径下运行

如 `1_test.go` 的运行结果如下：

```text
[root@VM_120_245_centos /data/github_own/]# go test -v 1_test.go
=== RUN   Test

----------------------------------------------------------------------
FAIL: 1_test.go:16: MySuite.TestHelloWorld

1_test.go:17:
    c.Assert(42, Equals, "42")
... obtained int = 42
... expected string = "42"

OOPS: 0 passed, 1 FAILED
--- FAIL: Test (0.00s)
FAIL
FAIL    command-line-arguments  0.003s
FAIL
```


## 参考
-   [golang 分层测试之单元测试 - gocheck 使用](https://studygolang.com/articles/16503)