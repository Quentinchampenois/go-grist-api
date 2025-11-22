package grist

import (
	"encoding/json"
	"fmt"
)

type Columns struct {
	Columns []Column
}

type Column struct {
	ID     string               `json:"id,omitempty"`
	Label  string               `json:"label,omitempty"`
	Type   string               `json:"type,omitempty"`
	Fields map[string]CellValue `json:"fields,omitempty"`
}

type Field struct {
	Label string `json:"label"`
	Type  []byte
}

// CellValue represents a cell value in a Grist table.
// source: https://github.com/gristlabs/grist-core/blob/9ba8a2e30bebd329ca447e1311de7af2cbb0efe5/app/plugin/GristData.ts#L46
type CellValue struct {
	Number  *float64
	String  *string
	Boolean *bool
	Null    bool
	Object  *ObjectGrist
}

// ObjectGrist represents a Grist object
// source: https://support.getgrist.com/code/enums/GristData.GristObjCode/
type ObjectGrist struct {
	Code string
	Data []any
}

func (c *CellValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		c.Null = true
		return nil
	}

	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		c.Number = &num
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		c.String = &str
		return nil
	}

	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		c.Boolean = &b
		return nil
	}

	var arr []any
	if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
		if code, ok := arr[0].(string); ok {
			c.Object = &ObjectGrist{
				Code: code,
				Data: arr[1:],
			}
			return nil
		}
	}

	return fmt.Errorf("unknown CellValue type: %s", string(data))
}

func (c *CellValue) MarshalJSON() ([]byte, error) {
	fmt.Println("MarshalJSON")
	switch {
	case c.Number != nil:
		return json.Marshal(*c.Number)
	case c.String != nil:
		return json.Marshal(*c.String)
	case c.Boolean != nil:
		return json.Marshal(*c.Boolean)
	case c.Null:
		return []byte("null"), nil
	case c.Object != nil:
		arr := []any{c.Object.Code}
		arr = append(arr, c.Object.Data...)
		return json.Marshal(arr)
	default:
		return nil, fmt.Errorf("empty CellValue")
	}
}
