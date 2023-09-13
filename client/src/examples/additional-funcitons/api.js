export function fetchData(formData = null, endPointUrl, method) {
	// console.log(endPointUrl)
	const requestOptions = {
		method: method,
		headers: {
			"Content-Type": "application/json",
		},
		credentials: "include",
	}

	if (formData !== null) {
		requestOptions.body = JSON.stringify(formData)
	}

	return fetch(`http://localhost:8080/${endPointUrl}`, requestOptions)
		.then((response) => {
			if (!response.ok) {
				throw new Error(
					`Request failed with status: ${response.status}`
				)
			}
			return response.json()
		})
		.catch((error) => {
			// Handle any network or fetch-related errors here
			throw new Error("Failed to fetch data: " + error.message)
		})
}
