package flam

type PubSubHandlerType[I ID, C Channel] func(channel C, data ...any) error
