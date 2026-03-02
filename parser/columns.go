package parser

// buildColumns converts normalized headers into column metadata.
// It is intentionally dumb: no SQL rules, no sanitization.
func buildColumns(headers []string) []Column {
	cols := make([]Column, 0, len(headers))

	for _, h := range headers {
		if h == "" {
			// still keep placeholder; downstream decides what to do
			cols = append(cols, Column{
				Name: "",
				Type: "STRING",
			})
			continue
		}

		cols = append(cols, Column{
			Name: h,
			Type: "STRING", // logical type, not SQL type
		})
	}

	return cols
}
