package constant

import "time"

const (
	AWS_S3          = "AWS_S3"
	PUBLIC          = "public"
	PRIVATE         = "private"
	STORAGE_DEFAULT = AWS_S3
)

const (
	EXPIRY_MEDIA_DURATION = time.Minute * 10
)
