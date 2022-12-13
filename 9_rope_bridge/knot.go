package main

type knot struct {
	h point
	t point
}

func (k *knot) up() point {
	(*k).h.y++

	newTail := tailPos((*k).h, (*k).t)

	(*k).t = newTail

	return newTail
}

func (k *knot) down() point {
	(*k).h.y--

	newTail := tailPos((*k).h, (*k).t)

	(*k).t = newTail

	return newTail

}

func (k *knot) left() point {
	(*k).h.x--

	newTail := tailPos((*k).h, (*k).t)

	(*k).t = newTail

	return newTail
}

func (k *knot) right() point {
	(*k).h.x++

	newTail := tailPos((*k).h, (*k).t)

	(*k).t = newTail

	return newTail
}

func (k *knot) tail() point {
	return k.t
}

func newKnot() knot {
	h := point{0, 0}
	t := point{0, 0}

	return knot{h, t}
}
