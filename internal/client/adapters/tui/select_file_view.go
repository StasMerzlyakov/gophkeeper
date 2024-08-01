package tui

import (
	"os"
	"path/filepath"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	fileperm "github.com/wneessen/go-fileperm"
)

func (tApp *tuiApp) SelectFileView(current *domain.FileInfo, selectFn func(info *domain.FileInfo)) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("SelectFileView start")
			tApp.fileTreeView.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("SelectFile")
			tApp.fileTreeView.Box = box

			rootDir := "."
			root := tview.NewTreeNode(rootDir).
				SetColor(tcell.ColorRed)
			tree := tview.NewTreeView().
				SetRoot(root).
				SetCurrentNode(root) // A helper function which adds the files and directories of the given path
			// to the given target node.
			add := func(target *tview.TreeNode, path string) {

				inf, err := os.Stat(path)
				if err != nil {
					panic(err)
				}

				if inf.IsDir() {
					files, err := os.ReadDir(path)
					if err != nil {
						log.Errorf("can't create readdir %s - %s", path, err.Error())
						panic(err)
					}
					for _, file := range files {
						path := filepath.Join(path, file.Name())

						up, err := fileperm.New(path)
						if err != nil {
							log.Errorf("can't create fileperm %s", err.Error())
							panic(err)
						}
						node := tview.NewTreeNode(file.Name()).
							SetReference(path).
							SetSelectable(up.UserReadable())
						if file.IsDir() {
							node.SetColor(tcell.ColorGreen)
						}
						target.AddChild(node)
					}
				} else {
					info := &domain.FileInfo{
						Name: filepath.Base(path),
						Path: path,
					}
					selectFn(info)
				}

			}

			// Add the current directory to the root node.
			add(root, rootDir)

			// If a directory was selected, open it.
			tree.SetSelectedFunc(func(node *tview.TreeNode) {
				reference := node.GetReference()
				if reference == nil {
					return // Selecting the root node does nothing.
				}
				children := node.GetChildren()
				if len(children) == 0 {
					// Load and show files in this directory.
					path := reference.(string)
					add(node, path)
				} else {
					// Collapse if visible, expand if collapsed.
					node.SetExpanded(!node.IsExpanded())
				}
			})

			tApp.fileTreeView.
				AddItem(tree, 0, 1, true).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						selectFn(current)
					}
					return event
				})
			tApp.pages.SwitchToPage(FileTreePagh)
			log.Debug("SelectFileView shown")
		})
	}()
}

func (tApp *tuiApp) SelectDirectoryView(current *domain.FileInfo, selectFn func(info *domain.FileInfo)) {
	go func() {
		tApp.app.QueueUpdateDraw(func() {
			log := app.GetMainLogger()
			log.Debugf("SelectFileView start")
			tApp.fileTreeView.Clear()

			box := tview.NewBox().SetBorder(true).SetTitle("SelectDirectory")
			tApp.fileTreeView.Box = box

			rootDir := "."
			root := tview.NewTreeNode(rootDir).
				SetColor(tcell.ColorRed)
			tree := tview.NewTreeView().
				SetRoot(root).
				SetCurrentNode(root) // A helper function which adds the files and directories of the given path
			// to the given target node.
			add := func(target *tview.TreeNode, path string) {

				inf, err := os.Stat(path)
				if err != nil {
					panic(err)
				}

				if inf.IsDir() {
					files, err := os.ReadDir(path)
					if err != nil {
						log.Errorf("can't create readdir %s - %s", path, err.Error())
						panic(err)
					}
					for _, file := range files {
						path := filepath.Join(path, file.Name())

						node := tview.NewTreeNode(file.Name()).
							SetReference(path).
							SetSelectable(true)
						if file.IsDir() {
							node.SetColor(tcell.ColorGreen)
						}
						target.AddChild(node)
					}
				} else {
					dir := filepath.Dir(path)
					path := filepath.Join(dir, current.Name)
					info := &domain.FileInfo{
						Name: current.Name,
						Path: path,
					}
					selectFn(info)
				}

			}

			// Add the current directory to the root node.
			add(root, rootDir)

			// If a directory was selected, open it.
			tree.SetSelectedFunc(func(node *tview.TreeNode) {
				reference := node.GetReference()
				if reference == nil {
					return // Selecting the root node does nothing.
				}
				children := node.GetChildren()
				if len(children) == 0 {
					// Load and show files in this directory.
					path := reference.(string)
					add(node, path)
				} else {
					// Collapse if visible, expand if collapsed.
					node.SetExpanded(!node.IsExpanded())
				}
			})

			tApp.fileTreeView.
				AddItem(tree, 0, 1, true).
				SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyCtrlQ:
						tApp.app.Stop()
					case tcell.KeyCtrlB:
						selectFn(current)
					}
					return event
				})
			tApp.pages.SwitchToPage(FileTreePagh)
			log.Debug("SelectFileView shown")
		})
	}()
}
