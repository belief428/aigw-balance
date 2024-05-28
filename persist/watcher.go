package persist

type IWatcher interface {
	SetExecCallback() error
}
