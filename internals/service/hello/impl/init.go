package impl

import (
	"context"

	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/entities"
	"github.com/triasbrata/adios/internals/repositories/words"
	"github.com/triasbrata/adios/internals/service/hello"
)

type srv struct {
	cfg  *config.Config
	repo words.WordRepository
}

// FetchHelloWorld implements hello.HelloService.
func (s *srv) FetchHelloWorld(ctx context.Context, param entities.FetchHelloWorldParam) (res entities.FetchHelloWorldRes, err error) {
	res.MapData = make(map[string][]string)
	dataWords, err := s.repo.GetWord(ctx, entities.GetWordParam{})
	if err != nil {
		return res, err
	}
	res.MapData["hello"] = dataWords.Data
	return res, nil
}

func NewServiceHello(cfg *config.Config, repo words.WordRepository) hello.HelloService {
	return &srv{repo: repo, cfg: cfg}
}
