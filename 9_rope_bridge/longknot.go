package main

// knots[0] = head, knots[9] = tail
type longknot struct {
	knots []point
}

func (k *longknot) updateLongTail() {
	knotSlice := (*k).knots

	for i := 1; i < 10; i++ {
		sectionHead := knotSlice[i-1]
		sectionTail := knotSlice[i]

		newp := tailPos(sectionHead, sectionTail)

		knotSlice[i] = newp
	}

	// (*k).knots = knotSlice
}

func (k *longknot) up() point {
	// update knots[0]
	k.knots[0].y++
	k.updateLongTail()
	return k.knots[len(k.knots)-1]
}

func (k *longknot) down() point {
	// update knots[0]
	k.knots[0].y--
	k.updateLongTail()
	return k.knots[len(k.knots)-1]

}

func (k *longknot) left() point {
	// update knots[0]
	k.knots[0].x--
	k.updateLongTail()
	return k.knots[len(k.knots)-1]
}

func (k *longknot) right() point {
	// update knots[0]
	k.knots[0].x++
	k.updateLongTail()
	return k.knots[len(k.knots)-1]
}

func (k longknot) tail() point {
	return k.knots[len(k.knots)-1]
}

func newLongKnot() longknot {
	return longknot{make([]point, 10)}
}
