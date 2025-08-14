package flam

type FactoryConfig interface {
	Get(path string, def ...any) Bag
}
