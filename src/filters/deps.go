package filters

type logger interface {
	Warnf(format string, args ...interface{})
}
