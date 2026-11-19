package enum

type ContentType int16

const (
	ContentTypeNone ContentType = iota
	ContentTypeTree
	ContentTypeSecret
	ContentTypeMonster
)
