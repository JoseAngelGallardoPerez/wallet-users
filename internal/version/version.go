package version

var (
	// DATE returns the build date
	DATE = "UNKNOWN"
	// TAG returns the git commit TAG
	TAG = "UNKNOWN"
	// COMMIT returns the sha from git
	COMMIT = "UNKNOWN"
)

type info struct {
	Date   *string
	Tag    *string
	Commit *string
}

var BuildInfo = info{
	Date:   &DATE,
	Tag:    &TAG,
	Commit: &COMMIT,
}
