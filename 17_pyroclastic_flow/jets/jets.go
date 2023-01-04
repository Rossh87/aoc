package jets

type Push rune

const (
	Left  Push = '<'
	Right Push = '>'
)

type Jets struct {
	directions []Push
	nextPush   int
}

func New(rawDirs string) *Jets {
	dirSlice := []Push(rawDirs)
	js := Jets{dirSlice, 0}
	return &js
}

func (j *Jets) Next() Push {
	dir := j.directions[j.nextPush]
	j.nextPush++
	j.nextPush %= len(j.directions)
	return dir
}

func (j Jets) Peek() int {
	return j.nextPush
}
