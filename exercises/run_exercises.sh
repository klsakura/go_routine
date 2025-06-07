#!/bin/bash

# Golang并发编程练习运行脚本
# 帮助选择和运行特定的练习文件

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印标题
print_title() {
    echo -e "${CYAN}======================================${NC}"
    echo -e "${CYAN}   Golang并发编程练习系统${NC}"
    echo -e "${CYAN}======================================${NC}"
    echo ""
    echo -e "${YELLOW}📚 学习提示：${NC}"
    echo -e "1. 先运行对应的demo了解概念"
    echo -e "2. 再做练习实现相同功能"
    echo -e "3. 对比分析不同的实现方案"
    echo ""
}

# 运行练习文件
run_exercise() {
    local file="$1"
    local name="$2"
    
    if [[ ! -f "$file" ]]; then
        echo -e "${RED}错误: 文件 $file 不存在${NC}"
        return 1
    fi
    
    echo -e "${CYAN}=== 练习: ${name} ===${NC}"
    echo -e "${YELLOW}文件: ${file}${NC}"
    echo ""
    echo -e "${GREEN}💡 提示：${NC}"
    echo -e "- 这是一个练习文件，包含TODO标记的代码需要您来实现"
    echo -e "- 请先打开文件查看练习要求和提示"
    echo -e "- 实现代码后再运行测试"
    echo ""
    echo "按 Enter 键查看练习文件内容，或按 Ctrl+C 跳过..."
    read
    
    # 显示文件前50行，让用户了解练习内容
    echo -e "${BLUE}=== 练习文件内容预览 ===${NC}"
    head -50 "$file"
    echo ""
    echo -e "${YELLOW}... (更多内容请直接打开文件查看) ...${NC}"
    echo ""
    
    echo "按 Enter 键尝试运行（可能会有编译错误，这是正常的），或按 Ctrl+C 跳过..."
    read
    
    echo -e "${GREEN}======== 尝试运行练习 ========${NC}"
    go run "$file" 2>&1 || echo -e "${YELLOW}注意：出现错误是正常的，请根据TODO提示完成代码实现${NC}"
    echo -e "${GREEN}======== 运行结束 ========${NC}"
    echo ""
    echo "按 Enter 键继续下一个练习..."
    read
    clear
}

# Simple级别练习
run_simple_exercises() {
    echo -e "${GREEN}=== 简单级别练习 (Simple) ===${NC}"
    echo "这个级别包含基础并发编程概念的练习"
    echo ""
    
    local exercises=(
        "simple/01_basic_goroutine_exercise.go:基础Goroutine练习"
        "simple/02_waitgroup_basic_exercise.go:WaitGroup同步练习"  
        "simple/03_channel_basic_exercise.go:基础Channel通信练习"
    )
    
    for exercise in "${exercises[@]}"; do
        IFS=':' read -r file name <<< "$exercise"
        if [[ -f "$file" ]]; then
            run_exercise "$file" "$name"
        else
            echo -e "${YELLOW}练习文件 $file 暂未创建${NC}"
        fi
    done
}

# Medium级别练习
run_medium_exercises() {
    echo -e "${BLUE}=== 中等级别练习 (Medium) ===${NC}"
    echo "这个级别包含实际应用中的并发模式练习"
    echo ""
    
    local exercises=(
        "medium/01_producer_consumer_exercise.go:生产者消费者模式练习"
    )
    
    for exercise in "${exercises[@]}"; do
        IFS=':' read -r file name <<< "$exercise"
        if [[ -f "$file" ]]; then
            run_exercise "$file" "$name"
        else
            echo -e "${YELLOW}练习文件 $file 暂未创建${NC}"
        fi
    done
}

