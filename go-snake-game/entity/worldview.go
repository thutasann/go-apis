package entity

type worldView interface {
	GetEntities(tag string) []Entity
}
