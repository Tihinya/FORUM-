export async function registrationRequest(formData) {
	const response = await fetch(`https://localhost:8080/user/create`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(formData),
	})

	const data = await response.json()

	return data
}

export async function loginRequest(formData) {
	const response = await fetch(`https://localhost:8080/login`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		credentials: "include",
		body: JSON.stringify(formData),
	})

	const data = await response.json()
	return data
}

export async function logoutRequest() {
	fetch("https://localhost:8080/logout", {
		method: "POST",
		credentials: "include",
	})
}
