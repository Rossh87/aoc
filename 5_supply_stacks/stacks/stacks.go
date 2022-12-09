package stacks

type Stacks [][]byte

func (s *Stacks) Reverse() {
	for _, stack := range *s {
		low := 0
		high := len(stack) - 1

		for low < high {
			stack[low], stack[high] = stack[high], stack[low]
			low++
			high--
		}
	}
}

func (s *Stacks) AddRow(row []byte) {
	for idx, b := range row {
		if b == 'x' {
			continue
		}

		(*s)[idx] = append((*s)[idx], b)
	}
}

func (s *Stacks) Move(count, src, dest int) {
	ss := (*s)[src]
	ds := (*s)[dest]

	for count > 0 {
		ds = append(ds, ss[len(ss)-1])
		ss = ss[:len(ss)-1]
		count--
	}

	(*s)[src] = ss
	(*s)[dest] = ds
}

func (s *Stacks) MoveMultiple(count, src, dest int) {
	ss := (*s)[src]
	ds := (*s)[dest]
	cutOff := len(ss) - count
	toAppend := ss[cutOff:]

	ds = append(ds, toAppend...)

	(*s)[src] = ss[:cutOff]
	(*s)[dest] = ds
}

// trick is we can always bite off 3, then bite off the gap.
func New(line string) *Stacks {
	s := make(Stacks, 0)

	pos := 0

	for len(line) > 0 {
		if len(s) < pos+1 {
			s = append(s, []byte{})
		}

		if line = line[3:]; len(line) == 0 {
			break
		}

		line = line[1:]
		pos++
	}

	return &s
}
