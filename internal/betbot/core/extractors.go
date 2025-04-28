package core

import "errors"

func extractFirstPoints(plays []play) (play, error) {
	for _, p := range plays {
		if p.ScoringPlay {
			return p, nil
		}
	}
	return play{}, errors.New("scoring play not found")
}

func extractFirstShotAttempt(plays []play) (play, error) {
	for _, p := range plays {
		if p.ShootingPlay {
			return p, nil
		}
	}
	return play{}, errors.New("shooting play not found")
}