# Hard级别练习
run_hard_exercises() {
    echo -e "${PURPLE}=== 困难级别练习 (Hard) ===${NC}"
    echo "这个级别包含企业级复杂系统的练习"
    echo ""
    
    local exercises=(
        "hard/01_distributed_worker_exercise.go:分布式工作者系统练习"
    )
    
    for exercise in "${exercises[@]}"; do
        IFS=':' read -r file name <<< "$exercise"
        if [[ -f "$file" ]]; then
            run_exercise "$file" "$name"
        else
            echo -e "${YELLOW}练习文件 $file 暂未创建${NC}"
        fi
    done
}

# 打开特定练习文件
open_exercise() {
    echo -e "${CYAN}请选择要打开的练习文件：${NC}"
    echo ""
    
    # 列出所有练习文件
    local count=1
    local files=()
    
    for level in simple medium hard; do
        if [[ -d "$level" ]]; then
            echo -e "${YELLOW}=== $level 级别 ===${NC}"
            for file in $level/*_exercise.go; do
                if [[ -f "$file" ]]; then
                    echo "$count) $(basename "$file")"
                    files[$count]="$file"
                    ((count++))
                fi
            done
            echo ""
        fi
    done
    
    read -p "请输入文件编号 (1-$((count-1))): " choice
    
    if [[ -n "${files[$choice]}" ]]; then
        local file="${files[$choice]}"
        echo -e "${GREEN}打开练习文件: $file${NC}"
        
        # 尝试用不同编辑器打开
        if command -v code >/dev/null; then
            code "$file"
        elif command -v vim >/dev/null; then
            vim "$file"
        elif command -v nano >/dev/null; then
            nano "$file"
        else
            echo "请手动打开文件: $file"
        fi
    else
        echo -e "${RED}无效选择${NC}"
    fi
}

# 创建新练习文件
create_exercise() {
    echo -e "${CYAN}练习文件创建助手${NC}"
    echo ""
    
    read -p "选择级别 (simple/medium/hard): " level
    read -p "输入练习文件名（不含扩展名）: " filename
    read -p "输入练习主题: " topic
    
    local filepath="$level/${filename}_exercise.go"
    
    if [[ -f "$filepath" ]]; then
        echo -e "${YELLOW}文件已存在: $filepath${NC}"
        return 1
    fi
    
    # 创建练习文件模板
    cat > "$filepath" << EOF
/*
Golang并发编程练习 - ${level^}级别
练习文件：${filename}_exercise.go
练习主题：$topic

练习目标：
1. TODO: 填写学习目标
2. TODO: 填写学习目标
3. TODO: 填写学习目标

练习任务：
- 任务1：TODO
- 任务2：TODO
- 任务3：TODO

运行方式：go run exercises/$level/${filename}_exercise.go
*/

package main

import (
	"fmt"
)

// TODO: 定义需要的结构体和接口

// TODO: 实现练习函数

func main() {
	fmt.Println("=== $topic 练习 ===")
	
	// 任务1：TODO
	fmt.Println("\\n任务1：TODO")
	// TODO: 在这里实现您的代码
	
	fmt.Println("任务1完成\\n")
	
	// 任务2：TODO  
	fmt.Println("任务2：TODO")
	// TODO: 在这里实现您的代码
	
	fmt.Println("任务2完成\\n")
	
	fmt.Println("所有练习完成！")
	
	// 反思问题：
	fmt.Println("\\n思考题：")
	fmt.Println("1. TODO：添加思考题")
	fmt.Println("2. TODO：添加思考题")
}
EOF

    echo -e "${GREEN}练习文件创建成功: $filepath${NC}"
}

