package srv

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/db/gorm"
)

const (
	cacheTTL = 5 * time.Minute
)

var (
	selectOptionFunc gorm.SelectOptionFunc
)
