import Gachi, {
	useNavigate,
	useContext,
	useState,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { NavBar } from "../navbar/navbar"
import Posts from "../posts/posts"
import Header from "../header/header"
import PersonalNavBar from "../personalNavBar/personalNavBar"
import Promotions from "./promotions"
import Demotions from "./demotions"
import Categories from "./categories"
import { fetchData } from "../../additional-funcitons/api"
import ReportWindow from "./reportWindow.jsx"
import ReportAnswer from "./reportAnswer.jsx"

export default function ProfilePage() {
	// const defaultposts = createContext("user/posts")
	const { props } = useContext("currentProps")
	const [promotions, setPromotions] = useState([])
	Gachi.createContext("promotions", { promotions, setPromotions })
	const { categories, setCategories } = useContext("categories")
	const { userRole } = useContext("currentUserRole")
	const { displayNavbar, setDisplayNavbar } = useContext("displayNavbar")
	const { moderators, setModerators } = useContext("currentModerators")

	useEffect(() => {
		fetchData(null, `categories`, "GET").then((resultInJson) => {
			setCategories(resultInJson)
		})
	}, [])

	useEffect(() => {
		if (userRole === "admin") {
			fetchData(null, `users/get`, "GET").then((usersInJson) => {
				const moderators = usersInJson.filter(
					(user) => user.role_id == 2
				) // Filter users who have moderator role
				setModerators(moderators)
			})

			fetchData(null, `promotions/get`, "GET").then((resultInJson) => {
				setPromotions(resultInJson)
			})
		}

		if (props !== "demotions") {
			setDisplayNavbar(true)
		}
	}, [props])

	return (
		<div>
			<Header />
			<NavBar />
			<PersonalNavBar />
			{props === "promotions" ? (
				<Promotions />
			) : props === "demotions" ? (
				<Demotions />
			) : props === "postreport/get" ? (
				<ReportWindow endPointUrl={props} />
			) : props === "postreport/answer/get" ? (
				<ReportAnswer endPointUrl={props} />
			) : props === "categories" ? (
				<Categories />
			) : (
				<Posts endPointUrl={props} userId={""} />
			)}
		</div>
	)
}
