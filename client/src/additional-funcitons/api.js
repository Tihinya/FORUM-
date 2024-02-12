export function fetchData(formData = null, endPointUrl, method) {
	const frontendHost = process.env.FRONTEND_HOST
	console.log(`Frontend Host: ${frontendHost}`)

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

	return fetch(`${frontendHost}/${endPointUrl}`, requestOptions).then(
		(response) => {
			if (response.headers.get("content-length") === "0") {
				return null
			}
			return response.json()
		}
	)
}
