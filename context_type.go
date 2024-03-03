package go_cake

type ContextType int

const (
	ctxDbDriverTotal ContextType = iota
	ctxDbDriverFind
	ctxDbDriverInsert
	ctxDbDriverDelete
	ctxDbDriverUpdate
)
