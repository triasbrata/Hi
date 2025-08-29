package impl

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/internals/entities"
	"github.com/triasbrata/adios/internals/repositories/words"
	"github.com/triasbrata/adios/internals/service/hello"
	"github.com/triasbrata/adios/pkgs/messagebroker/publisher"
)

type srv struct {
	cfg  *config.Config
	repo words.WordRepository
	pub  publisher.Publisher
	at   atomic.Int64
}

// FetchHelloWorld implements hello.HelloService.
func (s *srv) FetchHelloWorld(ctx context.Context, param entities.FetchHelloWorldParam) (res entities.FetchHelloWorldRes, err error) {
	res.MapData = make(map[string][]string)
	dataWords, err := s.repo.GetWord(ctx, entities.GetWordParam{})
	if err != nil {
		return res, err
	}
	res.MapData["hello"] = dataWords.Data
	start := time.Now()
	err = s.pub.PublishToQueue(ctx, "test_consumer", publisher.PublishPayload{
		Body: []byte("{'hello':'world'}"),
	})
	fmt.Printf("time.Since(start).Milliseconds(): %v\n", time.Since(start).Milliseconds())
	s.at.Add(1)
	fmt.Printf("s.at.Load(): %v\n", s.at.Load())

	return res, nil
}

func NewServiceHello(cfg *config.Config, repo words.WordRepository, pub publisher.Publisher) hello.HelloService {
	return &srv{repo: repo, cfg: cfg, pub: pub}
}
