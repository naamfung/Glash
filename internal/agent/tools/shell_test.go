package tools

import (
	"context"
	"encoding/json"
	"testing"

	"charm.land/fantasy"
	"github.com/charmbracelet/glash/internal/config"
	"github.com/charmbracelet/glash/internal/permission"
	"github.com/charmbracelet/glash/internal/pubsub"
	"github.com/charmbracelet/glash/internal/shell"
	"github.com/stretchr/testify/require"
)

type mockShellPermissionService struct {
	*pubsub.Broker[permission.PermissionRequest]
}

func (m *mockShellPermissionService) Request(ctx context.Context, req permission.CreatePermissionRequest) (bool, error) {
	return true, nil
}

func (m *mockShellPermissionService) Grant(req permission.PermissionRequest) bool { return true }

func (m *mockShellPermissionService) Deny(req permission.PermissionRequest) bool { return true }

func (m *mockShellPermissionService) GrantPersistent(req permission.PermissionRequest) bool {
	return true
}

func (m *mockShellPermissionService) AutoApproveSession(sessionID string) {}

func (m *mockShellPermissionService) SetSkipRequests(skip bool) {}

func (m *mockShellPermissionService) SkipRequests() bool {
	return false
}

func (m *mockShellPermissionService) SubscribeNotifications(ctx context.Context) <-chan pubsub.Event[permission.PermissionNotification] {
	return make(<-chan pubsub.Event[permission.PermissionNotification])
}

func TestShellTool_DefaultAutoBackgroundThreshold(t *testing.T) {
	workingDir := t.TempDir()
	tool := newShellToolForTest(workingDir)
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	resp := runShellTool(t, tool, ctx, ShellParams{
		Description: "default threshold",
		Command:     "echo done",
	})

	require.False(t, resp.IsError)
	var meta ShellResponseMetadata
	require.NoError(t, json.Unmarshal([]byte(resp.Metadata), &meta))
	require.False(t, meta.Background)
	require.Empty(t, meta.ShellID)
	require.Contains(t, meta.Output, "done")
}

func TestShellTool_CustomAutoBackgroundThreshold(t *testing.T) {
	workingDir := t.TempDir()
	tool := newShellToolForTest(workingDir)
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	resp := runShellTool(t, tool, ctx, ShellParams{
		Description:         "custom threshold",
		Command:             "sleep 1.5 && echo done",
		AutoBackgroundAfter: 1,
	})

	require.False(t, resp.IsError)
	var meta ShellResponseMetadata
	require.NoError(t, json.Unmarshal([]byte(resp.Metadata), &meta))
	require.True(t, meta.Background)
	require.NotEmpty(t, meta.ShellID)
	require.Contains(t, resp.Content, "moved to background")

	bgManager := shell.GetBackgroundShellManager()
	require.NoError(t, bgManager.Kill(meta.ShellID))
}

type recordingPermissionService struct {
	*pubsub.Broker[permission.PermissionRequest]
	requestCount int
	allow        bool
}

func (m *recordingPermissionService) Request(ctx context.Context, req permission.CreatePermissionRequest) (bool, error) {
	m.requestCount++
	return m.allow, nil
}

func (m *recordingPermissionService) Grant(req permission.PermissionRequest) bool { return true }

func (m *recordingPermissionService) Deny(req permission.PermissionRequest) bool { return true }

func (m *recordingPermissionService) GrantPersistent(req permission.PermissionRequest) bool {
	return true
}

func (m *recordingPermissionService) AutoApproveSession(sessionID string) {}

func (m *recordingPermissionService) SetSkipRequests(skip bool) {}

func (m *recordingPermissionService) SkipRequests() bool {
	return false
}

func (m *recordingPermissionService) SubscribeNotifications(ctx context.Context) <-chan pubsub.Event[permission.PermissionNotification] {
	return make(<-chan pubsub.Event[permission.PermissionNotification])
}

func newShellToolForTest(workingDir string) fantasy.AgentTool {
	permissions := &mockShellPermissionService{Broker: pubsub.NewBroker[permission.PermissionRequest]()}
	attribution := &config.Attribution{TrailerStyle: config.TrailerStyleNone}
	return NewShellTool(permissions, workingDir, attribution, "test-model")
}

func newShellToolWithRecordingPerms(workingDir string, allow bool) (fantasy.AgentTool, *recordingPermissionService) {
	perms := &recordingPermissionService{
		Broker: pubsub.NewBroker[permission.PermissionRequest](),
		allow:  allow,
	}
	attribution := &config.Attribution{TrailerStyle: config.TrailerStyleNone}
	return NewShellTool(perms, workingDir, attribution, "test-model"), perms
}

func TestShellTool_ChainedCommandsRequirePermission(t *testing.T) {
	workingDir := t.TempDir()
	tool, perms := newShellToolWithRecordingPerms(workingDir, true)
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	// ls && echo should trigger permission check.
	resp := runShellTool(t, tool, ctx, ShellParams{
		Description: "chained ls",
		Command:     "ls && echo done",
	})

	require.False(t, resp.IsError)
	require.Equal(t, 1, perms.requestCount, "chained command should trigger permission request")

	// Plain ls should NOT trigger permission check.
	perms.requestCount = 0
	resp = runShellTool(t, tool, ctx, ShellParams{
		Description: "plain ls",
		Command:     "ls -la",
	})

	require.False(t, resp.IsError)
	require.Equal(t, 0, perms.requestCount, "plain ls should not trigger permission request")
}

func TestShellTool_ChainedCommandsDenied(t *testing.T) {
	workingDir := t.TempDir()
	tool, perms := newShellToolWithRecordingPerms(workingDir, false)
	ctx := context.WithValue(context.Background(), SessionIDContextKey, "test-session")

	resp := runShellTool(t, tool, ctx, ShellParams{
		Description: "chained ls denied",
		Command:     "ls && rm -rf /",
	})

	require.Equal(t, 1, perms.requestCount)
	require.Contains(t, resp.Content, "User denied permission")
}

func runShellTool(t *testing.T, tool fantasy.AgentTool, ctx context.Context, params ShellParams) fantasy.ToolResponse {
	t.Helper()

	input, err := json.Marshal(params)
	require.NoError(t, err)

	call := fantasy.ToolCall{
		ID:    "test-call",
		Name:  ShellToolName,
		Input: string(input),
	}

	resp, err := tool.Run(ctx, call)
	require.NoError(t, err)
	return resp
}
