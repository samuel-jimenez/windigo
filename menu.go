/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

var (
	nextMenuItemID  uint16 = 3
	actionsByID            = make(map[uint16]*MenuItem)
	shortcut2Action        = make(map[Shortcut]*MenuItem)
	menuItems              = make(map[w32.HMENU][]*MenuItem)
	radioGroups            = make(map[*MenuItem]*RadioGroup)
	initialised     bool
)

var NoShortcut = Shortcut{}

// Menu for main window and context menus on controls.
// Most methods used for both main window menu and context menu.
type Menu struct {
	hMenu w32.HMENU
	hwnd  w32.HWND // root window. hwnd might be nil if it is context menu.
}

type MenuItem struct {
	hMenu    w32.HMENU
	hSubMenu w32.HMENU // Non zero if this item is in itself a submenu.

	text     string
	toolTip  string
	image    *Bitmap
	shortcut Shortcut
	enabled  bool

	checkable bool
	checked   bool
	isRadio   bool

	id uint16

	onClick EventManager
}

type RadioGroup struct {
	members []*MenuItem
	hwnd    w32.HWND
}

func NewContextMenu() *MenuItem {
	hMenu := w32.CreatePopupMenu()
	if hMenu == 0 {
		panic("failed CreateMenu")
	}

	item := &MenuItem{
		hMenu:    hMenu,
		hSubMenu: hMenu,
	}
	return item
}

func (control *Menu) Dispose() {
	if control.hMenu != 0 {
		w32.DestroyMenu(control.hMenu)
		control.hMenu = 0
	}
}

func (control *Menu) IsDisposed() bool {
	return control.hMenu == 0
}

func initMenuItemInfoFromAction(item_info *w32.MENUITEMINFO, menu_item *MenuItem) {
	item_info.CbSize = uint32(unsafe.Sizeof(*item_info))
	item_info.FMask = w32.MIIM_FTYPE | w32.MIIM_ID | w32.MIIM_STATE | w32.MIIM_STRING
	if menu_item.image != nil {
		item_info.FMask |= w32.MIIM_BITMAP
		item_info.HbmpItem = menu_item.image.handle
	}
	if menu_item.IsSeparator() {
		item_info.FType = w32.MFT_SEPARATOR
	} else {
		item_info.FType = w32.MFT_STRING
		var text string
		if s := menu_item.shortcut; s.Key != 0 {
			text = fmt.Sprintf("%s\t%s", menu_item.text, s.String())
			shortcut2Action[menu_item.shortcut] = menu_item
		} else {
			text = menu_item.text
		}
		item_info.DwTypeData = syscall.StringToUTF16Ptr(text)
		item_info.Cch = uint32(len([]rune(menu_item.text)))
	}
	item_info.WID = uint32(menu_item.id)

	if menu_item.Enabled() {
		item_info.FState &^= w32.MFS_DISABLED
	} else {
		item_info.FState |= w32.MFS_DISABLED
	}

	if menu_item.Checkable() {
		item_info.FMask |= w32.MIIM_CHECKMARKS
	}
	if menu_item.Checked() {
		item_info.FState |= w32.MFS_CHECKED
	}

	if menu_item.hSubMenu != 0 {
		item_info.FMask |= w32.MIIM_SUBMENU
		item_info.HSubMenu = menu_item.hSubMenu
	}
}

// Show menu on the main window.
func (control *Menu) Show() {
	initialised = true
	updateRadioGroups()
	if !w32.DrawMenuBar(control.hwnd) {
		panic("DrawMenuBar failed")
	}
}

// AddSubMenu returns item that is used as submenu to perform AddItem(s).
func (control *Menu) AddSubMenu(text string) *MenuItem {
	hSubMenu := w32.CreateMenu()
	if hSubMenu == 0 {
		panic("failed CreateMenu")
	}
	return addMenuItem(control.hMenu, hSubMenu, text, Shortcut{}, nil, false)
}

// This method will iterate through the menu items, group radio items together, build a
// quick access map and set the initial items
func updateRadioGroups() {

	if !initialised {
		return
	}

	radioItemsChecked := []*MenuItem{}
	radioGroups = make(map[*MenuItem]*RadioGroup)
	var currentRadioGroupMembers []*MenuItem
	// Iterate the menus
	for _, menu := range menuItems {
		menuLength := len(menu)
		for index, menuItem := range menu {
			if menuItem.isRadio {
				currentRadioGroupMembers = append(currentRadioGroupMembers, menuItem)
				if menuItem.checked {
					radioItemsChecked = append(radioItemsChecked, menuItem)
				}

				// If end of menu
				if index == menuLength-1 {
					radioGroup := &RadioGroup{
						members: currentRadioGroupMembers,
						hwnd:    menuItem.hMenu,
					}
					// Save the group to each member iin the radiomap
					for _, member := range currentRadioGroupMembers {
						radioGroups[member] = radioGroup
					}
					currentRadioGroupMembers = []*MenuItem{}
				}
				continue
			}

			// Not a radio item
			if len(currentRadioGroupMembers) > 0 {
				radioGroup := &RadioGroup{
					members: currentRadioGroupMembers,
					hwnd:    menuItem.hMenu,
				}
				// Save the group to each member iin the radiomap
				for _, member := range currentRadioGroupMembers {
					radioGroups[member] = radioGroup
				}
				currentRadioGroupMembers = []*MenuItem{}
			}
		}
	}

	// Enable the checked items
	for _, item := range radioItemsChecked {
		radioGroup := radioGroups[item]
		startID := radioGroup.members[0].id
		endID := radioGroup.members[len(radioGroup.members)-1].id
		w32.SelectRadioMenuItem(item.id, startID, endID, radioGroup.hwnd)
	}

}

