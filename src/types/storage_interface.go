package types

type IStorage interface {
	DatabaseHas(table, column string, value any) bool
}