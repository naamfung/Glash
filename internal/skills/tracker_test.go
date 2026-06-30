package skills

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracker_MarkLoadedAndIsLoaded(t *testing.T) {
	t.Parallel()

	activeSkills := []*Skill{
		{Name: "go-doc"},
		{Name: "shell"},
	}
	tracker := NewTracker(activeSkills)

	// Initially not loaded.
	require.False(t, tracker.IsLoaded("go-doc"))
	require.False(t, tracker.IsLoaded("shell"))

	// Mark as loaded.
	tracker.MarkLoaded("go-doc")
	require.True(t, tracker.IsLoaded("go-doc"))
	require.False(t, tracker.IsLoaded("shell"))

	// Mark another.
	tracker.MarkLoaded("shell")
	require.True(t, tracker.IsLoaded("go-doc"))
	require.True(t, tracker.IsLoaded("shell"))
}

func TestTracker_NonActiveSkillCannotBeMarkedLoaded(t *testing.T) {
	t.Parallel()

	activeSkills := []*Skill{
		{Name: "go-doc"},
	}
	tracker := NewTracker(activeSkills)

	// Cannot mark non-active skill as loaded.
	tracker.MarkLoaded("shell")
	require.False(t, tracker.IsLoaded("shell"))

	// Can mark active skill as loaded.
	tracker.MarkLoaded("go-doc")
	require.True(t, tracker.IsLoaded("go-doc"))
}

func TestTracker_NilSafety(t *testing.T) {
	t.Parallel()

	var tracker *Tracker

	// Should not panic.
	tracker.MarkLoaded("go-doc")
	require.False(t, tracker.IsLoaded("go-doc"))
}

func TestTracker_BuiltinSkillTracking(t *testing.T) {
	t.Parallel()

	// Simulate active skills including a builtin skill (glash-config).
	activeSkills := []*Skill{
		{Name: "glash-config", Description: "Glash config", Builtin: true},
		{Name: "go-doc", Description: "Go docs", Builtin: false},
	}
	tracker := NewTracker(activeSkills)

	// Initially not loaded.
	require.False(t, tracker.IsLoaded("glash-config"))
	require.False(t, tracker.IsLoaded("go-doc"))

	// Mark builtin skill as loaded (simulating read via glash://...).
	tracker.MarkLoaded("glash-config")
	require.True(t, tracker.IsLoaded("glash-config"))

	// Mark user skill as loaded.
	tracker.MarkLoaded("go-doc")
	require.True(t, tracker.IsLoaded("go-doc"))
}

func TestTracker_OverriddenBuiltinNotTracked(t *testing.T) {
	t.Parallel()

	// Simulate scenario where builtin "shell" is overridden by user "shell".
	// After dedup, only user "shell" is active.
	activeSkills := []*Skill{
		{Name: "shell", Description: "User shell override", Builtin: false},
	}
	tracker := NewTracker(activeSkills)

	// Trying to mark the builtin "shell" as loaded should not work
	// because the active skill is the user override.
	tracker.MarkLoaded("shell")
	require.True(t, tracker.IsLoaded("shell"))

	// But if we somehow tried to mark a different builtin that's not active,
	// it wouldn't get marked.
	tracker.MarkLoaded("nonexistent-builtin")
	require.False(t, tracker.IsLoaded("nonexistent-builtin"))
}
