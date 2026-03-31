package editorApi

type EditorApi interface {
	SendQuitSignal()
	ToggleCommandPrompt(active bool)
}
