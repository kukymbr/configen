package logger

var opt Options

type Options struct {
	Silent bool
	Debug  bool
}

func SetOptions(o Options) {
	opt = o
}
