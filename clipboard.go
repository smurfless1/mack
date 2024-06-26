/*
** Mack: Clipboard
** Interact with the clipboard
 */

package mack

// Clipboard returns the content of the clipboard
func Clipboard() (string, error) {
	str, err := run(buildClipboard())
	// de-quote it on return
	return str[1 : len(str)-1], err
}

// Build the command
func buildClipboard() string {
	return build("the clipboard")
}

// SetClipboard changes the content of the clipboard
func SetClipboard(content string) error {
	_, err := run(buildSetClipboard(content))
	return err
}

// Wrap the content in quotes and build the command
func buildSetClipboard(content string) string {
	content = wrapInQuotes(content)
	return build("set the clipboard to", content)
}
