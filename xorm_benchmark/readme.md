##	xorm插入性能测试

建立测试表：

```sql
create table users(
    id BIGINT NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(32) NOT NULL,
    last_name VARCHAR(32) NOT NULL,
    extra1 VARCHAR(32),
    extra2 VARCHAR(32),
    version BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (id)
);
```

####	tools
```text
go test -bench .
```


####	case1
```text
[root@VM_120_245_centos ~/sql_benchmark/case1]# go test -bench .
goos: linux
goarch: amd64
cpu: AMD EPYC Processor
BenchmarkInsert1-8                   109          11175403 ns/op
BenchmarkInsert10-8                   94          12795190 ns/op
BenchmarkInsert100-8                  30          39982011 ns/op
BenchmarkInsert1000-8                  1        2348321509 ns/op
PASS
ok      _/root/sql_benchmark/case1      13.037s
```
