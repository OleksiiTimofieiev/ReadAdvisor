package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"ReadAdvisor/cache"
	"ReadAdvisor/clients/telegram"
	"ReadAdvisor/lib/e"
	"ReadAdvisor/storage"

	"mvdan.cc/xurls/v2"

	"time"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	ListCmd  = "/list"
)

var (
	postCache cache.PostCache
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)
	// p.tg.SetMenuCommands()

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case ListCmd:
		return p.sendListOfURL(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Processor) sendListOfURL(chatID int) (err error) {
	urls, _ := p.storage.List(context.Background())

	for _, url := range urls {
		link := url
		go p.tg.SendMessage(chatID, link)
	}
	return nil
}

func (p *Processor) savePeeData(chatID int, username string, quantity int) error {
	//TODO: create separate folder
	year, month, day := time.Now().Date()
	fmt.Println("Year   :", year)
	fmt.Println("Month  :", month)
	fmt.Println("Day    :", day)
	fileName := "files_storage/" + username + "/" + strconv.Itoa(day) + "_" + month.String() + "_" + strconv.Itoa(year)
	fmt.Println(fileName)

	if err := fileExists(fileName); !err {
		fmt.Println("here")

		createFile(fileName)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can`t do command: save page", err) }()

	// sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}

	if isExists {
		if err := p.tg.SendMessage(chatID, msgAlreadyExists); err != nil {
			return err
		}
		return p.storage.Remove(context.Background(), page)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can`t do command: can`t send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	xurlsStrict := xurls.Strict()
	output := xurlsStrict.FindAllString(text, -1)
	if len(output) == 0 {
		return false
	}

	u, err := url.Parse(output[0])

	return err == nil && u.Host != ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created successfully")
	defer file.Close()
}
