package insertbatchplan

import "fmt"

type Plan struct {
	Start int
	End   int
}

func Split(totalRows int, columns int, maxRows int, maxParams int) ([]Plan, error) {
	if totalRows < 0 {
		return nil, fmt.Errorf("total rows must not be negative")
	}
	if columns <= 0 {
		return nil, fmt.Errorf("columns must be positive")
	}
	if maxRows <= 0 {
		return nil, fmt.Errorf("max rows must be positive")
	}
	if maxParams < columns {
		return nil, fmt.Errorf("max params must be at least columns")
	}
	if totalRows == 0 {
		return nil, nil
	}

	rowsPerBatch := maxRows
	if paramLimited := maxParams / columns; paramLimited < rowsPerBatch {
		rowsPerBatch = paramLimited
	}
	if rowsPerBatch <= 0 {
		return nil, fmt.Errorf("rows per batch must be positive")
	}

	var plans []Plan
	for start := 0; start < totalRows; start += rowsPerBatch {
		end := start + rowsPerBatch
		if end > totalRows {
			end = totalRows
		}
		plans = append(plans, Plan{Start: start, End: end})
	}
	return plans, nil
}
