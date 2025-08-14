package flam

type PubSubHandlerType[I ID, C Channel] func(id I, channel C, data ...any) error
