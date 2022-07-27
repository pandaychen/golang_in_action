package hello_test

import (
	"io"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) } //继承testing的方法，可以直接使用go test命令运行

type MySuite struct{} //创建测试套件结构体

var _ = Suite(&MySuite{})

func (s *MySuite) TestHelloWorld(c *C) { //声明TestHelloWorld方法为MySuite套件的测试用例
	c.Assert(42, Equals, "42")
	//stop
	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
	c.Check(42, Equals, 42)
    	c.Assert(42, Equals, 42)
}
