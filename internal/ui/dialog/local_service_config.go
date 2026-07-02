package dialog

import (
	"fmt"
	"net/url"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/catwalk/pkg/catwalk"
	"charm.land/lipgloss/v2"
	"glash/internal/config"
	"glash/internal/ui/common"
	"glash/internal/ui/util"
	uv "github.com/charmbracelet/ultraviolet"
)

// LocalServiceConfigID is the identifier for the local service config dialog.
const LocalServiceConfigID = "localServiceConfig"

// LocalServiceType represents the type of local service.
type LocalServiceType string

const (
	LocalServiceOpenAI LocalServiceType = "openai"
	LocalServiceOllama LocalServiceType = "ollama"
	LocalServiceLiteLLM LocalServiceType = "litellm"
	LocalServiceLMStudio LocalServiceType = "lmstudio"
	LocalServiceLlamaCPP LocalServiceType = "llamacpp"
	LocalServiceOmlx     LocalServiceType = "omlx"
	LocalServiceAnthropic LocalServiceType = "anthropic"
)

// String returns the string representation of the [LocalServiceType].
func (lst LocalServiceType) String() string {
	switch lst {
	case LocalServiceOpenAI:
		return "OpenAI Compatible"
	case LocalServiceOllama:
		return "Ollama"
	case LocalServiceLiteLLM:
		return "LiteLLM"
	case LocalServiceLMStudio:
		return "LM Studio"
	case LocalServiceLlamaCPP:
		return "llama.cpp"
	case LocalServiceAnthropic:
		return "Anthropic Compatible"
	default:
		return "Unknown"
	}
}

// LocalServiceConfig represents a local service configuration dialog.
type LocalServiceConfig struct {
	com *common.Common

	serviceTypes           []LocalServiceType
	selectedServiceTypeIndex int
	interfaceType          LocalServiceType

	baseURLInput textinput.Model
	modelInput   textinput.Model

	keyMap struct {
		Next         key.Binding
		Previous     key.Binding
		Left         key.Binding
		Right        key.Binding
		Submit       key.Binding
		Close        key.Binding
		BackToOnline key.Binding
	}
	help         help.Model
	focusedInput int // 0: baseURL, 1: model
}

var _ Dialog = (*LocalServiceConfig)(nil)

// NewLocalServiceConfig creates a new LocalServiceConfig dialog.
func NewLocalServiceConfig(com *common.Common) *LocalServiceConfig {
	m := &LocalServiceConfig{
		com:                    com,
		serviceTypes:           []LocalServiceType{LocalServiceLlamaCPP, LocalServiceOllama, LocalServiceOpenAI, LocalServiceLiteLLM, LocalServiceLMStudio, LocalServiceAnthropic},
		selectedServiceTypeIndex: 0,
		interfaceType:          LocalServiceLlamaCPP,
		focusedInput:           0,
	}

	t := com.Styles

	m.baseURLInput = textinput.New()
	m.baseURLInput.Placeholder = "http://localhost:11434"
	m.baseURLInput.SetStyles(t.TextInput)
	m.baseURLInput.Focus()

	m.modelInput = textinput.New()
	m.modelInput.Placeholder = "model name (optional)"
	m.modelInput.SetStyles(t.TextInput)

	help := help.New()
	help.Styles = t.DialogHelpStyles()
	m.help = help

	m.keyMap.Next = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next field"),
	)
	m.keyMap.Previous = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous field"),
	)
	m.keyMap.Left = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "prev service type"),
	)
	m.keyMap.Right = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "next service type"),
	)
	m.keyMap.Submit = key.NewBinding(
		key.WithKeys("enter", "return"),
		key.WithHelp("enter", "submit config"),
	)
	m.keyMap.Close = CloseKey
	m.keyMap.BackToOnline = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	return m
}

// ID implements Dialog.
func (m *LocalServiceConfig) ID() string {
	return LocalServiceConfigID
}

// HandleMsg implements Dialog.
func (m *LocalServiceConfig) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Close, m.keyMap.BackToOnline):
			return ActionClose{}
		case key.Matches(msg, m.keyMap.Left):
			m.cycleServiceTypeUp()
		case key.Matches(msg, m.keyMap.Right):
			m.cycleServiceTypeDown()
		case key.Matches(msg, m.keyMap.Next):
			m.handleFieldNavigationNext()
		case key.Matches(msg, m.keyMap.Previous):
			m.handleFieldNavigationPrevious()
		case key.Matches(msg, m.keyMap.Submit):
			// Submit the config
			return m.submit()
		}
	case ActionSelectModelSource:
		// Ignore model source selection actions in this dialog
		return nil
	}

	// Handle text input messages
	switch m.focusedInput {
	case 0:
		m.baseURLInput, _ = m.baseURLInput.Update(msg)
	case 1:
		m.modelInput, _ = m.modelInput.Update(msg)
	}

	return nil
}

