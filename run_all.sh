#!/bin/bash

# Golang并发编程学习Demo运行脚本
# 用于方便地运行各个级别的并发编程示例

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
    echo -e "${CYAN}====================================${NC}"
    echo -e "${CYAN}    Golang并发编程学习Demo${NC}"
    echo -e "${CYAN}====================================${NC}"
    echo ""
}

# 打印分级说明
print_levels() {
    echo -e "${YELLOW}学习级别说明：${NC}"
    echo -e "${GREEN}Simple (简单级别)${NC} - 基础并发概念，适合初学者"
    echo -e "${BLUE}Medium (中等级别)${NC} - 实际应用模式，适合有一定基础者"
    echo -e "${PURPLE}Hard (困难级别)${NC} - 企业级系统，适合深入学习者"
    echo ""
}

# 运行指定文件
run_demo() {
    local file="$1"
    local name="$2"
    
    echo -e "${CYAN}运行: ${name}${NC}"
    echo -e "${YELLOW}文件: ${file}${NC}"
    echo "按 Enter 键开始运行，或按 Ctrl+C 退出..."
    read
    
    echo -e "${GREEN}======== 开始运行 ========${NC}"
    go run "$file"
    echo -e "${GREEN}======== 运行完成 ========${NC}"
    echo ""
    echo "按 Enter 键继续下一个demo..."
    read
    clear
}

# Simple级别demos
run_simple_demos() {
    echo -e "${GREEN}=== 简单级别 (Simple) ===${NC}"
    echo "这个级别包含10个基础并发编程示例"
    echo ""
    
    declare -A simple_demos=(
        ["simple/01_basic_goroutine.go"]="基础Goroutine使用"
        ["simple/02_waitgroup_basic.go"]="WaitGroup基础使用"
        ["simple/03_channel_basic.go"]="基础Channel通信"
        ["simple/04_buffered_channel.go"]="缓冲Channel使用"
        ["simple/05_select_basic.go"]="基础Select语句"
        ["simple/06_timeout_select.go"]="带超时的Select"
        ["simple/07_mutex_basic.go"]="基础互斥锁使用"
        ["simple/08_once_basic.go"]="sync.Once使用"
        ["simple/09_channel_pipeline.go"]="简单Channel管道"
        ["simple/10_goroutine_pool.go"]="简单Goroutine池"
    )
    
    for file in simple/0*.go; do
        if [[ -f "$file" ]]; then
            name="${simple_demos[$file]}"
            run_demo "$file" "$name"
        fi
    done
}

# Medium级别demos
run_medium_demos() {
    echo -e "${BLUE}=== 中等级别 (Medium) ===${NC}"
    echo "这个级别包含10个实际应用的并发模式"
    echo ""
    
    declare -A medium_demos=(
        ["medium/01_producer_consumer.go"]="生产者消费者模式"
        ["medium/02_worker_pool_advanced.go"]="高级工作池"
        ["medium/03_rate_limiter.go"]="速率限制器"
        ["medium/04_publish_subscribe.go"]="发布订阅模式"
        ["medium/05_context_cancellation.go"]="Context取消机制"
        ["medium/06_fan_in_fan_out.go"]="扇入扇出模式"
        ["medium/07_circuit_breaker.go"]="熔断器模式"
        ["medium/08_semaphore.go"]="信号量实现"
        ["medium/09_actor_model.go"]="Actor模型"
        ["medium/10_pipeline_processing.go"]="流水线处理"
    )
    
    for file in medium/0*.go; do
        if [[ -f "$file" ]]; then
            name="${medium_demos[$file]}"
            run_demo "$file" "$name"
        fi
    done
}

# Hard级别demos
run_hard_demos() {
    echo -e "${PURPLE}=== 困难级别 (Hard) ===${NC}"
    echo "这个级别包含4个企业级并发编程示例"
    echo ""
    
    declare -A hard_demos=(
        ["hard/01_distributed_worker.go"]="分布式工作者系统"
        ["hard/02_load_balancer.go"]="负载均衡器"
        ["hard/03_message_queue.go"]="消息队列系统"
        ["hard/04_connection_pool.go"]="连接池管理"
    )
    
    for file in hard/0*.go; do
        if [[ -f "$file" ]]; then
            name="${hard_demos[$file]}"
            run_demo "$file" "$name"
        fi
    done
}

# 主菜单
show_menu() {
    print_title
    print_levels
    
    echo -e "${YELLOW}请选择要运行的级别：${NC}"
    echo "1) Simple - 简单级别 (10个demo)"
    echo "2) Medium - 中等级别 (10个demo)"  
    echo "3) Hard - 困难级别 (4个demo)"
    echo "4) All - 运行所有demo (24个demo)"
    echo "5) 退出"
    echo ""
    
    read -p "请输入选择 (1-5): " choice
    
    case $choice in
        1)
            clear
            run_simple_demos
            ;;
        2)
            clear
            run_medium_demos
            ;;
        3)
            clear
            run_hard_demos
            ;;
        4)
            clear
            echo -e "${CYAN}开始运行所有Demo...${NC}"
            echo ""
            run_simple_demos
            run_medium_demos
            run_hard_demos
            echo -e "${GREEN}所有Demo运行完成！${NC}"
            ;;
        5)
            echo -e "${GREEN}谢谢使用！祝学习愉快！${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}无效选择，请重新选择${NC}"
            echo ""
            show_menu
            ;;
    esac
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

# 使用说明
show_usage() {
    echo -e "${YELLOW}使用说明：${NC}"
    echo "• 每个demo都是独立的程序，演示特定的并发编程概念"
    echo "• 建议按照 Simple -> Medium -> Hard 的顺序学习"
    echo "• 运行时注意观察输出，理解并发执行的特点"
    echo "• 可以修改代码参数来观察不同的行为"
    echo ""
    echo -e "${YELLOW}注意事项：${NC}"
    echo "• linter可能会显示'main redeclared'警告，这是正常现象"
    echo "• 每个文件都是独立的main程序，单独运行即可"
    echo "• 如果遇到问题，请查看README.md获取更多信息"
    echo ""
}

# 主程序
main() {
    clear
    check_go
    show_usage
    
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