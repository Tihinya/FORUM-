import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api.js"

export default function LikesAndDislikes({ post, method, endPointUrl }) {
	const likeUrl = "user/liked"
	const disLikeUrl = "user/disliked"
	const { setPosts } = useContext("currentPosts")
	const { setErrorMessage } = useContext("currentErrorMessage")
	const { setComments } = useContext("currentComment")
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])

	const fetchLikes = () => {
		fetchData(null, likeUrl, "GET").then((resultInJson) => {
			setLikedPosts(resultInJson)
		})
	}
	const fetchDislikes = () => {
		fetchData(null, disLikeUrl, "GET").then((resultInJson) => {
			setDislikedPosts(resultInJson)
		})
	}
	const handleLike = (type, postId) => {
		fetchLikes()
		fetchDislikes()
		const isLiking =
			!likedPosts.some((obj) => obj.id === postId) &&
			!dislikedPosts.includes(postId)

		const endpoint = isLiking ? `${postId}/${type}` : `${postId}/un${type}`
		console.log(`${method}/${endpoint}`)
		fetchData(null, `${method}/${endpoint}`, "POST").then(
			(resultInJson) => {
				if (resultInJson.status === "success") {
					fetchData(null, endPointUrl, "GET").then((resultInJson) => {
						if (
							endPointUrl === "posts" ||
							endPointUrl === "user/posts" ||
							endPointUrl === "post"
						) {
							setPosts(resultInJson)
						} else {
							setComments(resultInJson)
						}
						fetchLikes()
						fetchDislikes()
					})
				} else if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
				}
			}
		)
	}

	return (
		<div className="post__likes">
			<img
				onClick={() => handleLike("like", post.id)}
				src="../img/thumbs-up.svg"
			/>
			<p onClick={() => handleLike("like", post.id)}>{post.likes}</p>
			<img
				onClick={() => handleLike("dislike", post.id)}
				src="../img/thumbs-down.svg"
			/>
			<p onClick={() => handleLike("dislike", post.id)}>
				{post.dislikes}
			</p>
		</div>
	)
}
