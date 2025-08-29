/*
 * Copyright (C) 2025 The Windigo Authors. All Rights Reserved.
 */

package windigo

/* BaseLabelable
 *
 */
type BaseLabelable interface {
	ComponentFrame
	BaseController
}

/* DiffLabelable
 *
 */
type DiffLabelable interface {
	Label() *Label
	SetLabeledSize(label_width, control_width, height int)
}

/* Labeled
 *
 */
type Labeled struct {
	FieldLabel *Label
}

func (control Labeled) Label() *Label { return control.FieldLabel }

/* Labelable
 *
 */
type Labelable interface {
	BaseLabelable
	DiffLabelable
}

/* LabeledEditable
 *
 */
type LabeledEditable interface {
	Labelable
	DiffEditable
}

/* LabeledEdit
 *
 */
type LabeledEdit struct {
	ComponentFrame
	*Edit //Editable
	*Labeled
}

func NewLabeledEdit(parent Controller, label_text string) *LabeledEdit {

	panel := NewAutoPanel(parent)

	label := NewLabel(panel)
	label.SetText(label_text)

	field := NewEdit(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledEdit{panel, field, &Labeled{label}}
}

func NewSizedLabeledEdit(parent Controller, label_width, control_width, height int, label_text string) *LabeledEdit {

	panel := NewAutoPanel(parent)
	panel.SetSize(label_width+control_width, height)

	label := NewLabel(panel)
	label.SetSize(label_width, height)
	label.SetText(label_text)

	field := NewEdit(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledEdit{panel, field, &Labeled{label}}
}

func (control *LabeledEdit) SetFont(font *Font) {
	control.Edit.SetFont(font)
	control.Label().SetFont(font)
}
func (control *LabeledEdit) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	control.Label().SetSize(label_width, height)
}

/* LabeledComboBoxable
 *
 */
type LabeledComboBoxable interface {
	Labelable
	DiffComboBoxable
}

/* LabeledComboBox
 *
 */
type LabeledComboBox struct {
	ComponentFrame
	*ComboBox //ComboBoxable
	*Labeled
}

func NewLabeledComboBox(parent Controller, label_text string) *LabeledComboBox {

	panel := NewAutoPanel(parent)

	label := NewLabel(panel)
	label.SetText(label_text)

	field := NewComboBox(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledComboBox{panel, field, &Labeled{label}}

}

func NewSizedLabeledComboBox(parent Controller, label_width, control_width, height int, label_text string) *LabeledComboBox {

	panel := NewAutoPanel(parent)
	panel.SetSize(label_width+control_width, height)

	label := NewLabel(panel)
	label.SetSize(label_width, height)
	label.SetText(label_text)

	field := NewComboBox(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledComboBox{panel, field, &Labeled{label}}

}

func NewLabeledListComboBox(parent Controller, label_text string) *LabeledComboBox {

	panel := NewAutoPanel(parent)

	label := NewLabel(panel)
	label.SetText(label_text)

	field := NewListComboBox(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledComboBox{panel, field, &Labeled{label}}

}

func NewSizedLabeledListComboBox(parent Controller, label_width, control_width, height int, label_text string) *LabeledComboBox {

	panel := NewAutoPanel(parent)
	panel.SetSize(label_width+control_width, height)

	label := NewLabel(panel)
	label.SetSize(label_width, height)
	label.SetText(label_text)

	field := NewListComboBox(panel)
	field.SetText("")

	panel.Dock(label, Left)
	panel.Dock(field, Fill)
	return &LabeledComboBox{panel, field, &Labeled{label}}

}

func (control *LabeledComboBox) SetFont(font *Font) {
	control.ComboBox.SetFont(font)
	control.Label().SetFont(font)
}

func (control *LabeledComboBox) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	control.Label().SetSize(label_width, height)
}

func (control *LabeledComboBox) SetFGColor(color Color) {
	control.ComponentFrame.(*AutoPanel).SetFGColor(color)
	control.ComboBox.SetFGColor(color)
}

func (control *LabeledComboBox) ClearFGColor() {
	control.ComponentFrame.(*AutoPanel).ClearFGColor()
	control.ComboBox.ClearFGColor()
}

func (control *LabeledComboBox) SetBGColor(color Color) {
	control.ComponentFrame.(*AutoPanel).SetBGColor(color)
	control.ComboBox.SetBGColor(color)
}

func (control *LabeledComboBox) ClearBGColor() {
	control.ComponentFrame.(*AutoPanel).ClearBGColor()
	control.ComboBox.ClearBGColor()
}

func (control *LabeledComboBox) SetHighlightColor(color Color) {
	control.ComponentFrame.(*AutoPanel).SetHighlightColor(color)
	control.ComboBox.SetHighlightColor(color)
}

func (control *LabeledComboBox) ClearHighlightColor() {
	control.ComponentFrame.(*AutoPanel).ClearHighlightColor()
	control.ComboBox.ClearHighlightColor()
}

/* LabeledLabelable
 *
 */
type LabeledLabelable interface {
	BaseLabelable
}

/* LabeledLabel
 *
 */
type LabeledLabel struct {
	ComponentFrame
	*Label // BaseController
}

func NewLabeledLabel(parent Controller, text string) *LabeledLabel {

	panel := NewAutoPanel(parent)

	label := NewLabel(panel)
	label.SetText(text)

	panel.Dock(label, Fill)
	return &LabeledLabel{panel, label}
}

func NewSizedLabeledLabel(parent Controller, width, height int, text string) *LabeledLabel {

	label := NewLabeledLabel(parent, text)
	label.SetSize(width, height)
	return label
}

func (control *LabeledLabel) SetFont(font *Font) {
	control.Label.SetFont(font)
}

/* LabeledCheckBoxable
 *
 */
type LabeledCheckBoxable interface {
	BaseLabelable
	DiffButtonable
}

/* LabeledCheckBox
 *
 */
type LabeledCheckBox struct {
	ComponentFrame
	*CheckBox // Buttonable
}

func NewLabeledCheckBox(parent Controller, text string) *LabeledCheckBox {

	panel := NewAutoPanel(parent)

	label := NewCheckBox(panel)
	label.SetText(text)

	panel.Dock(label, Fill)
	return &LabeledCheckBox{panel, label}
}

func (control *LabeledCheckBox) SetFont(font *Font) {
	control.CheckBox.SetFont(font)
}
