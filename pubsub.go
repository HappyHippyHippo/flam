package flam

import "sync"

type PubSub[I ID, C Channel] interface {
	Subscribe(id I, channel C, handler PubSubHandlerType[I, C]) PubSub[I, C]
	Unsubscribe(id I, channel C) PubSub[I, C]
	Publish(channel C, data ...any) error
}

type pubsub[I ID, C Channel] struct {
	locker   sync.Locker
	handlers map[C]map[I]PubSubHandlerType[I, C]
}

func NewPubSub[I ID, C Channel]() PubSub[I, C] {
	return &pubsub[I, C]{
		locker:   &sync.Mutex{},
		handlers: map[C]map[I]PubSubHandlerType[I, C]{},
	}
}

func (pubsub *pubsub[I, C]) Subscribe(
	id I,
	channel C,
	handler PubSubHandlerType[I, C],
) PubSub[I, C] {
	pubsub.locker.Lock()
	defer pubsub.locker.Unlock()

	if _, ok := pubsub.handlers[channel]; !ok {
		pubsub.handlers[channel] = map[I]PubSubHandlerType[I, C]{}
	}

	pubsub.handlers[channel][id] = handler

	return pubsub
}

func (pubsub *pubsub[I, C]) Unsubscribe(
	id I,
	channel C,
) PubSub[I, C] {
	pubsub.locker.Lock()
	defer pubsub.locker.Unlock()

	if subs, ok := pubsub.handlers[channel]; ok {
		delete(subs, id)
	}

	return pubsub
}

func (pubsub *pubsub[I, C]) Publish(
	channel C,
	data ...any,
) error {
	pubsub.locker.Lock()
	defer pubsub.locker.Unlock()

	if subs, ok := pubsub.handlers[channel]; ok {
		for id, handler := range subs {
			if e := handler(id, channel, data...); e != nil {
				return e
			}
		}
	}

	return nil
}
