package constant

import "time"

const ClassinJobUsersCacheKey string = "cache:classin:job_users"
const ClassinJobUsersCacheDelayKey string = "cache:classin:sync:delay"
const ClassinJobUsersCacheExpiry time.Duration = time.Hour * 24 * 30
