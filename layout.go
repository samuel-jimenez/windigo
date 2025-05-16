/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

// Marginal component must satisfy interface to have margins.
type Marginal interface {
	MarginTop() int
	MarginBtm() int
	MarginLeft() int
	MarginRight() int

	SetMarginsAll(int)
	SetMarginsVH(margin_vertical, margin_horizontal int)
	SetMargins(margin_top, margin_right, margin_btm, margin_left int)
	SetMarginTop(int)
	SetMarginBtm(int)
	SetMarginLeft(int)
	SetMarginRight(int)
}

// Dockable component must satisfy interface to be docked.
type Dockable interface {
	Marginal
	Handle() w32.HWND

	Pos() (x, y int)
	Width() int
	Height() int
	Visible() bool

	SetPos(x, y int)
	SetSize(width, height int)
	Show()
	Hide()

	OnMouseMove() *EventManager
	OnLBUp() *EventManager
}

// DockAllow is window, panel or other component that satisfies interface.
type DockAllow interface {
	Handle() w32.HWND
	ClientWidth() int
	ClientHeight() int
	SetLayout(mng LayoutManager)
}

// Various layout managers
type Direction int

const (
	Top Direction = iota
	Bottom
	Left
	Right
	Fill
)

type LayoutControl struct {
	child Dockable
	dir   Direction
}

type LayoutControls []*LayoutControl

type SimpleDock struct {
	parent      DockAllow
	layoutCtl   LayoutControls
	loadedState bool

	padding_top, padding_btm,
	padding_left, padding_right int
}

// DockState gets saved and loaded from json
type CtlState struct {
	X, Y, Width, Height int
}

type LayoutState struct {
	WindowState string
	Controls    []*CtlState
}

func (lc LayoutControls) Len() int           { return len(lc) }
func (lc LayoutControls) Swap(i, j int)      { lc[i], lc[j] = lc[j], lc[i] }
func (lc LayoutControls) Less(i, j int) bool { return lc[i].dir < lc[j].dir }

func NewSimpleDock(parent DockAllow) *SimpleDock {
	d := &SimpleDock{parent: parent}
	parent.SetLayout(d)
	return d
}

// Layout management for the child controls.
func (control *SimpleDock) Dock(child Dockable, dir Direction) {
	control.layoutCtl = append(control.layoutCtl, &LayoutControl{child, dir})
}

// Padded
func (control *SimpleDock) PaddingTop() int {
	return control.padding_top
}

func (control *SimpleDock) PaddingBtm() int {
	return control.padding_btm
}
func (control *SimpleDock) PaddingLeft() int {
	return control.padding_left
}
func (control *SimpleDock) PaddingRight() int {
	return control.padding_right
}

func (control *SimpleDock) SetPaddingsAll(padding int) {
	control.padding_top = padding
	control.padding_right = padding
	control.padding_btm = padding
	control.padding_left = padding
}

func (control *SimpleDock) SetPaddingsHV(padding_vertical, padding_horizontal int) {
	control.padding_top = padding_vertical
	control.padding_right = padding_horizontal
	control.padding_btm = padding_vertical
	control.padding_left = padding_horizontal
}

func (control *SimpleDock) SetPaddings(padding_top, padding_right, padding_btm, padding_left int) {
	control.padding_top = padding_top
	control.padding_right = padding_right
	control.padding_btm = padding_btm
	control.padding_left = padding_left
}

func (control *SimpleDock) SetPaddingTop(padding int) {
	control.padding_top = padding
}

func (control *SimpleDock) SetPaddingBtm(padding int) {
	control.padding_btm = padding
}

func (control *SimpleDock) SetPaddingLeft(padding int) {
	control.padding_left = padding
}

func (control *SimpleDock) SetPaddingRight(padding int) {
	control.padding_right = padding
}

// SaveState of the layout. Only works for Docks with parent set to main form.
func (control *SimpleDock) SaveState(w io.Writer) error {
	var ls LayoutState

	var wp w32.WINDOWPLACEMENT
	wp.Length = uint32(unsafe.Sizeof(wp))
	if !w32.GetWindowPlacement(control.parent.Handle(), &wp) {
		return fmt.Errorf("GetWindowPlacement failed")
	}

	ls.WindowState = fmt.Sprint(
		wp.Flags, wp.ShowCmd,
		wp.PtMinPosition.X, wp.PtMinPosition.Y,
		wp.PtMaxPosition.X, wp.PtMaxPosition.Y,
		wp.RcNormalPosition.Left, wp.RcNormalPosition.Top,
		wp.RcNormalPosition.Right, wp.RcNormalPosition.Bottom)

	for _, c := range control.layoutCtl {
		x, y := c.child.Pos()
		w, h := c.child.Width(), c.child.Height()

		ctl := &CtlState{X: x, Y: y, Width: w, Height: h}
		ls.Controls = append(ls.Controls, ctl)
	}

	if err := json.NewEncoder(w).Encode(ls); err != nil {
		return err
	}

	return nil
}

