import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function PostFilter() {
	const [filter, setFilter] = useState({
		title: "",
		description: "",
	})

	const filteredPosts = posts.filter((post) => {
		return (
			post.title.toLowerCase().includes(filter.title.toLowerCase()) &&
			post.content
				.toLowerCase()
				.includes(filter.description.toLowerCase())
		)
	})
	return (
		<div>
			<input
				type="text"
				placeholder="Search by title"
				value={filter.title}
				onChange={(e) =>
					setFilter({ ...filter, title: e.target.value })
				}
			/>
		</div>
	)
}
