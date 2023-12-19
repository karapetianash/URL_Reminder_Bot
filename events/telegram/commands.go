package telegram

import (
	"URLReminderBot/clients/telegram"
	"URLReminderBot/libraries/myErr"
	"URLReminderBot/storage"
	"context"
	"errors"
	"log"
	url2 "net/url"
	"strings"
)

const (
	RandCmd  = "/rnd"
	StartCmd = "/start"
	HelpCmd  = "/help"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.doCmd(text, chatID, username)
	}

	switch text {
	case RandCmd:
		return p.sendRandom(chatID, username)
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(context.TODO(), chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = myErr.WrapIfErr("can't save page: ", err) }()

	sendMsg := newMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}
	if isExist {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}

	if err = p.tg.SendMessage(context.TODO(), chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = myErr.WrapIfErr("can't send random: ", err) }()

	sendMsg := newMessageSender(chatId, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if err = sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	sendMsg := newMessageSender(chatID, p.tg)
	return sendMsg(msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	sendMsg := newMessageSender(chatID, p.tg)
	return sendMsg(msgHello)
}

func newMessageSender(chatID int, client *telegram.Client) func(string) error {
	return func(msg string) error {
		return client.SendMessage(context.TODO(), chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

// TODO: upgrade this function
func isURL(text string) bool {
	url, err := url2.Parse(text)

	return err != nil && url.Host != ""
}
