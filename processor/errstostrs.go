package processor

func errsToStrs(es []error) []string {
	acc := make([]string, len(es))
	for i, e := range es {
		acc[i] = e.Error()
	}
	return acc
}
