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

	str := SetColor("pandaychen", 1, 41, 31)
	fmt.Println(str)
	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}

	// 对号的 Unicode 代码
	CheckSymbol := "\u2714 "
	// 叉号的 Unicode 代码
	CrossSymbol := "\u2716 "
	fmt.Println(CheckSymbol, CrossSymbol)
}
