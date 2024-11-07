package queue

type Queue interface {
	Enqueue(string)
	Dequeue() (int64, string)

	Close()
}
