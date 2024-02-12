import Gachi, {
	useState,
	useEffect,
} from "../../Gachi.js/src/core/framework.ts"

import ErrorPage from "./../components/errors/error-page.jsx"
import { App } from "../index.js"
import { fetchData } from "./api.js"

const ErrorTooManyRequests = {
	message: "Too many Requests",
	status: "429",
}

// let baseURL = "https://ec2-51-20-1-125.eu-north-1.compute.amazonaws.com:8080"
// const frontendHost = process.env.FRONTEND_HOST

// export async function fetchRateLimited() {
// 	const response = await fetch(frontendHost + "/ratelimited")

// 	if (response.status == "200") {
// 		return false
// 	}

// 	return true
// }

// export function fetchRateLimited() {}

export function RateLimiter() {
	const [rateLimited, setRateLimited] = useState(null)

	useEffect(() => {
		async function checkRateLimit() {
			fetchData(null, `ratelimited`, "GET").then((resultInJson) => {
				resultInJson.status == "200"
					? setRateLimited(false)
					: setRateLimited(true)
			})

			// const limited = await fetchRateLimited()
			// setRateLimited(limited)
		}

		checkRateLimit()
	}, [])

	if (rateLimited === null) {
		return <div>Loading...</div>
	}

	if (rateLimited) {
		return <ErrorPage error={ErrorTooManyRequests} />
	}

	return <App />
}
