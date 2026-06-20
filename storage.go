package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const dataFile = "tools_data.json"

const (
	ToolTypeElectric  = "电动工具"
	ToolTypeManual    = "手动工具"
	ToolTypeMeasure   = "测量工具"
	ToolTypeGardening = "园艺工具"
)

const (
	StatusAvailable = "在库"
	StatusBorrowed  = "借出中"
)

const (
	ConditionGood       = "完好"
	ConditionMinor      = "轻微损耗"
	ConditionDamaged    = "有损坏"
)

var ValidToolTypes = []string{ToolTypeElectric, ToolTypeManual, ToolTypeMeasure, ToolTypeGardening}
var ValidConditions = []string{ConditionGood, ConditionMinor, ConditionDamaged}
var ValidStatuses = []string{StatusAvailable, StatusBorrowed}

type Tool struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Deposit  int    `json:"deposit"`
	Status   string `json:"status"`
	Borrower string `json:"borrower,omitempty"`
	Phone    string `json:"phone,omitempty"`
	BorrowDate string `json:"borrow_date,omitempty"`
	ExpectedReturnDate string `json:"expected_return_date,omitempty"`
	ActualReturnDate   string `json:"actual_return_date,omitempty"`
	ReturnCondition    string `json:"return_condition,omitempty"`
}

type BorrowRecord struct {
	ToolID             string `json:"tool_id"`
	ToolName           string `json:"tool_name"`
	ToolType           string `json:"tool_type"`
	Borrower           string `json:"borrower"`
	Phone              string `json:"phone"`
	BorrowDate         string `json:"borrow_date"`
	ExpectedReturnDate string `json:"expected_return_date"`
	ActualReturnDate   string `json:"actual_return_date,omitempty"`
	ReturnCondition    string `json:"return_condition,omitempty"`
}

type DataStore struct {
	Tools   []Tool         `json:"tools"`
	Records []BorrowRecord `json:"records"`
}

func IsValidToolType(t string) bool {
	for _, vt := range ValidToolTypes {
		if vt == t {
			return true
		}
	}
	return false
}

func IsValidCondition(c string) bool {
	for _, vc := range ValidConditions {
		if vc == c {
			return true
		}
	}
	return false
}

func IsValidStatus(s string) bool {
	for _, vs := range ValidStatuses {
		if vs == s {
			return true
		}
	}
	return false
}

func validateStore(store *DataStore) error {
	for i, tool := range store.Tools {
		if !IsValidToolType(tool.Type) {
			return fmt.Errorf("工具 [%s] 类型无效: %q，可选值: %v", tool.ID, tool.Type, ValidToolTypes)
		}
		if !IsValidStatus(tool.Status) {
			return fmt.Errorf("工具 [%s] 状态无效: %q，可选值: %v", tool.ID, tool.Status, ValidStatuses)
		}
		if tool.ReturnCondition != "" && !IsValidCondition(tool.ReturnCondition) {
			return fmt.Errorf("工具 [%s] 归还状况无效: %q，可选值: %v", tool.ID, tool.ReturnCondition, ValidConditions)
		}
		_ = i
	}
	for i, record := range store.Records {
		if !IsValidToolType(record.ToolType) {
			return fmt.Errorf("借阅记录[%d] 工具类型无效: %q，可选值: %v", i, record.ToolType, ValidToolTypes)
		}
		if record.ReturnCondition != "" && !IsValidCondition(record.ReturnCondition) {
			return fmt.Errorf("借阅记录[%d] 归还状况无效: %q，可选值: %v", i, record.ReturnCondition, ValidConditions)
		}
	}
	return nil
}

func LoadData() (*DataStore, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return &DataStore{
			Tools:   []Tool{},
			Records: []BorrowRecord{},
		}, nil
	}
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, fmt.Errorf("读取数据文件失败: %w", err)
	}
	var store DataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("解析数据文件失败: %w", err)
	}
	if err := validateStore(&store); err != nil {
		return nil, fmt.Errorf("数据校验失败: %w", err)
	}
	return &store, nil
}

func SaveData(store *DataStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		return fmt.Errorf("写入数据文件失败: %w", err)
	}
	return nil
}

func FindTool(store *DataStore, id string) *Tool {
	for i := range store.Tools {
		if store.Tools[i].ID == id {
			return &store.Tools[i]
		}
	}
	return nil
}

func ParseDate(s string) (time.Time, error) {
	loc := time.UTC
	return time.ParseInLocation("2006-01-02", s, loc)
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
