package sqlsetdiff

import "sort"

type Diff struct {
	ToInsert []int64
	ToDelete []int64
}

func Build(current []int64, desired []int64) Diff {
	currentSet := uniqueSet(current)
	desiredSet := uniqueSet(desired)

	diff := Diff{
		ToInsert: make([]int64, 0),
		ToDelete: make([]int64, 0),
	}

	for id := range desiredSet {
		if !currentSet[id] {
			diff.ToInsert = append(diff.ToInsert, id)
		}
	}
	for id := range currentSet {
		if !desiredSet[id] {
			diff.ToDelete = append(diff.ToDelete, id)
		}
	}

	sort.Slice(diff.ToInsert, func(i, j int) bool { return diff.ToInsert[i] < diff.ToInsert[j] })
	sort.Slice(diff.ToDelete, func(i, j int) bool { return diff.ToDelete[i] < diff.ToDelete[j] })
	return diff
}

func uniqueSet(ids []int64) map[int64]bool {
	out := make(map[int64]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out
}
