package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func AddTool(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: add-tool <工具ID> <工具名称> --type <类型> --deposit <押金>")
	}
	id := args[0]
	name := args[1]
	toolType := ""
	depositStr := ""

	i := 2
	for i < len(args) {
		switch args[i] {
		case "--type":
			if i+1 >= len(args) {
				return fmt.Errorf("--type 缺少参数值")
			}
			toolType = args[i+1]
			i += 2
		case "--deposit":
			if i+1 >= len(args) {
				return fmt.Errorf("--deposit 缺少参数值")
			}
			depositStr = args[i+1]
			i += 2
		default:
			return fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if toolType == "" {
		return fmt.Errorf("缺少 --type 参数，可选值: %s", strings.Join(ValidToolTypes, "/"))
	}
	if !IsValidToolType(toolType) {
		return fmt.Errorf("工具类型无效，可选值: %s", strings.Join(ValidToolTypes, "/"))
	}
	if depositStr == "" {
		return fmt.Errorf("缺少 --deposit 参数")
	}
	deposit, err := strconv.Atoi(depositStr)
	if err != nil {
		return fmt.Errorf("押金必须是整数（单位：分）: %w", err)
	}
	if deposit < 0 {
		return fmt.Errorf("押金不能为负数")
	}

	store, err := LoadData()
	if err != nil {
		return err
	}
	if FindTool(store, id) != nil {
		return fmt.Errorf("工具ID %s 已存在", id)
	}

	tool := Tool{
		ID:      id,
		Name:    name,
		Type:    toolType,
		Deposit: deposit,
		Status:  StatusAvailable,
	}
	store.Tools = append(store.Tools, tool)

	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("工具添加成功: [%s] %s (%s), 押金: %d分\n", id, name, toolType, deposit)
	return nil
}

func BorrowTool(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: borrow <工具ID> --borrower <借用人> --phone <电话> --days <天数>")
	}
	id := args[0]
	borrower := ""
	phone := ""
	daysStr := ""

	i := 1
	for i < len(args) {
		switch args[i] {
		case "--borrower":
			if i+1 >= len(args) {
				return fmt.Errorf("--borrower 缺少参数值")
			}
			borrower = args[i+1]
			i += 2
		case "--phone":
			if i+1 >= len(args) {
				return fmt.Errorf("--phone 缺少参数值")
			}
			phone = args[i+1]
			i += 2
		case "--days":
			if i+1 >= len(args) {
				return fmt.Errorf("--days 缺少参数值")
			}
			daysStr = args[i+1]
			i += 2
		default:
			return fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if borrower == "" {
		return fmt.Errorf("缺少 --borrower 参数")
	}
	if phone == "" {
		return fmt.Errorf("缺少 --phone 参数")
	}
	if daysStr == "" {
		return fmt.Errorf("缺少 --days 参数")
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return fmt.Errorf("天数必须是整数: %w", err)
	}
	if days <= 0 {
		return fmt.Errorf("天数必须大于0")
	}

	store, err := LoadData()
	if err != nil {
		return err
	}
	tool := FindTool(store, id)
	if tool == nil {
		return fmt.Errorf("工具ID %s 不存在", id)
	}
	if tool.Status != StatusAvailable {
		return fmt.Errorf("工具 [%s] %s 当前状态为 %s，无法借出", id, tool.Name, tool.Status)
	}

	now := time.Now()
	borrowDate := FormatDate(now)
	expectedReturn := FormatDate(now.AddDate(0, 0, days))

	tool.Status = StatusBorrowed
	tool.Borrower = borrower
	tool.Phone = phone
	tool.BorrowDate = borrowDate
	tool.ExpectedReturnDate = expectedReturn
	tool.ActualReturnDate = ""
	tool.ReturnCondition = ""

	record := BorrowRecord{
		ToolID:             tool.ID,
		ToolName:           tool.Name,
		ToolType:           tool.Type,
		Borrower:           borrower,
		Phone:              phone,
		BorrowDate:         borrowDate,
		ExpectedReturnDate: expectedReturn,
	}
	store.Records = append(store.Records, record)

	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("借出成功: [%s] %s\n", id, tool.Name)
	fmt.Printf("  借用人: %s\n", borrower)
	fmt.Printf("  联系电话: %s\n", phone)
	fmt.Printf("  借出日期: %s\n", borrowDate)
	fmt.Printf("  预计归还: %s\n", expectedReturn)
	return nil
}

func ReturnTool(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: return <工具ID> --condition <完好/轻微损耗/有损坏>")
	}
	id := args[0]
	condition := ""

	i := 1
	for i < len(args) {
		switch args[i] {
		case "--condition":
			if i+1 >= len(args) {
				return fmt.Errorf("--condition 缺少参数值")
			}
			condition = args[i+1]
			i += 2
		default:
			return fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if condition == "" {
		return fmt.Errorf("缺少 --condition 参数，可选值: %s", strings.Join(ValidConditions, "/"))
	}
	if !IsValidCondition(condition) {
		return fmt.Errorf("归还状况无效，可选值: %s", strings.Join(ValidConditions, "/"))
	}

	store, err := LoadData()
	if err != nil {
		return err
	}
	tool := FindTool(store, id)
	if tool == nil {
		return fmt.Errorf("工具ID %s 不存在", id)
	}
	if tool.Status != StatusBorrowed {
		return fmt.Errorf("工具 [%s] %s 当前状态为 %s，不是借出中", id, tool.Name, tool.Status)
	}

	returnDate := FormatDate(time.Now())
	tool.Status = StatusAvailable
	tool.ActualReturnDate = returnDate
	tool.ReturnCondition = condition

	for i := range store.Records {
		r := &store.Records[i]
		if r.ToolID == id && r.ActualReturnDate == "" {
			r.ActualReturnDate = returnDate
			r.ReturnCondition = condition
			break
		}
	}

	tool.Borrower = ""
	tool.Phone = ""
	tool.BorrowDate = ""
	tool.ExpectedReturnDate = ""

	if err := SaveData(store); err != nil {
		return err
	}
	fmt.Printf("归还成功: [%s] %s\n", id, tool.Name)
	fmt.Printf("  归还日期: %s\n", returnDate)
	fmt.Printf("  工具状况: %s\n", condition)
	return nil
}

func ListOverdue(_ []string) error {
	store, err := LoadData()
	if err != nil {
		return err
	}

	today := time.Now()
	todayStr := FormatDate(today)
	hasOverdue := false

	fmt.Println("=== 超期未还工具列表 ===")
	for _, tool := range store.Tools {
		if tool.Status == StatusBorrowed && tool.ExpectedReturnDate != "" {
			expected, err := ParseDate(tool.ExpectedReturnDate)
			if err != nil {
				continue
			}
			if today.After(expected) {
				days := int(today.Sub(expected).Hours() / 24)
				if days == 0 {
					days = 1
				}
				hasOverdue = true
				fmt.Printf("  [%s] %s\n", tool.ID, tool.Name)
				fmt.Printf("    借用人: %s\n", tool.Borrower)
				fmt.Printf("    预计归还: %s (今日: %s)\n", tool.ExpectedReturnDate, todayStr)
				fmt.Printf("    超期天数: %d 天\n", days)
				fmt.Println()
			}
		}
	}

	if !hasOverdue {
		fmt.Println("暂无超期未还的工具")
	}
	return nil
}

func MonthlyReport(args []string) error {
	month := ""
	i := 0
	for i < len(args) {
		switch args[i] {
		case "--month":
			if i+1 >= len(args) {
				return fmt.Errorf("--month 缺少参数值")
			}
			month = args[i+1]
			i += 2
		default:
			return fmt.Errorf("未知参数: %s", args[i])
		}
	}

	if month == "" {
		return fmt.Errorf("用法: monthly --month YYYY-MM (如 2024-03)")
	}
	if _, err := time.Parse("2006-01", month); err != nil {
		return fmt.Errorf("月份格式错误，应为 YYYY-MM，如 2024-03")
	}

	store, err := LoadData()
	if err != nil {
		return err
	}

	totalBorrow := 0
	typeCount := make(map[string]int)
	for _, t := range ValidToolTypes {
		typeCount[t] = 0
	}

	for _, r := range store.Records {
		if strings.HasPrefix(r.BorrowDate, month) {
			totalBorrow++
			typeCount[r.ToolType]++
		}
	}

	fmt.Printf("=== %s 月度借出报告 ===\n", month)
	fmt.Printf("  总借出次数: %d 次\n", totalBorrow)
	fmt.Println()
	fmt.Println("  各工具类型分布:")
	if totalBorrow == 0 {
		fmt.Println("    (本月无借出记录)")
	} else {
		for _, t := range ValidToolTypes {
			count := typeCount[t]
			pct := float64(count) / float64(totalBorrow) * 100
			fmt.Printf("    %-8s: %d 次 (%.1f%%)\n", t, count, pct)
		}
	}
	return nil
}
