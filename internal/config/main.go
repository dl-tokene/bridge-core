package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/tokend/bridge/core/internal/data/mem"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/signature"
)

type Config interface {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	signature.Signerer
	mem.Chainer
}

type config struct {
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	signature.Signerer
	mem.Chainer
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Signerer:   signature.NewSignerer(getter),
		Chainer:    mem.NewChainer(getter),
	}
}
