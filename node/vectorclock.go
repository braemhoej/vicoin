package node

type VectorClock struct {
	vector map[Peer]int
}

func NewVectorClock() VectorClock {
	return VectorClock{
		vector: make(map[Peer]int),
	}
}

func (vectorClock *VectorClock) Get(peer Peer) int {
	return vectorClock.vector[peer]
}

func (vectorClock *VectorClock) Increment(peer Peer) {
	vectorClock.vector[peer] += 1
}

func Update(v1 VectorClock, v2 VectorClock) VectorClock {
	updatedClock := v1
	for peer, time := range v2.vector {
		updatedClock.vector[peer] = max(time, v2.vector[peer])
	}
	return updatedClock
}

func max(i1 int, i2 int) int {
	if i1 >= i2 {
		return i1
	} else {
		return i2
	}
}