func (m *LocalServiceConfig) handleFieldNavigationNext() {
	m.focusedInput++
	if m.focusedInput > 1 {
		m.focusedInput = 0
	}
	m.focusInput()
}

func (m *LocalServiceConfig) handleFieldNavigationPrevious() {
	m.focusedInput--
	if m.focusedInput < 0 {
		m.focusedInput = 1
	}
	m.focusInput()
}

func (m *LocalServiceConfig) focusInput() {
	// Focus/blur inputs
	m.baseURLInput.Blur()
	m.modelInput.Blur()

	switch m.focusedInput {
	case 0:
		m.baseURLInput.Focus()
	case 1:
		m.modelInput.Focus()
	}
}

func (m *LocalServiceConfig) cycleServiceTypeDown() {
	m.selectedServiceTypeIndex++
	if m.selectedServiceTypeIndex >= len(m.serviceTypes) {
		m.selectedServiceTypeIndex = 0
	}
	m.interfaceType = m.serviceTypes[m.selectedServiceTypeIndex]
}

func (m *LocalServiceConfig) cycleServiceTypeUp() {
	m.selectedServiceTypeIndex--
	if m.selectedServiceTypeIndex < 0 {
		m.selectedServiceTypeIndex = len(m.serviceTypes) - 1
	}
	m.interfaceType = m.serviceTypes[m.selectedServiceTypeIndex]
}

func (m *LocalServiceConfig) submit() Action {
	// Validate base URL
	baseURL := strings.TrimSpace(m.baseURLInput.Value())
	if baseURL == "" {
		return util.NewErrorMsg(fmt.Errorf("base URL is required"))
	}
	
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return util.NewErrorMsg(fmt.Errorf("invalid base URL: %s", baseURL))
	}

	// Create provider config
	providerID := fmt.Sprintf("local-%s", strings.ToLower(string(m.interfaceType)))
	typeStr := string(m.interfaceType)
	providerName := "Local " + strings.ToUpper(typeStr[:1]) + strings.ToLower(typeStr[1:])

	models := []catwalk.Model{}
	if model := strings.TrimSpace(m.modelInput.Value()); model != "" {
		models = []catwalk.Model{
			{ID: model, Name: model},
		}
	}

	providerConfig := config.ProviderConfig{
		ID:                 providerID,
		Name:               providerName,
		BaseURL:            baseURL,
		Type:               catwalk.Type(m.interfaceType),
		Disable:            false,
		AutoDiscoverModels: nil,
		Models:             models,
	}

	return ActionAddLocalService{
		Type:     m.interfaceType,
		BaseURL:  baseURL,
		Model:    strings.TrimSpace(m.modelInput.Value()),
		Provider: providerConfig,
	}
}

// Draw implements Dialog.
func (m *LocalServiceConfig) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := m.com.Styles
	width := max(0, min(defaultModelsDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))

	rc := NewRenderContext(t, width)
	rc.Title = "配置本地模型服务"

	// Service type selection display
	serviceTypeOptions := []string{}
	for i, st := range m.serviceTypes {
		marker := "  "
		if i == m.selectedServiceTypeIndex {
			marker = "▶ "
		}
		serviceTypeOptions = append(serviceTypeOptions, marker+st.String())
	}
	serviceTypeText := lipgloss.JoinVertical(lipgloss.Left, serviceTypeOptions...)
	typeDisplay := fmt.Sprintf("服务类型: [按 ←/→ 切换] %s", m.interfaceType.String())
	rc.AddPart(t.Dialog.InputPrompt.Render(typeDisplay))
	rc.AddPart(t.Dialog.InputPrompt.Render(serviceTypeText))

	// Base URL input
	m.baseURLInput.SetWidth(max(0, width-t.Dialog.InputPrompt.GetHorizontalFrameSize()-1))
	baseURLView := t.Dialog.InputPrompt.Render(m.baseURLInput.View())
	rc.AddPart(baseURLView)

	// Model input
	m.modelInput.SetWidth(max(0, width-t.Dialog.InputPrompt.GetHorizontalFrameSize()-1))
	modelView := t.Dialog.InputPrompt.Render(m.modelInput.View())
	rc.AddPart(modelView)

	rc.Help = m.help.View(m)

	view := rc.Render()
	DrawCenter(scr, area, view)
	return InputCursor(t, m.baseURLInput.Cursor())
}

// ShortHelp returns the short help view.
func (m *LocalServiceConfig) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.BackToOnline,
		m.keyMap.Submit,
	}
}

// FullHelp returns the full help view.
func (m *LocalServiceConfig) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keyMap.Left, m.keyMap.Right, m.keyMap.Previous, m.keyMap.Next, m.keyMap.Submit, m.keyMap.BackToOnline},
	}
}

// ActionAddLocalService is a message to add a local service.
type ActionAddLocalService struct {
	Type      LocalServiceType
	BaseURL   string
	Model     string
	Provider  config.ProviderConfig
}

// ActionFocusProvidersList is a message to focus the providers list.
type ActionFocusProvidersList struct{}
