package libs

func Enque(queue []string, element string) []string {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func Dequeue(queue []string) (string, []string) {
	first := queue[0]
	if len(queue) == 1 {
		var tmp = []string{}
		return first, tmp
	}

	return first, queue[1:]
}