# 显示使用统计
show_stats() {
    echo -e "${CYAN}=== 练习文件统计 ===${NC}"
    echo ""
    
    for level in simple medium hard; do
        if [[ -d "$level" ]]; then
            local count=$(find "$level" -name "*_exercise.go" | wc -l)
            echo -e "${YELLOW}$level 级别:${NC} $count 个练习文件"
        fi
    done
    
    echo ""
    echo -e "${CYAN}=== 完成情况检查 ===${NC}"
    echo "注意：以下只是简单检查，不代表代码质量"
    echo ""
    
    for level in simple medium hard; do
        if [[ -d "$level" ]]; then
            for file in $level/*_exercise.go; do
                if [[ -f "$file" ]]; then
                    local todo_count=$(grep -c "TODO" "$file" 2>/dev/null || echo "0")
                    local filename=$(basename "$file")
                    
                    if [[ $todo_count -eq 0 ]]; then
                        echo -e "${GREEN}✓${NC} $filename (可能已完成)"
                    elif [[ $todo_count -lt 5 ]]; then
                        echo -e "${YELLOW}◐${NC} $filename (部分完成，剩余 $todo_count 个TODO)"
                    else
                        echo -e "${RED}◯${NC} $filename (未开始，有 $todo_count 个TODO)"
                    fi
                fi
            done
        fi
    done
}

# 主菜单
show_menu() {
    print_title
    
    echo -e "${YELLOW}请选择操作：${NC}"
    echo "1) 运行 Simple 级别练习"
    echo "2) 运行 Medium 级别练习"
    echo "3) 运行 Hard 级别练习"
    echo "4) 打开特定练习文件编辑"
    echo "5) 查看练习统计信息"
    echo "6) 创建新练习文件"
    echo "7) 查看使用说明"
    echo "8) 退出"
    echo ""
    
    read -p "请输入选择 (1-8): " choice
    
    case $choice in
        1) clear; run_simple_exercises ;;
        2) clear; run_medium_exercises ;;
        3) clear; run_hard_exercises ;;
        4) clear; open_exercise ;;
        5) clear; show_stats ;;
        6) clear; create_exercise ;;
        7) clear; show_help ;;
        8) echo -e "${GREEN}感谢使用练习系统！祝学习愉快！${NC}"; exit 0 ;;
        *) echo -e "${RED}无效选择，请重新选择${NC}"; echo ""; show_menu ;;
    esac
}

# 显示使用说明
show_help() {
    echo -e "${CYAN}=== 练习系统使用说明 ===${NC}"
    echo ""
    echo -e "${YELLOW}1. 练习流程：${NC}"
    echo "   • 选择对应级别的练习"
    echo "   • 查看练习文件中的TODO标记"
    echo "   • 根据提示实现代码"
    echo "   • 运行测试验证结果"
    echo ""
    echo -e "${YELLOW}2. 文件结构：${NC}"
    echo "   • 每个练习文件都有详细的任务说明"
    echo "   • TODO标记指示需要实现的代码"
    echo "   • 提示信息帮助理解实现思路"
    echo ""
    echo -e "${YELLOW}3. 调试技巧：${NC}"
    echo "   • 使用 fmt.Printf 打印调试信息"
    echo "   • 使用 go run -race 检测竞态条件"
    echo "   • 逐步实现，每完成一部分就测试"
    echo ""
    echo -e "${YELLOW}4. 获取帮助：${NC}"
    echo "   • 参考对应的demo文件"
    echo "   • 查看练习文件中的提示"
    echo "   • 阅读 exercises/README.md"
    echo ""
}

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未找到Go环境，请先安装Go${NC}"
        echo "下载地址: https://golang.org/dl/"
        exit 1
    fi
    
    echo -e "${GREEN}Go环境检查通过: $(go version)${NC}"
    echo ""
}

# 主程序
main() {
    # 检查当前目录
    if [[ ! -d "simple" ]] && [[ ! -d "medium" ]] && [[ ! -d "hard" ]]; then
        echo -e "${RED}错误: 请在 exercises 目录下运行此脚本${NC}"
        exit 1
    fi
    
    clear
    check_go
    
    while true; do
        show_menu
        echo ""
        echo -e "${CYAN}按Enter键返回主菜单，或Ctrl+C退出${NC}"
        read
        clear
    done
}

# 运行主程序
main 