// LoadState of the layout. Only works for Docks with parent set to main form.
func (control *SimpleDock) LoadState(r io.Reader) error {
	var ls LayoutState

	if err := json.NewDecoder(r).Decode(&ls); err != nil {
		return err
	}

	var wp w32.WINDOWPLACEMENT
	if _, err := fmt.Sscan(ls.WindowState,
		&wp.Flags, &wp.ShowCmd,
		&wp.PtMinPosition.X, &wp.PtMinPosition.Y,
		&wp.PtMaxPosition.X, &wp.PtMaxPosition.Y,
		&wp.RcNormalPosition.Left, &wp.RcNormalPosition.Top,
		&wp.RcNormalPosition.Right, &wp.RcNormalPosition.Bottom); err != nil {
		return err
	}
	wp.Length = uint32(unsafe.Sizeof(wp))

	if !w32.SetWindowPlacement(control.parent.Handle(), &wp) {
		return fmt.Errorf("SetWindowPlacement failed")
	}

	// if number of controls in the saved layout does not match
	// current number on screen - something changed and we do not reload
	// rest of control sizes from json
	if len(control.layoutCtl) != len(ls.Controls) {
		return nil
	}

	for i, c := range control.layoutCtl {
		c.child.SetPos(ls.Controls[i].X, ls.Controls[i].Y)
		c.child.SetSize(ls.Controls[i].Width, ls.Controls[i].Height)
	}
	return nil
}

// SaveStateFile convenience function.
func (control *SimpleDock) SaveStateFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	return control.SaveState(f)
}

// LoadStateFile loads state ignores error if file is not found.
func (control *SimpleDock) LoadStateFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return nil // if file is not found or not accessible ignore it
	}
	return control.LoadState(f)
}

// Update is called to resize child items based on layout directions.
func (control *SimpleDock) Update() {

	x0, x1, y0, y1 := control.padding_left, control.parent.ClientWidth()-control.padding_right, control.padding_top, control.parent.ClientHeight()-control.padding_btm
	control_width, control_height := x1-x0,
		y1-y0

	for _, c := range control.layoutCtl {
		// Non visible controls do not preserve space.
		if !c.child.Visible() {
			continue
		}
		child_width := c.child.Width()
		child_height := c.child.Height()
		child_h_margin := c.child.MarginLeft() + c.child.MarginRight()
		child_v_margin := c.child.MarginTop() + c.child.MarginBtm()
		total_child_width := child_width + child_h_margin
		total_child_height := child_height + child_v_margin

		switch c.dir {
		case Top:
			c.child.SetPos(x0+c.child.MarginLeft(), y0+c.child.MarginTop())
			c.child.SetSize(control_width-child_h_margin, child_height)
			control_height -= total_child_height
			y0 += total_child_height
		case Bottom:
			c.child.SetPos(x0+c.child.MarginLeft(), y1-total_child_height+c.child.MarginTop())
			c.child.SetSize(control_width-child_h_margin, child_height)
			control_height -= total_child_height
			y1 -= total_child_height
		case Left:
			c.child.SetPos(x0+c.child.MarginLeft(), y0+c.child.MarginTop())
			c.child.SetSize(child_width, control_height-child_v_margin)
			control_width -= total_child_width
			x0 += total_child_width
		case Right:
			c.child.SetPos(x1-total_child_width+c.child.MarginLeft(), y0+c.child.MarginTop())
			c.child.SetSize(child_width, control_height-child_v_margin)
			control_width -= total_child_width
			x1 -= total_child_width

		case Fill:
			// fill available space
			c.child.SetPos(x0+c.child.MarginLeft(), y0+c.child.MarginTop())
			c.child.SetSize(control_width-child_h_margin, control_height-child_v_margin)
		}
		//c.child.Invalidate(true)
	}
}

type AutoPanel struct {
	Pane
	*SimpleDock
}

func NewAutoPanel(parent Controller) AutoPanel {
	panel := NewPanel(parent)
	dock := NewSimpleDock(panel)

	return AutoPanel{panel, dock}
}

func NewGroupAutoPanel(parent Controller) AutoPanel {
	panel := NewGroupPanel(parent)
	dock := NewSimpleDock(panel)

	return AutoPanel{panel, dock}
}
