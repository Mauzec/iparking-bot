package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mauzec/ibot-things/pkg/retry"
)

type DataPayload struct {
	Distance float64 `json:"distance"`
}

func ReadDistance(path string) (float64, error) {
	var payload DataPayload

	err := retry.Do(3, 100*time.Millisecond, func() error {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("[ReadDistance] %w: %s", err, path)
		}
		if err := json.Unmarshal(b, &payload); err != nil {
			return fmt.Errorf("[ReadDistance] %w: %s", err, path)
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to read distance: %w", err)
	}
	return payload.Distance, nil
}
