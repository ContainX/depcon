package compose

type Context struct {
	ComposeFile          string
	ProjectName          string
	EnvParams            map[string]string
	ErrorOnMissingParams bool
}
