package telemetry

import "github.com/mika-rafie/alamino/internal/types"

func StepName(name types.StepName) string {
	if name.IsCustomGate() {
		return "gate"
	}
	return string(name)
}
