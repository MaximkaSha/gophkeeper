package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/MaximkaSha/gophkeeper/internal/client"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var logo = `   

░██████╗░░█████╗░██████╗░██╗░░██╗██╗░░██╗███████╗███████╗██████╗░███████╗██████╗░
██╔════╝░██╔══██╗██╔══██╗██║░░██║██║░██╔╝██╔════╝██╔════╝██╔══██╗██╔════╝██╔══██╗
██║░░██╗░██║░░██║██████╔╝███████║█████═╝░█████╗░░█████╗░░██████╔╝█████╗░░██████╔╝
██║░░╚██╗██║░░██║██╔═══╝░██╔══██║██╔═██╗░██╔══╝░░██╔══╝░░██╔═══╝░██╔══╝░░██╔══██╗
╚██████╔╝╚█████╔╝██║░░░░░██║░░██║██║░╚██╗███████╗███████╗██║░░░░░███████╗██║░░██║
░╚═════╝░░╚════╝░╚═╝░░░░░╚═╝░░╚═╝╚═╝░░╚═╝╚══════╝╚══════╝╚═╝░░░░░╚══════╝╚═╝░░╚═╝    
  `

func UI(ctx context.Context, client client.Client) {
	user := models.User{}
	app := tview.NewApplication()
	//box := tview.NewBox().SetBorder(true).SetTitle("GOPHKEEPER")
	logo := tview.NewTextView().SetText(logo)
	status := tview.NewTextView()
	form := tview.NewForm().
		AddInputField("Email: ", user.Email, 20, nil, func(text string) {
			user.Email = text
		}).
		AddPasswordField("Password:", user.Password, 20, '*', func(text string) {
			user.Password = text
		}).
		AddButton("Login", func() {
			err := client.UserLogin(ctx, user)
			if err != nil {
				status.SetText(err.Error())
			} else {
				status.SetText("Logged In")
				app.Stop()
				loggedIn(ctx, client)
			}
		}).
		AddButton("Register", func() {
			err := client.UserRegister(ctx, user)
			if err != nil {
				status.SetText(err.Error())
			} else {
				status.SetText("Registred")
			}
		})
	grid := tview.NewGrid().
		AddItem(logo, 1, 1, 1, 3, 0, 0, false).
		AddItem(form, 2, 1, 1, 3, 0, 0, true).
		AddItem(status, 3, 1, 1, 3, 0, 0, false)
	err := app.SetRoot(grid, true).Run()
	if err != nil {
		panic(err)
	}
}

func DrawError(err error) {

}

func UpdateTable(ctx context.Context, client client.Client, table *tview.Table) *tview.Table {
	for i := 0; i < len(client.LocalStorage.PasswordStorage); i++ {
		table.SetCellSimple(i+1, 0, client.LocalStorage.PasswordStorage[i].Tag)
	}
	for i := 0; i < len(client.LocalStorage.CCStorage); i++ {
		table.SetCellSimple(i+1, 1, client.LocalStorage.CCStorage[i].Tag)
	}
	for i := 0; i < len(client.LocalStorage.TextStorage); i++ {
		table.SetCellSimple(i+1, 2, client.LocalStorage.TextStorage[i].Tag)
	}
	for i := 0; i < len(client.LocalStorage.DataStorage); i++ {
		table.SetCellSimple(i+1, 3, client.LocalStorage.DataStorage[i].Tag)
	}
	return table
}

func DrawEditFromTable(ctx context.Context, client client.Client, app *tview.Application, grid *tview.Grid, row int, dType int) *tview.Form {
	form := tview.NewForm()
	switch dType {
	case 1:
		isChanged := false
		cc := client.LocalStorage.CCStorage[row-1]
		form.AddInputField("CCNum: ", cc.CardNum, 16, nil, func(text string) {
			isChanged = true
			cc.CardNum = text
		}).
			AddInputField("Exp: ", cc.Exp, 16, nil, func(text string) {
				isChanged = true
				cc.Exp = text
			}).
			AddInputField("Name: ", cc.Name, 16, nil, func(text string) {
				isChanged = true
				cc.Name = text
			}).
			AddInputField("CVV: ", cc.CVV, 16, nil, func(text string) {
				isChanged = true
				cc.CVV = text
			}).
			AddInputField("Tag: ", cc.Tag, 16, nil, func(text string) {
				isChanged = true
				cc.Tag = text
			}).
			AddButton("Add/Update", func() {
				if isChanged {
					err := client.AddData(ctx, cc)
					if err != nil {
						DrawError(err)
					}
				}
				app.Stop()
				loggedIn(ctx, client)
			}).AddButton("Del", func() {
			client.DelData(ctx, cc)
			app.Stop()
			loggedIn(ctx, client)
		})
	case 0:
		isChanged := false
		pass := client.LocalStorage.PasswordStorage[row-1]
		form.AddInputField("Login: ", pass.Login, 16, nil, func(text string) {
			isChanged = true
			pass.Login = text
		}).
			AddInputField("Password: ", pass.Password, 16, nil, func(text string) {
				isChanged = true
				pass.Password = text
			}).
			AddInputField("Tag: ", pass.Tag, 16, nil, func(text string) {
				isChanged = true
				pass.Login = text
			}).
			AddButton("Add/Update", func() {
				if isChanged {
					err := client.AddData(ctx, pass)
					if err != nil {
						DrawError(err)
					}
				}
				app.Stop()
				loggedIn(ctx, client)
			}).AddButton("Del", func() {
			client.DelData(ctx, pass)
			app.Stop()
			loggedIn(ctx, client)
		})
	case 2:
		isChanged := false
		txt := client.LocalStorage.TextStorage[row-1]
		form.AddInputField("Text: ", txt.Data, 16, nil, func(text string) {
			isChanged = true
			txt.Data = text
		}).
			AddInputField("Tag: ", txt.Tag, 16, nil, func(text string) {
				isChanged = true
				txt.Tag = text
			}).
			AddButton("Add/Update", func() {
				if isChanged {
					err := client.AddData(ctx, txt)
					if err != nil {
						DrawError(err)
					}
				}
				app.Stop()
				loggedIn(ctx, client)
			}).AddButton("Del", func() {
			client.DelData(ctx, txt)
			app.Stop()
			loggedIn(ctx, client)
		})
	case 3:
		isChanged := false
		data := client.LocalStorage.DataStorage[row-1]
		form.AddInputField("Data: ", string(data.Data), 16, nil, func(text string) {
			isChanged = true
			data.Data = []byte(text)
		}).
			AddInputField("Tag: ", data.Tag, 16, nil, func(text string) {
				isChanged = true
				data.Tag = text
			}).
			AddButton("Add/Update", func() {
				if isChanged {
					err := client.AddData(ctx, data)
					if err != nil {
						DrawError(err)
					}
				}
				app.Stop()
				loggedIn(ctx, client)
			}).AddButton("Del", func() {
			client.DelData(ctx, data)
			app.Stop()
			loggedIn(ctx, client)
		})
		//return form

	}
	return form
}

func loggedIn(ctx context.Context, client client.Client) {

	app := tview.NewApplication()
	err := client.GetAllDataFromDB(ctx)
	if err != nil {
		log.Println(err)
	}

	status := tview.NewTextView().SetText("Ctrl + (A)dd,  (E)xit")
	table := tview.NewTable()
	table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("Passwords(%v)", len(client.LocalStorage.PasswordStorage))).SetExpansion(1).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color100))
	table.SetCell(0, 1, tview.NewTableCell(fmt.Sprintf("Credit Cards(%v)", len(client.LocalStorage.CCStorage))).SetExpansion(1).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color100))
	table.SetCell(0, 2, tview.NewTableCell(fmt.Sprintf("Texts(%v)", len(client.LocalStorage.TextStorage))).SetExpansion(1).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color100))
	table.SetCell(0, 3, tview.NewTableCell(fmt.Sprintf("Data(%v)", len(client.LocalStorage.DataStorage))).SetExpansion(1).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color100))
	table = UpdateTable(ctx, client, table)
	table.SetSelectable(true, true)

	grid := tview.NewGrid().
		AddItem(table, 0, 0, 1, 1, 0, 0, false).
		AddItem(status, 1, 0, 1, 1, 0, 0, false)
	table.SetSelectedFunc(func(row int, column int) {
		if row > 0 {
			form := DrawEditFromTable(ctx, client, app, grid, row, column)
			app.SetRoot(form, true)
		}

	})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE:
			modal := tview.NewModal().
				SetText("Do you want to quit the application?").
				AddButtons([]string{"Quit", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Quit" {
						app.SetRoot(grid, true)
						app.Stop()
					}
					if buttonLabel == "Cancel" {
						app.SetRoot(grid, true)
					}
				})
			app.SetRoot(modal, true)
		case tcell.KeyCtrlA:
			formAdd := tview.NewForm().AddDropDown("Type: ", []string{"PASSWORD", "CREDIT CARD", "TEXT", "DATA"}, 0, func(option string, i int) {
				switch option {
				case "PASSWORD":
					password := models.Password{}
					formAddPass := tview.NewForm().AddInputField("Login: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { password.Login = text }).
						AddInputField("Password: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { password.Password = text }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { password.Tag = text }).
						AddButton("Add/Update", func() {
							client.AddData(ctx, password)
							//	app.SetRoot(grid, true)
							app.Stop()
							loggedIn(ctx, client)
						})
					app.SetRoot(formAddPass, true)
				case "CREDIT CARD":
					cc := models.CreditCard{}
					formAddCC := tview.NewForm().AddInputField("CC Num: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.CardNum = text }).
						AddInputField("EXP: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.Exp = text }).
						AddInputField("Name : ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.Name = text }).
						AddInputField("CVV: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.CVV = text }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.Tag = text }).
						AddButton("Add/Update", func() {
							client.AddData(ctx, cc)
							//	app.SetRoot(grid, true)
							app.Stop()
							loggedIn(ctx, client)
						})
					app.SetRoot(formAddCC, true)
				case "TEXT":
					txt := models.Text{}
					formAddText := tview.NewForm().AddInputField("Text: ", "", 100, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { txt.Data = text }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { txt.Tag = text }).
						AddButton("Add/Update", func() {
							client.AddData(ctx, txt)
							//	app.SetRoot(grid, true)
							app.Stop()
							loggedIn(ctx, client)
						})
					app.SetRoot(formAddText, true)
				case "DATA":
					data := models.Data{}
					// TODO: ADD FILE PATH PARSING!!
					formAddText := tview.NewForm().AddInputField("Path to file: ", "", 100, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { data.Data = []byte(text) }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { data.Tag = text }).
						AddButton("Add/Update", func() {
							client.AddData(ctx, data)
							//	app.SetRoot(grid, true)
							app.Stop()
							loggedIn(ctx, client)
						})
					app.SetRoot(formAddText, true)
				}
			}).AddButton("Add", nil).AddButton("Back", func() { app.SetRoot(grid, true) })
			app.SetRoot(formAdd, true)

		}

		return event
	})

	err = app.SetRoot(grid, true).EnableMouse(true).SetFocus(table).Run()
	if err != nil {
		panic(err)
	}

}
