import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

import Header from "../header/header"
import DropdownMenu from "../header/dropdown"
import { NavBar } from "../navbar/navbar"
import Posts from "../posts/posts"

export default function MainPage() {
	const navigate = useNavigate()

	return (
		<div>
			<Header />
			{/* <DropdownMenu /> */}
			<NavBar />
			<Posts />
		</div>
	)
}
