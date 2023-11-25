package main

import (
	// "fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Contact struct {
	firstName	string
	lastName	string
	email		string
	phoneNumber string
	country		string
	business	bool
}

// I don't like too much doing this
var countries = []string{"AR", "BR", "BO", "CL", "PA"}
var contacts = make([]Contact, 0)
var contactsList = tview.NewList().ShowSecondaryText(false)
var contactText = tview.NewTextView()

func selectContact() {
	contactsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		setConcatText(&contacts[index])
	})
}

// Is this the best way? I don't think so, we could pass the contact as a parameter
func addContactTolist() {
	contactsList.Clear()

	for idx, contact := range contacts {
		contactsList.AddItem(contact.firstName + " " + contact.lastName, " ", rune(49 + idx), nil)
	}
}

func addContact(form *tview.Form, pages *tview.Pages, app *tview.Application) {
	contact := Contact{}

	form.AddInputField("First Name", "", 20, nil, func(firstName string) {
		contact.firstName = firstName
	})
	form.AddInputField("Last Name", "", 20, nil, func(lastname string){
		contact.lastName = lastname
	})
	form.AddInputField("Email", "", 30, nil, func(email string){
		contact.email = email
	})
	form.AddInputField("Phone number", "", 30, nil, func(phone string){
		contact.phoneNumber = phone
	})
	form.AddDropDown("Country", countries, 0, func(country string, idx int){
		contact.country = country
	})
	form.AddCheckbox("Business", false, func(business bool) {
		contact.business = business
	})

	form.AddButton("Save", func() {
		contacts = append(contacts, contact)
		addContactTolist()
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})

}

func setConcatText(contact *Contact) {
	contactText.Clear()
	text := contact.firstName + " " + contact.lastName + "\n" + contact.email + "\n" + contact.phoneNumber
	contactText.SetText(text)
}

func main() {
	var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contact \n(q) to quit")
	flexView := tview.NewFlex()
	flexView.SetDirection(tview.FlexRow).
			AddItem(tview.NewFlex().
			AddItem(contactsList, 0, 1, true).
			AddItem(contactText, 0, 4, true), 0, 6, false).
		AddItem(text, 0, 1, false)

	flexView.SetBorder(true).SetTitle("Simple contacts")

	form := tview.NewForm()
	form.SetTitle("Add new contact").
		SetBorder(true)
	pages := tview.NewPages()

	pages.AddPage("Menu", flexView, true, true)
	pages.AddPage("Add Contact", form, true, false)

	currItem := 0
	app := tview.NewApplication().SetRoot(pages, true).EnableMouse(true)
	flexView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		} else if event.Rune() == 'a' {
			form.Clear(true)
			addContact(form, pages, app)
			pages.SwitchToPage("Add Contact")
		} else if event.Key() == tcell.KeyUp {
			currItem -= 1
			currItem = currItem % contactsList.GetItemCount()
			contactsList.SetCurrentItem(currItem)
		} else if event.Key() == tcell.KeyDown {
			currItem += 1
			currItem = currItem % contactsList.GetItemCount()
			contactsList.SetCurrentItem(currItem)
		}
		return event
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
