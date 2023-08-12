import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
}

// Add hover-over date to get full creation date

from "../core/framework.ts"
import { importCss } from "../modules/cssLoader.js"
import Button from "./button.jsx"
importCss("index.css")

export function PostContainerAuth() {
	const [posts, setPosts] = useState([])
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])

	// For displaying liked icon, if the post is already liked (TODO)
	const fetchLikedPosts = () => {
		fetch("https://localhost:8080/user/liked", {
			credentials: 'include'
		})
			.then(response => response.json())
			.then(data => setLikedPosts(data))
			.catch(error => console.error("Error fetching liked posts:", error));
	}

	const fetchDislikedPosts = () => {
		fetch("https://localhost:8080/user/disliked", {
			credentials: 'include'
		})
			.then(response => response.json())
			.then(data => setDislikedPosts(data))
			.catch(error => console.error("Error fetching liked posts:", error));
	}

	// Initialize posts/likes/dislikes upon page load
	useEffect(() => {
		// Make a GET request to fetch post data
		const fetchPosts = () => {
			fetch("https://localhost:8080/posts")
				.then(response => response.json())
				.then(data => setPosts(data))
				.catch(error => console.error("Error fetching posts:", error));
		};

		fetchPosts()
		fetchDislikedPosts()
		fetchLikedPosts()
	}, [])

	const handleLike = async (type, postId) => {
		try {
			if (!likedPosts.includes(postId) && !dislikedPosts.includes(postId)) {
				const response = await fetch(`https://localhost:8080/post/${postId}/${type}`, {
					method: 'POST',
					credentials: 'include',
				});
	
				if (response.ok) {
					setPosts(prevPosts => {
						return prevPosts.map(post => {
							if (post.id === postId) {
								if (type === 'like') {
									return { ...post, likes: post.likes + 1 };
								} else {
									return { ...post, dislikes: post.dislikes + 1 };
								}
							}
							return post;
						});
					});
					fetchLikedPosts()
					fetchDislikedPosts()
				}
			} else {
				const response = await fetch(`https://localhost:8080/post/${postId}/un${type}`, {
					method: 'POST',
					credentials: 'include',
				});
	
				if (response.ok) {
					setPosts(prevPosts => {
						return prevPosts.map(post => {
							if (post.id === postId) {
								if (type === 'like') {
									return { ...post, likes: post.likes - 1 };
								} else {
									return { ...post, dislikes: post.dislikes - 1 };
								}
							}
							return post;
						});
					});
				fetchLikedPosts()
				fetchDislikedPosts()
				}
			}
		} catch {
			console.error(`You are not logged in`)
		}
	}

	return (
		<div className="post__container">
		{posts.map((post) => {

			return (
				<div className="post__box">
				<div className="post__header">
				<div className="user__info">
					<div className="user__info_picture">
					<a href="/src/html/profile-page.html"
						><img src={post.user_info.avatar}
					/></a>
					</div>
					<div className="user__info_name">
					<p className="name">{post.user_info.username}</p>
					<p className="date">{convertTime(new Date(post.creation_date))}</p>
					</div>
				</div>
				</div>
				<div className="post__content">
				<h3>{post.title}</h3>
				<p className="post__text">
					{post.content}
				</p>
				</div>
				<div className="post__info">
				<div className="post__tags">
				{post.categories.map((category) => {
					return (
						<p className="tag">{category}</p>
					)
					})}
				</div>
				<div className="post__likes">
					<a href="/src/html/post-comment.html"
					><img src="/src/img/message-square.svg"
					/></a>
					<p>{post.comment_count}</p>
					<img src="/src/img/thumbs-up.svg" />
					<p onClick={() => handleLike("like", post.id)}>{post.likes}</p>
					<img  src="/src/img/thumbs-down.svg" />
					<p onClick={() => handleLike("dislike", post.id)}>{post.dislikes}</p>
				</div>
				</div>
			</div>
			)
		})}
		</div>
	)
}

function convertTime(creationDate) {
	const timeSinceCreation = (Date.now() - creationDate) / 1000 / 60 / 60
	switch (true) {
		case (timeSinceCreation < 1):
			const minutes = Math.floor(timeSinceCreation * 60)
			return `${minutes} minute${minutes == 1 ? "" : "s"} ago`
		case (timeSinceCreation < 24):
			const hours = Math.floor(timeSinceCreation)
			return `${hours} hour${hours == 1 ? "" : "s"} ago`
		case (timeSinceCreation < 168):
			const days = Math.floor(timeSinceCreation / 24)
			return `${days} day${days == 1 ? "" : "s"} ago`
		case (timeSinceCreation < 720):
			const weeks = Math.floor(timeSinceCreation / 24 / 7)
			return `${weeks} week${weeks == 1 ? "" : "s"} ago`
		case (timeSinceCreation < 8760):
			const months = Math.floor(timeSinceCreation / 24 / 30)
			return `${months} month${months == 1 ? "" : "s"} ago`
		case (timeSinceCreation > 8760):
			const years = Math.floor(timeSinceCreation / 24 / 365)
			return `${years} year${years == 1 ? "" : "s"} ago`
		default:
			return null
	}
}

