package spine

import (
	"github.com/enbility/eebus-go/spine/model"
)

type Operations struct {
	Read, Write bool
}

func NewOperations(read, write bool) *Operations {
	return &Operations{
		Read:  read,
		Write: write,
	}
}

func (r *Operations) String() string {
	switch {
	case r.Read && !r.Write:
		return "RO"
	case r.Read && r.Write:
		return "RW"
	default:
		return "--"
	}
}

func (r *Operations) Information() *model.PossibleOperationsType {
	res := new(model.PossibleOperationsType)
	if r.Read {
		res.Read = &model.PossibleOperationsReadType{}
	}
	if r.Write {
		res.Write = &model.PossibleOperationsWriteType{}
	}

	return res
}
