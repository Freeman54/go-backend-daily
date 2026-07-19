package fieldmaskmerge

import "fmt"

func Apply(base map[string]string, updates map[string]*string, mask []string, allowed map[string]struct{}) (map[string]string, error) {
	next := make(map[string]string, len(base))
	for key, value := range base {
		next[key] = value
	}

	for _, path := range mask {
		if _, ok := allowed[path]; !ok {
			return nil, fmt.Errorf("field %q is not allowed", path)
		}

		value, exists := updates[path]
		if !exists {
			return nil, fmt.Errorf("field %q is missing from updates", path)
		}

		if value == nil {
			delete(next, path)
			continue
		}
		next[path] = *value
	}

	return next, nil
}
