package counter

type PageViewCounter interface {
	Increment(pageID string) error
	GetCount(pageID string) (int, error)
}

type Counter struct {
	// Поля для хранилища (например, Redis, in-memory map, DB)
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Increment(pageID string) error {
	return nil
}

func (c *Counter) GetCount(pageID string) (int, error) {
	return 0, nil
}
