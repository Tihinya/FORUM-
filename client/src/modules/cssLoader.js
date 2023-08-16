export function importCss(link) {
	// console.log(document)
	const temp = document.createElement("link")
	temp.setAttribute("rel", "stylesheet")
	temp.setAttribute("type", "text/css")
	temp.setAttribute("href", link)

	document.head.appendChild(temp)
	// console.log(temp)
}
