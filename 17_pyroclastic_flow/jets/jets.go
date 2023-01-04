package jets

type Push rune

const (
	Left  Push = '<'
	Right Push = '>'
)

type Jets struct {
	directions []Push
	currByte   int
}

func New(rawDirs string) *Jets {
	dirSlice := []Push(rawDirs)
	js := Jets{dirSlice, 0}
	return &js
}

func (j *Jets) Next() Push {
	dir := j.directions[j.currByte]
	j.currByte = (j.currByte + 1) % len(j.directions)
	return dir
}
