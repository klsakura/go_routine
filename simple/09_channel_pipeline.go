package main

import (
	"fmt"
)

// 简单的channel管道演示
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	fmt.Println("=== 简单Channel管道演示 ===")

	// 设置管道: generator -> square
	numbers := generator(1, 2, 3, 4, 5)
	squares := square(numbers)

	// 消费结果
	fmt.Println("原数字 -> 平方:")
	for result := range squares {
		fmt.Printf("结果: %d\n", result)
	}

	fmt.Println("管道处理完成！")
}
