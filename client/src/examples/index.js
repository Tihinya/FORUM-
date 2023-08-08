import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
} from "../core/framework.ts"
import { importCss } from "../modules/cssLoader.js"
import Button from "./button.jsx"
importCss("./index.css")

const container = document.getElementById("root")

function App() {
	const [users, setUsers] = useState([])

	useEffect(() => {
		// Make a GET request to fetch user data
		fetch("https://localhost:8080/users/get")
			.then((response) => response.json())
			.then((data) => setUsers(data))
			.catch((error) => console.error("Error fetching users:", error))
	}, [])

	return (
		<div>
			<h1>User List</h1>
			<ul>
				{users.map((user) => (
					<li key={user.id}>
						<p>Username: {user.username}</p>
						<p>Email: {user.email}</p>
					</li>
				))}
			</ul>
		</div>
	)
}

Gachi.render(<App />, container)
