package gui

import (
	"fmt"

	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/bf2-map-mod-installer/internal"
)

const (
	windowWidth  = 300
	windowHeight = 249
)

type finder interface {
	GetInstallDirFromSomewhere(configs []software_finder.Config) (string, error)
}

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

func CreateMainWindow(f finder) (*walk.MainWindow, error) {
	icon, err := walk.NewIconFromResourceIdWithSize(2, walk.Size{Width: 256, Height: 256})
	if err != nil {
		return nil, err
	}

	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)

	var mw *walk.MainWindow
	var itemLabel *walk.Label
	var pathTE *walk.TextEdit
	var installBtn *walk.PushButton
	var uninstallBtn *walk.PushButton
	var config *internal.Config
	var bf2InstallPath string

	enableInstallActions := func(path string) {
		_ = pathTE.SetText(path)
		_ = pathTE.SetToolTipText(path)
		installBtn.SetEnabled(true)
		uninstallBtn.SetEnabled(true)
	}

	if err = (declarative.MainWindow{
		AssignTo: &mw,
		Title:    "BF2 map mod installer",
		Name:     "BF2 map mod installer",
		Bounds: declarative.Rectangle{
			X:      int((screenWidth - windowWidth) / 2),
			Y:      int((screenHeight - windowHeight) / 2),
			Width:  windowWidth,
			Height: windowHeight,
		},
		Layout:  declarative.VBox{},
		Icon:    icon,
		ToolBar: declarative.ToolBar{},
		Children: []declarative.Widget{
			declarative.Label{
				AssignTo:   &itemLabel,
				Text:       "Config contains no items",
				Alignment:  declarative.AlignHCenterVCenter,
				TextColor:  walk.Color(win.GetSysColor(win.COLOR_CAPTIONTEXT)),
				Background: declarative.SolidColorBrush{Color: walk.Color(win.GetSysColor(win.COLOR_BTNFACE))},
			},
			declarative.VSpacer{Size: 1},
			declarative.GroupBox{
				Title:  "Installation folder",
				Name:   "Installation folder",
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.TextEdit{
						AssignTo: &pathTE,
						Name:     "Installation folder",
						ReadOnly: true,
					},
					declarative.HSplitter{
						Children: []declarative.Widget{
							declarative.PushButton{
								Text: "Detect",
								OnClicked: func() {
									detected, err2 := detectInstallPath(f)
									if err2 != nil {
										walk.MsgBox(mw, "Warning", "Could not detect game installation folder, please choose the path manually", walk.MsgBoxIconWarning)
										return
									}

									enableInstallActions(detected)
								},
							},
							declarative.PushButton{
								Text: "Choose",
								OnClicked: func() {
									dlg := &walk.FileDialog{
										Title: "Choose installation folder",
									}

									ok, err2 := dlg.ShowBrowseFolder(mw)
									if err2 != nil {
										walk.MsgBox(mw, "Error", fmt.Sprintf("Failed to choose installation folder: %s", err2.Error()), walk.MsgBoxIconError)
										return
									} else if !ok {
										// User canceled dialog
										return
									}

									enableInstallActions(dlg.FilePath)
								},
							},
						},
					},
				},
			},
			declarative.VSpacer{Size: 1},
			declarative.PushButton{
				AssignTo: &installBtn,
				Text:     "Install",
				Enabled:  false,
				OnClicked: func() {
					mw.SetEnabled(false)
					err = internal.InstallItems(bf2InstallPath, config.InstallItems)
					if err != nil {
						walk.MsgBox(mw, "Error", fmt.Sprintf("Installation failed\n%s", err), walk.MsgBoxIconError)
					} else {
						walk.MsgBox(mw, "Success", "Installation succeeded", walk.MsgBoxIconInformation)
					}
					mw.SetEnabled(true)
				},
			},
			declarative.PushButton{
				AssignTo: &uninstallBtn,
				Text:     "Uninstall",
				Enabled:  false,
				OnClicked: func() {
					mw.SetEnabled(false)
					err = internal.UninstallItems(bf2InstallPath, config.InstallItems)
					if err != nil {
						walk.MsgBox(mw, "Error", fmt.Sprintf("Uninstallation failed\n%s", err), walk.MsgBoxIconError)
					} else {
						walk.MsgBox(mw, "Success", "Uninstallation succeeded", walk.MsgBoxIconInformation)
					}
					mw.SetEnabled(true)
				},
			},
			declarative.Label{
				Text:       "BF2 map mod installer v0.1.1",
				Alignment:  declarative.AlignHCenterVCenter,
				TextColor:  walk.Color(win.GetSysColor(win.COLOR_GRAYTEXT)),
				Background: declarative.SolidColorBrush{Color: walk.Color(win.GetSysColor(win.COLOR_BTNFACE))},
			},
		},
	}).Create(); err != nil {
		return nil, err
	}

	// Disable minimize/maximize buttons and fix size
	win.SetWindowLong(mw.Handle(), win.GWL_STYLE, win.GetWindowLong(mw.Handle(), win.GWL_STYLE) & ^win.WS_MINIMIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX)

	// Load install item config
	config, err = internal.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load config")
		walk.MsgBox(mw, "Error", fmt.Sprintf("Failed to load config\n%s", err), walk.MsgBoxIconError)
		return mw, err
	}

	// Update item label
	if len(config.InstallItems) > 0 {
		_ = itemLabel.SetText(fmt.Sprintf(
			"Configuration contains %d mod(s), %d map(s)",
			len(config.GetItemOfType(internal.ItemTypeMod)),
			len(config.GetItemOfType(internal.ItemTypeMap)),
		))
	}

	// Automatically try to detect install path once, pre-filling path if path is detected
	detected, err := detectInstallPath(f)
	if err == nil {
		enableInstallActions(detected)
	}

	return mw, nil
}

func detectInstallPath(f finder) (string, error) {
	// Copied from https://github.com/cetteup/joinme.click-launcher/blob/089fb595adc426aab775fe40165431501a5c38c3/internal/titles/bf2.go#L37
	dir, err := f.GetInstallDirFromSomewhere([]software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
			RegistryValueName: "InstallDir",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyCurrentUser,
			RegistryPath:      "SOFTWARE\\BF2Hub Systems\\BF2Hub Client",
			RegistryValueName: "bf2Dir",
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to determine Battlefield 2 install directory: %w", err)
	}

	return dir, err
}
