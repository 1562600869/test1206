package main

import (
	"fmt"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("社区共享工具站 - 工具借用归还管理系统")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run main.go <命令> [参数...]")
	fmt.Println()
	fmt.Println("可用命令:")
	fmt.Println("  add-tool <ID> <名称> --type <类型> --deposit <押金(分)>")
	fmt.Println("      添加新工具")
	fmt.Println("      类型可选: 电动工具/手动工具/测量工具/园艺工具")
	fmt.Println()
	fmt.Println("  borrow <ID> --borrower <借用人> --phone <电话> --days <天数>")
	fmt.Println("      借出工具（仅在库状态可借出）")
	fmt.Println()
	fmt.Println("  return <ID> --condition <状况>")
	fmt.Println("      归还工具")
	fmt.Println("      状况可选: 完好/轻微损耗/有损坏")
	fmt.Println()
	fmt.Println("  overdue")
	fmt.Println("      列出超期未还的工具")
	fmt.Println()
	fmt.Println("  monthly --month YYYY-MM")
	fmt.Println("      某月借出统计报告")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	var err error
	switch cmd {
	case "add-tool":
		err = AddTool(args)
	case "borrow":
		err = BorrowTool(args)
	case "return":
		err = ReturnTool(args)
	case "overdue":
		err = ListOverdue(args)
	case "monthly":
		err = MonthlyReport(args)
	case "-h", "--help", "help":
		printUsage()
		return
	default:
		fmt.Printf("未知命令: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
