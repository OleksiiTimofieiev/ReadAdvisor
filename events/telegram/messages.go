package telegram

const msgHelp = `In order to save the page, just send me a link to it.
In order to get a random page from your list,
send me command /rnd.
Caution! `

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages="You have no save pages"
	msgSaved = "Saved! 😘 "
	msgAlreadyExists = "You have already have this page in your list. URL will be deleted from DB"
	msgDeleted = "URL was deleted"
)