package telegram

import (
	"context"
	"errors"

	"URLReminderBot/clients/telegram"
	"URLReminderBot/events"
	"URLReminderBot/libraries/myErr"
	"URLReminderBot/storage"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(context.TODO(), p.offset, limit)
	if err != nil {
		return nil, myErr.Wrap("can't get updates: ", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, updatesToEvents(u))

	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = myErr.Wrap("can't process the message: ", err) }()
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err = p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}

	return res, nil
}

func updatesToEvents(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message

}
