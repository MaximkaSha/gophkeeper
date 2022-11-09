package ui

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/MaximkaSha/gophkeeper/internal/client"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/theplant/luhn"
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
	status.SetText("Client version: " + client.BuildVersion + ", client build time: " + client.BuildTime)

	grid2 := tview.NewFlex().AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(logo, 0, 2, false).
			AddItem(form, 0, 1, true).
			AddItem(status, 0, 1, false), 0, 5, true).
		AddItem(tview.NewBox(), 0, 1, false)

	/*grid := tview.NewGrid().
	AddItem(logo, 1, 1, 1, 3, 0, 0, false).
	AddItem(form, 2, 1, 1, 3, 0, 0, true).
	AddItem(status, 3, 1, 1, 3, 0, 0, false) */
	err := app.SetRoot(grid2, true).Run()
	if err != nil {
		panic(err)
	}
}

func DrawError(err error) {
	app := tview.NewApplication()

	modal := tview.NewModal().
		SetText("Wrong CC Num or EXP").AddButtons([]string{"", "OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
			}
		}).SetTextColor(tcell.ColorRed)
	app.SetRoot(modal, true)
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

func CheckCCNum(textToCheck string) bool {
	ccNum, _ := strconv.Atoi(textToCheck)
	return luhn.Valid(ccNum)
}
func CheckCCExp(textToCheck string) bool {
	exp := strings.Split(textToCheck, "/")
	if len(exp) < 2 {
		exp = strings.Split(textToCheck, "\\")
		if len(exp) < 2 {
			return false
		}

	}
	month, _ := strconv.Atoi(exp[0])
	year, _ := strconv.Atoi(exp[1])
	return !(month < 12) || !(year < 70)
}

func DrawEditFromTable(ctx context.Context, client client.Client, app *tview.Application, grid *tview.Grid, row int, dType int) *tview.Form {
	form := tview.NewForm()
	switch dType {
	case 1:
		isChanged := false
		cc := client.LocalStorage.CCStorage[row-1]
		form.AddInputField("CCNum: ", cc.CardNum, 17, nil, func(text string) {
			isChanged = true
			cc.CardNum = text
		}).
			AddInputField("Exp: ", cc.Exp, 8, nil, func(text string) {
				isChanged = true
				cc.Exp = text
			}).
			AddInputField("Name: ", cc.Name, 16, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) {
				isChanged = true
				cc.Name = text
			}).
			AddInputField("CVV: ", cc.CVV, 5, func(textToCheck string, lastChar rune) bool { return len(textToCheck) == 4 }, func(text string) {
				isChanged = true
				cc.CVV = text
			}).
			AddInputField("Tag: ", cc.Tag, 16, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) {
				isChanged = true
				cc.Tag = text
			}).
			AddButton("Add/Update", func() {
				if isChanged {
					if !CheckCCNum(cc.CardNum) || !CheckCCExp(cc.Exp) {
						DrawError(errors.New("wrong cc num or exp"))
						app.Stop()
						loggedIn(ctx, client)
					}
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
				pass.Tag = text
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

		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		data := client.LocalStorage.DataStorage[row-1]
		form.AddInputField("Path to write: ", path+data.Tag, len(path+data.Tag), nil, func(text string) {
			path = text
		}).
			AddInputField("Tag: ", data.Tag, 16, nil, func(text string) {

				data.Tag = text
			}).
			AddButton("Write file", func() {
				os.WriteFile(path, data.Data, fs.FileMode(os.O_WRONLY))
				app.Stop()
				loggedIn(ctx, client)
			}).AddButton("Del", func() {
			client.DelData(ctx, data)
			app.Stop()
			loggedIn(ctx, client)
		}).AddButton("Exit", func() { app.Stop() })

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
	table.SetCell(0, 3, tview.NewTableCell(fmt.Sprintf("Files(%v)", len(client.LocalStorage.DataStorage))).SetExpansion(1).SetAlign(tview.AlignCenter).SetBackgroundColor(tcell.Color100))
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
			formAdd := tview.NewForm().AddDropDown("Type: ", []string{"PASSWORD", "CREDIT CARD", "TEXT", "FILE"}, 0, func(option string, i int) {
				switch option {
				case "PASSWORD":
					password := models.Password{}
					formAddPass := tview.NewForm().AddInputField("Login: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { password.Login = text }).
						AddInputField("Password: ", "", 16, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { password.Password = text }).
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
					formAddCC := tview.NewForm().AddInputField("CC Num: ", "", 17, nil, func(text string) { cc.CardNum = text }).
						AddInputField("EXP: ", "", 7, nil, func(text string) { cc.Exp = text }).
						AddInputField("Name : ", "", 16, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.Name = text }).
						AddInputField("CVV: ", "", 5, func(textToCheck string, lastChar rune) bool { return len(textToCheck) < 4 }, func(text string) { cc.CVV = text }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { cc.Tag = text }).
						AddButton("Add/Update", func() {
							if !CheckCCNum(cc.CardNum) || !CheckCCExp(cc.Exp) {
								app.Stop()
							}
							client.AddData(ctx, cc)
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
					path := Tree()
					formAddText := tview.NewForm().AddInputField("Path to file: ", path, len(path), func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { data.Data = []byte(text) }).
						AddInputField("Tag: ", "", 10, func(textToCheck string, lastChar rune) bool { return textToCheck != "" }, func(text string) { data.Tag = text }).
						AddButton("Add/Update", func() {
							b, err := ioutil.ReadFile(path)
							if err != nil {
								log.Fatal(err)
							}
							data.Data = b
							client.AddData(ctx, data)
							app.Stop()
							loggedIn(ctx, client)
						}).SetFocus(1)
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

func Tree() string {
	app := tview.NewApplication()
	rootDir := "."
	if runtime.GOOS == "windows" {
		rootDir = "C:\\"
	}
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name()+"f")).
				SetSelectable(true)
			if file.IsDir() {
				node = tview.NewTreeNode(file.Name()).
					SetReference(filepath.Join(path, file.Name()+"d")).
					SetSelectable(true)
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}

	}

	// Add the current directory to the root node.
	add(root, rootDir)
	path := ""
	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {

		reference := node.GetReference()
		lastChar := reference.(string)[len(reference.(string))-1:]
		if lastChar == "f" {
			path = reference.(string)[:len(reference.(string))-1]
			app.Stop()
			return

		}
		if (reference == nil) || (reference == "d") || (reference == "f") {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path[:len(path)-1])
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := app.SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
	return path
}
