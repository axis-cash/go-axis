package consensus

import (
	"github.com/axis-cash/go-axis/axisdb"
)

type DB interface {
	CurrentTri() axisdb.Tri
	GlobalGetter() axisdb.Getter
}

type CItem interface {
	CopyTo() (ret CItem)
	CopyFrom(CItem)
}

type PItem interface {
	CItem
	Id() (ret []byte)
	State() (ret []byte)
}
