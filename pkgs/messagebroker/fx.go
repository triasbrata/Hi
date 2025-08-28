package messagebroker

import (
	"github.com/triasbrata/adios/internals/config"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker"
	"github.com/triasbrata/adios/pkgs/messagebroker/broker/impl"
	"go.uber.org/fx"
)

func LoadMessageBrokerAmqp(cfg ...impl.AmqpConfig) fx.Option {
	provider := []fx.Option{}
	if len(cfg) == 1 {
		provider = append(provider, fx.Supply(cfg[0], fx.Private))
	}
	if len(cfg) == 0 {
		provider = append(provider, fx.Provide(func(c *config.Config) impl.AmqpConfig {
			return impl.AmqpConfig{
				ConnectionName: c.Consumer.Amqp.ConnectionName,
				URI:            c.Consumer.Amqp.URI,
			}
		}, fx.Private))
	}
	provider = append(provider, fx.Provide(func(cfg impl.AmqpConfig) (broker.Broker, error) {
		return impl.CreateNewBroker(impl.WithAmqpBroker(cfg))

	}))
	return fx.Module("pkg/messagebroker", provider...)

}
