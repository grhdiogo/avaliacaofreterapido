package quote

type Repository interface {
	Store(e Entity) error
	List(limit int) ([]Entity, error)
}
