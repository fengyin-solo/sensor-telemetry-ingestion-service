package enrich

import "telemetry/internal/config"

func Apply(profile config.Profile, value int) (config.Profile, error) {
	labels := make(map[string]string, len(profile.Labels)+1)
	for key, label := range profile.Labels {
		labels[key] = label
	}
	profile.Labels = labels
	profile.Labels["source"] = "ingestion"
	if config.ValidatorAvailable(profile.Validator) {
		if err := profile.Validator.Validate(value); err != nil {
			return profile, err
		}
	}
	return profile, nil
}
