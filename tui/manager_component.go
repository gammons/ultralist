package tui

// ManagerComponent is a component within the tui manager.
// Some ManagerComponents process their own input events with input fields.  In this case, we need to tell the global manager to stop processing global events (e.g. hitting "q" to quit).
type ManagerComponent interface {
	ProcessGlobalEvents() bool
}
