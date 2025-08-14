package flam

type ResourceCreator[R Resource] interface {
	Accept(config Bag) bool
	Create(config Bag) (R, error)
}
