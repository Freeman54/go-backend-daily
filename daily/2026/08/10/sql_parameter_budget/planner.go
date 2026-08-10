// Package parametbudget plans SQL batches within a parameter-count limit.
package parametbudget

import "errors"

var ErrInvalidBudget = errors.New("invalid SQL parameter budget")

type Planner struct{ MaxParameters int }

// BatchSizes returns row counts for statements with fixed parameters and parametersPerRow.
func (p Planner) BatchSizes(rows, parametersPerRow, fixed int) ([]int, error) {
	if p.MaxParameters <= 0 || rows < 0 || parametersPerRow <= 0 || fixed < 0 || fixed >= p.MaxParameters {
		return nil, ErrInvalidBudget
	}
	perBatch := (p.MaxParameters - fixed) / parametersPerRow
	if perBatch == 0 {
		return nil, ErrInvalidBudget
	}
	if rows == 0 {
		return []int{}, nil
	}
	sizes := make([]int, 0, (rows+perBatch-1)/perBatch)
	for rows > 0 {
		size := min(rows, perBatch)
		sizes = append(sizes, size)
		rows -= size
	}
	return sizes, nil
}
