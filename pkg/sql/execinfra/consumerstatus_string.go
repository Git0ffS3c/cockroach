// Code generated by "stringer"; DO NOT EDIT.

package execinfra

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NeedMoreRows-0]
	_ = x[SwitchToAnotherPortal-1]
	_ = x[DrainRequested-2]
	_ = x[ConsumerClosed-3]
}

const _ConsumerStatus_name = "NeedMoreRowsSwitchToAnotherPortalDrainRequestedConsumerClosed"

var _ConsumerStatus_index = [...]uint8{0, 12, 33, 47, 61}

func (i ConsumerStatus) String() string {
	if i >= ConsumerStatus(len(_ConsumerStatus_index)-1) {
		return "ConsumerStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ConsumerStatus_name[_ConsumerStatus_index[i]:_ConsumerStatus_index[i+1]]
}
