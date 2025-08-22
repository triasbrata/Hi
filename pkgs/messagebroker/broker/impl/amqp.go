package impl

type AmqpConfigTls struct {
}
type AmqpConfig struct {
	ConnectionName string
	URI            string
	TLS            AmqpConfigTls
}

func WithAmqpBroker(config AmqpConfig) brokerConfig {
	return func(brk *brk) {
		brk.cfg.amqp = &brokerConAMQP{
			name: config.ConnectionName,
			uri:  config.URI,
		}
	}
}