func (control *MenuItem) OnClick() *EventManager {
	return &control.onClick
}

func (control *MenuItem) AddSeparator() {
	addMenuItem(control.hSubMenu, 0, "-", Shortcut{}, nil, false)
}

// AddItem adds plain menu item.
func (control *MenuItem) AddItem(text string, shortcut Shortcut) *MenuItem {
	return addMenuItem(control.hSubMenu, 0, text, shortcut, nil, false)
}

// AddItemCheckable adds plain menu item that can have a checkmark.
func (control *MenuItem) AddItemCheckable(text string, shortcut Shortcut) *MenuItem {
	return addMenuItem(control.hSubMenu, 0, text, shortcut, nil, true)
}

// AddItemRadio adds plain menu item that can have a checkmark and is part of a radio group.
func (control *MenuItem) AddItemRadio(text string, shortcut Shortcut) *MenuItem {
	menuItem := addMenuItem(control.hSubMenu, 0, text, shortcut, nil, true)
	menuItem.isRadio = true
	return menuItem
}

// AddItemWithBitmap adds menu item with shortcut and bitmap.
func (control *MenuItem) AddItemWithBitmap(text string, shortcut Shortcut, image *Bitmap) *MenuItem {
	return addMenuItem(control.hSubMenu, 0, text, shortcut, image, false)
}

// AddSubMenu adds a submenu.
func (control *MenuItem) AddSubMenu(text string) *MenuItem {
	hSubMenu := w32.CreatePopupMenu()
	if hSubMenu == 0 {
		panic("failed CreatePopupMenu")
	}
	return addMenuItem(control.hSubMenu, hSubMenu, text, Shortcut{}, nil, false)
}

// AddItem to the menu, set text to "-" for separators.
func addMenuItem(hMenu, hSubMenu w32.HMENU, text string, shortcut Shortcut, image *Bitmap, checkable bool) *MenuItem {
	item := &MenuItem{
		hMenu:     hMenu,
		hSubMenu:  hSubMenu,
		text:      text,
		shortcut:  shortcut,
		image:     image,
		enabled:   true,
		id:        nextMenuItemID,
		checkable: checkable,
		isRadio:   false,
		//visible:  true,
	}
	nextMenuItemID++
	actionsByID[item.id] = item
	menuItems[hMenu] = append(menuItems[hMenu], item)

	var item_info w32.MENUITEMINFO
	initMenuItemInfoFromAction(&item_info, item)

	index := -1
	if !w32.InsertMenuItem(hMenu, uint32(index), true, &item_info) {
		panic("InsertMenuItem failed")
	}
	return item
}

func indexInObserver(menu_item *MenuItem) int {
	var idx int
	for _, control := range menuItems[menu_item.hMenu] {
		if control == menu_item {
			return idx
		}
		idx++
	}
	return -1
}

func findMenuItemByID(id int) *MenuItem {
	return actionsByID[uint16(id)]
}

func (control *MenuItem) update() {
	var item_info w32.MENUITEMINFO
	initMenuItemInfoFromAction(&item_info, control)

	if !w32.SetMenuItemInfo(control.hMenu, uint32(indexInObserver(control)), true, &item_info) {
		panic("SetMenuItemInfo failed")
	}
	if control.isRadio {
		control.updateRadioGroup()
	}
}

func (control *MenuItem) IsSeparator() bool { return control.text == "-" }
func (control *MenuItem) SetSeparator()     { control.text = "-" }

func (control *MenuItem) Enabled() bool     { return control.enabled }
func (control *MenuItem) SetEnabled(b bool) { control.enabled = b; control.update() }

func (control *MenuItem) Checkable() bool     { return control.checkable }
func (control *MenuItem) SetCheckable(b bool) { control.checkable = b; control.update() }

func (control *MenuItem) Checked() bool { return control.checked }
func (control *MenuItem) SetChecked(b bool) {
	if control.isRadio {
		radioGroup := radioGroups[control]
		if radioGroup != nil {
			for _, member := range radioGroup.members {
				member.checked = false
			}
		}

	}
	control.checked = b
	control.update()
}

func (control *MenuItem) Text() string     { return control.text }
func (control *MenuItem) SetText(s string) { control.text = s; control.update() }

func (control *MenuItem) Image() *Bitmap     { return control.image }
func (control *MenuItem) SetImage(b *Bitmap) { control.image = b; control.update() }

func (control *MenuItem) ToolTip() string     { return control.toolTip }
func (control *MenuItem) SetToolTip(s string) { control.toolTip = s; control.update() }

func (control *MenuItem) updateRadioGroup() {
	radioGroup := radioGroups[control]
	if radioGroup == nil {
		return
	}
	startID := radioGroup.members[0].id
	endID := radioGroup.members[len(radioGroup.members)-1].id
	w32.SelectRadioMenuItem(control.id, startID, endID, radioGroup.hwnd)
}
