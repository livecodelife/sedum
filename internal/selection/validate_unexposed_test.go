package selection

import (
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
)

func hiddenHelperInvocation(kwargs map[string]any) []recording.Invocation {
	return []recording.Invocation{{Action: "hiddenHelper", Kwargs: kwargs}}
}

// A recording written by an adopter names actions its package hides from the
// model. Validate takes them only when told to.
func TestValidateRejectsAnUnexposedActionByDefault(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	violations := Validate(files, hiddenHelperInvocation(map[string]any{"controller": "users"}), false)

	if len(violations) == 0 {
		t.Fatal("an unexposed action was accepted")
	}
	if violations[0].Rule != RuleUnknownAction {
		t.Errorf("rule = %v, want %v", violations[0].Rule, RuleUnknownAction)
	}
}

func TestValidateAcceptsAnUnexposedActionWhenAllowed(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	if violations := Validate(files, hiddenHelperInvocation(map[string]any{"controller": "users"}), true); len(violations) > 0 {
		t.Errorf("an unexposed action was rejected although allowed: %v", violations)
	}
}

func TestValidateAllowingUnexposedActionsKeepsTheOtherRules(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	violations := Validate(files, hiddenHelperInvocation(map[string]any{}), true)

	if len(violations) == 0 {
		t.Fatal("an unexposed action missing a required kwarg was accepted")
	}
	if violations[0].Rule == RuleUnknownAction {
		t.Errorf("rule = %v, want a kwarg rule", violations[0].Rule)
	}
}
