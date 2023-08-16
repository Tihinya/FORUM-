import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { convertTime } from "../../additional-funcitons/postTime"

export default function Posts() {
	const navigate = useNavigate()

	const [posts, setPosts] = useState([])

	useEffect(() => {
		const interval = setInterval(() => {
			fetch("https://localhost:8080/posts")
				.then((response) => response.json())
				.then((data) => setPosts(data))
				.catch((error) => console.error("Error fetching users:", error))
		}, 1000)

		return () => clearInterval(interval)
	}, [])

	return (
		<div className="post__container">
			{posts.map((post) => (
				<div className="post__box">
					<div className="post__header">
						<div className="user__info">
							<div className="user__info_picture">
								<a onClick={() => navigate("/profile-page")}>
									<img src="../img/avatarka.jpeg" />
								</a>
							</div>
							<div className="user__info_name">
								<p className="name">{post.user_info.username}</p>
								<p className="date">{convertTime(post.creation_date)}</p>
							</div>
						</div>
					</div>
					<div className="post__content">
						<h3>{post.title}</h3>
						<p className="post__text">{post.content}</p>
					</div>
					<div className="post__info">
						<div className="post__tags">
							{post.categories.map((categories) => (
								<p className="tag">{categories}</p>
							))}
						</div>
						<div className="post__likes">
							<a onClick={() => navigate("/post-comment")}>
								<img src="../img/message-square.svg" />
							</a>
							<p>3</p>
							<img src="../img/thumbs-up.svg" />
							<p>{post.likes}</p>
							<img src="../img/thumbs-down.svg" />
							<p>{post.dislikes}</p>
						</div>
					</div>
				</div>
			))}
		</div>
	)
}
