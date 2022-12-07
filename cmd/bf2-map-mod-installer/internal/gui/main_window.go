package gui

import (
	"fmt"

	"github.com/cetteup/bf2-map-mod-installer/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/rs/zerolog/log"
)

const (
	windowWidth  = 300
	windowHeight = 125
)

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

func CreateMainWindow(finder *software_finder.SoftwareFinder) (*walk.MainWindow, error) {
	icon, err := walk.NewIconFromResourceIdWithSize(2, walk.Size{Width: 256, Height: 256})
	if err != nil {
		return nil, err
	}

	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)

	var mw *walk.MainWindow
	var itemLabel *walk.Label
	var installBtn *walk.PushButton
	var uninstallBtn *walk.PushButton
	var config *internal.Config
	var bf2InstallPath string

	if err := (declarative.MainWindow{
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
				Text:       "BF2 map mod installer v0.1.0",
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

	if len(config.InstallItems) > 0 {
		// Update item label
		err = itemLabel.SetText(fmt.Sprintf("Configuration contains %d mod(s), %d map(s)", len(config.GetItemOfType(internal.ItemTypeMod)), len(config.GetItemOfType(internal.ItemTypeMap))))
		if err != nil {
			log.Error().Err(err).Msg("Failed to update item label text")
			return mw, err
		}

		// Enable buttons
		installBtn.SetEnabled(true)
		uninstallBtn.SetEnabled(true)
	}

	// Determine BF2 install path
	bf2InstallPath, err = finder.GetInstallDir(software_finder.Config{
		ForType:           software_finder.RegistryFinder,
		RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
		RegistryValueName: "InstallDir",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to determine BF2 install path")
		walk.MsgBox(mw, "Error", fmt.Sprintf("Failed to determine path to Battlefield 2 installation\n%s", err), walk.MsgBoxIconError)
		return mw, err
	}

	return mw, nil
}
