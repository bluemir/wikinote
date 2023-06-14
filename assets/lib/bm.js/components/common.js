export const css = `
/* prevent FOUC */
*:not(:defined) {
	display:none;
}
:host {
	display: block;
}
`
