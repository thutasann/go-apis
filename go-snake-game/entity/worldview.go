package entity

type worldView interface {
	GetEntities(tag string) []Entity
	GetFirstEntity(tag string) (Entity, bool)
}
