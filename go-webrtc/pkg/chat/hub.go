package chat

type Hub struct {
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) Run() {}
