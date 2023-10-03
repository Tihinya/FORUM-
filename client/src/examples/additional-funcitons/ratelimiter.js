import Gachi, {
	useState,
    useEffect,
} from "../../core/framework.ts"

import ErrorPage from "./../components/errors/error-page.jsx"
import { App } from "../index.js"

const ErrorTooManyRequests = {
	message: "Too many Requests",
	status: "429",
}

export async function fetchRateLimited() {
    const response = await fetch("https://localhost:8080/ratelimited")

    if (response.status == "200") {
        return false
    }

    return true
}

export function RateLimiter() {
    const [rateLimited, setRateLimited] = useState(null)
    
    useEffect(() => {
        async function checkRateLimit() {
            const limited = await fetchRateLimited();
            setRateLimited(limited)
        }

        checkRateLimit()
    }, [])

    if (rateLimited === null) {
        return <div>Loading...</div>
    }

    if (rateLimited) {
        return <ErrorPage error={ErrorTooManyRequests} />
    }

    return (
        <App />
    )
}