package impl

import (
	"context"

	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/entities"
	"github.com/triasbrata/adios/internals/repositories/words"
)

type repo struct {
	cfg *config.Config
}

// GetWord implements words.Repository.
func (r *repo) GetWord(ctx context.Context, param entities.GetWordParam) (res entities.GetWordRes, err error) {
	return entities.GetWordRes{
		Data: []string{"hello", "world"},
	}, nil
}

func NewWordRepository(cfg *config.Config) words.WordRepository {
	return &repo{cfg: cfg}
}
