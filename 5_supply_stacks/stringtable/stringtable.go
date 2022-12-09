package stringtable

func ParseTableRow(line string) []byte {
	out := []byte{}

	for len(line) > 0 {
		if line[0] == '[' {
			out = append(out, line[1])
		}

		if line[0] == ' ' {
			out = append(out, 'x')
		}

		if line = line[3:]; len(line) == 0 {
			break
		}

		line = line[1:]
	}

	return out
}
