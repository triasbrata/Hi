package words

import (
	"context"

	"github.com/triasbrata/adios/internals/entities"
)

type WordRepository interface {
	GetWord(ctx context.Context, param entities.GetWordParam) (res entities.GetWordRes, err error)
}
