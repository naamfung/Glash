package dialog

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"glash/internal/ui/common"
	uv "github.com/charmbracelet/ultraviolet"
)

// ModelSourceSelectID is the identifier for the model source select dialog.
const ModelSourceSelectID = "modelSourceSelect"

// ModelSourceType represents the type of model source to select.
type ModelSourceType int

const (
	ModelSourceOnline ModelSourceType = iota
	ModelSourceLocal
)

// String returns the string representation of the [ModelSourceType].
func (ms ModelSourceType) String() string {
	switch ms {
	case ModelSourceOnline:
		return "在线模型"
	case ModelSourceLocal:
		return "本地模型"
	default:
		return "Unknown"
	}
}

// ModelSourceSelect represents a model source selection dialog.
type ModelSourceSelect struct {
	com           *common.Common
	selected      ModelSourceType
	back          bool
	keyMap        struct {
		LeftRight  key.Binding
		Enter      key.Binding
		Close      key.Binding
		BackToMenu key.Binding
	}
}

var _ Dialog = (*ModelSourceSelect)(nil)

// NewModelSourceSelect creates a new ModelSourceSelect dialog.
func NewModelSourceSelect(com *common.Common) *ModelSourceSelect {
	m := &ModelSourceSelect{
		com:       com,
		selected:  ModelSourceLocal,
		back:      false,
	}

	m.keyMap.LeftRight = key.NewBinding(
		key.WithKeys("left", "right", "tab", "shift+tab"),
		key.WithHelp("←/→", "switch source"),
	)
	m.keyMap.Enter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	)
	m.keyMap.Close = CloseKey
	m.keyMap.BackToMenu = key.NewBinding(
		key.WithKeys("h", "b"),
		key.WithHelp("h/b", "back"),
	)

	return m
}

// ID implements Dialog.
func (m *ModelSourceSelect) ID() string {
	return ModelSourceSelectID
}

// HandleMsg implements Dialog.
func (m *ModelSourceSelect) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Close, m.keyMap.BackToMenu):
			return ActionClose{}
		case key.Matches(msg, m.keyMap.LeftRight):
			if m.selected == ModelSourceOnline {
				m.selected = ModelSourceLocal
			} else {
				m.selected = ModelSourceOnline
			}
		case key.Matches(msg, m.keyMap.Enter):
			return ActionSelectModelSource{Source: m.selected}
		}
	}
	return nil
}

// Draw implements Dialog.
func (m *ModelSourceSelect) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := m.com.Styles

	width := max(0, min(defaultModelsDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))

	rc := NewRenderContext(t, width)
	rc.Title = "选择模型来源"

	onlineStyle := t.Radio.Off.Padding(0, 2)
	localStyle := t.Radio.Off.Padding(0, 2)

	if m.selected == ModelSourceOnline {
		onlineStyle = t.Radio.On.Padding(0, 2)
	} else {
		localStyle = t.Radio.On.Padding(0, 2)
	}

	options := lipgloss.JoinHorizontal(lipgloss.Center,
		onlineStyle.Render("在线模型"),
		lipgloss.NewStyle().Padding(0, 2).Render(" | "),
		localStyle.Render("本地模型"),
	)

	contentPanelWidth := width - t.Dialog.View.GetHorizontalFrameSize()
	contentPanelStyle := t.Dialog.ContentPanel.Width(contentPanelWidth).Align(lipgloss.Center)
	rc.AddPart(contentPanelStyle.Render(options))

	view := rc.Render()
	DrawCenter(scr, area, view)
	return nil
}

// ShortHelp implements help.KeyMap.
func (m *ModelSourceSelect) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.LeftRight,
		m.keyMap.Enter,
	}
}

// FullHelp implements help.KeyMap.
func (m *ModelSourceSelect) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keyMap.LeftRight, m.keyMap.Enter, m.keyMap.Close},
	}
}

// ActionSelectModelSource is a message indicating a model source has been selected.
type ActionSelectModelSource struct {
	Source ModelSourceType
}
