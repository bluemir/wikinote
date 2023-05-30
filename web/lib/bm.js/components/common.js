export const css = `
@import '/static/lib/icons/icons.css';

/* prevent FOUC */
*:not(:defined) {
	display:none;
}
:host {
	display: block;
}
`
