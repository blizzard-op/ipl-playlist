package playlist

import (
	"time"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

type Playlist struct {
	StartsAt, EndsAt time.Time
	Config yaml.File
}