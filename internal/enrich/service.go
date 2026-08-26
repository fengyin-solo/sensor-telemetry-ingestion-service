package enrich

import "telemetry/internal/config"

func Apply(profile config.Profile, value int) (config.Profile, error) {
	profile.Labels["source"] = "ingestion"
	if profile.Validator != nil {
		if err := profile.Validator.Validate(value); err != nil {
			return profile, err
		}
	}
	return profile, nil
}
