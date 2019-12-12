package main

import (
    "fmt"
    "time"

    "gopkg.in/go-playground/validator.v9"
)

type Address struct {
    Street   string    `validate:"required"`
    City     string    `validate:"required"`
    Planet   string    `validate:"required"`
    Phone    string    `validate:"required"`
    Age      int       `validate:"min=12,max=15"`
    CreateAt time.Time `validate:"myParam=this is called param"`
}

func main() {
    address := &Address{
        Street: "Eavesdown Docks",
        City:   "beijing",
        Planet: "Persphone",
        Phone:  "none",
        Age:    16,
    }
    validate := validator.New()
        //自己定义tag标签以及与之对应的处理逻辑
    validate.RegisterValidation("myParam", mytimeFunc)
        //查看是否符合验证
    err := validate.Struct(address)
    fmt.Println(err)
}

func mytimeFunc(fl validator.FieldLevel) bool {
    fmt.Println("FieldName:", fl.FieldName())
    fmt.Println("StructFieldName", fl.StructFieldName())
    fmt.Println("Parm:", fl.Param())
    return false
}
