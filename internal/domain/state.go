package domain

import "github.com/zhanglei10281852-gif/mining-eco-restoration/internal/apperror"

var transitions = map[string]map[string]bool{
	"planned": {"assigned": true}, "assigned": {"in_progress": true}, "in_progress": {"submitted": true},
	"submitted": {"accepted": true, "rejected": true}, "rejected": {"in_progress": true}, "accepted": {},
}

func CanTransition(from, to string) bool { return transitions[from][to] }
func Transition(from, to string) error {
	if !CanTransition(from, to) {
		return apperror.ErrConflict
	}
	return nil
}
func TaskStatuses() []string {
	return []string{"planned", "assigned", "in_progress", "submitted", "accepted", "rejected"}
}
