package main

import "fmt"

func SetColor(msg string, conf, bg, text int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, text, msg, 0x1B)
}

func main() {
	fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	conf := 1  // 配置、终端默认设置
	bg := 40   // 背景色、终端默认设置
	text := 31 // 前景色、红色
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "testPrintColor", 0x1B)

	str := SetColor("pandaychen", 1, 40, 31)
	fmt.Println(str)
}
