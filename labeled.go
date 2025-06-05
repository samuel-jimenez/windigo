/*
 * Copyright (C) 2025 The Windigo Authors. All Rights Reserved.
 */

package windigo

/* Labelable
 *
 */
type Labelable interface {
	ComponentFrame
	BaseController
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
}

func NewLabeledEdit(parent Controller, label_width, control_width, height int, label_text string) *LabeledEdit {

	panel := NewAutoPanel(parent)
	panel.SetSize(label_width+control_width, height)

	edit_label := NewLabel(panel)
	edit_label.SetSize(label_width, height)
	edit_label.SetText(label_text)

	edit_field := NewEdit(panel)
	edit_field.SetText("")

	panel.Dock(edit_label, Left)
	panel.Dock(edit_field, Fill)
	return &LabeledEdit{panel, edit_field}
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
}

func NewLabeledComboBox(parent Controller, label_width, control_width, height int, label_text string) *LabeledComboBox {

	panel := NewAutoPanel(parent)
	panel.SetSize(label_width+control_width, height)

	combobox_label := NewLabel(panel)
	combobox_label.SetSize(label_width, height)
	combobox_label.SetText(label_text)

	combobox_field := NewComboBox(panel)
	combobox_field.SetText("")

	panel.Dock(combobox_label, Left)
	panel.Dock(combobox_field, Fill)
	return &LabeledComboBox{panel, combobox_field}
}

/* LabeledLabelable
 *
 */
type LabeledLabelable interface {
	Labelable
}

/* LabeledLabel
 *
 */
type LabeledLabel struct {
	ComponentFrame
	BaseController
}

func NewLabeledLabel(parent Controller, width, height int, text string) *LabeledLabel {

	panel := NewAutoPanel(parent)
	panel.SetSize(width, height)

	label := NewLabel(panel)
	// label.SetSize(width, height)
	label.SetText(text)

	// panel.Dock(label, Left)
	panel.Dock(label, Fill)
	return &LabeledLabel{panel, label}
}

/* LabeledCheckBoxable
 *
 */
type LabeledCheckBoxable interface {
	Labelable
	DiffButtonable
}

/* LabeledCheckBox
 *
 */
type LabeledCheckBox struct {
	ComponentFrame
	*CheckBox // Buttonable
}

func NewLabeledCheckBox(parent Controller, width, height int, text string) *LabeledCheckBox {

	panel := NewAutoPanel(parent)
	panel.SetSize(width, height)

	label := NewCheckBox(panel)
	// label.SetSize(width, height)
	label.SetText(text)

	// panel.Dock(label, Left)
	panel.Dock(label, Fill)
	return &LabeledCheckBox{panel, label}
}
