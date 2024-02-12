import Gachi, {
	useContext,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api"

export function UserRole() {
	const { setUserRole } = useContext("currentUserRole")
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	useEffect(() => {
		if (isLoggin) {
			fetchUserRole()
		}
	}, [isLoggin])

	function fetchUserRole() {
		fetchData(null, `user/role`, "GET").then((resultInJson) => {
			setUserRole(resultInJson)
		})
	}
}

export function UserId() {
	const { setUserId } = useContext("currentUserId")
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	useEffect(() => {
		if (isLoggin) {
			fetchUserId()
		}
	}, [isLoggin])

	function fetchUserId() {
		fetchData(null, `user/me`, "GET").then((resultInJson) => {
			setUserId(resultInJson.id)
		})
	}
}
