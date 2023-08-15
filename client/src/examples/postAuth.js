import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
}

// Add hover-over date to get full creation date

from "../core/framework.ts"
import { importCss } from "../modules/cssLoader.js"
import { convertTime } from "./helpers.js"
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
	
	const fetchPosts = () => {
		fetch("https://localhost:8080/posts")
			.then(response => response.json())
			.then(data => setPosts(data))
			.catch(error => console.error("Error fetching posts:", error));
	}

	// Initialize posts/likes/dislikes upon page load
	useEffect(() => {
		fetchPosts()
		fetchDislikedPosts()
		fetchLikedPosts()
	}, [])

	const createPost = async (title, content, categories) => { // img to be added
		try {
			const response = await fetch(`https://localhost:8080/post`, {
				method: 'POST',
				credentials: 'include',
				headers: {
					'Content-Type': 'application/json',
					'Accept': 'application/json',
				},
				body: JSON.stringify({title: title, content: content, categories: categories})
			});

			if (response.ok) {
				fetchPosts()
			} else {
				const errorData = await response.json()
				console.error(response.status, response.statusText, "-", errorData.message)
			}
		} catch {
			console.error(response.status, response.statusText, "-", errorData.message)
		}
	}

	const handleLike = async (type, postId) => {
		try {
			if (!likedPosts.includes(postId) && !dislikedPosts.includes(postId)) {
				const response = await fetch(`https://localhost:8080/post/${postId}/${type}`, {
					method: 'POST',
					credentials: 'include',
				});
	
				const errorData = await response.json()

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
				} else {
					console.error(error)
				}
			} else {
				const response = await fetch(`https://localhost:8080/post/${postId}/un${type}`, {
					method: 'POST',
					credentials: 'include',
				});

				const errorData = await response.json()
	
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
				} else {
					console.error(response.status, response.statusText, "-", errorData.message)
				}
			}
		} catch {
			console.error(response.status, response.statusText, "-", errorData.message)
		}
	}

	function handleSubmit(e) {
		// Prevent the browser from reloading the page
		e.preventDefault();
		// Read the form data
		const form = e.target;
		const formData = new FormData(form);
		const formJson = Object.fromEntries(formData.entries());
		createPost(formJson.title, formJson.content, formJson.categories) // TODO CATEGORIES
	}

	return (
		
		<div className="post__container">
		<form onSubmit={handleSubmit} className="add-thread">
          <div className="thread-button" id="add-a-thread">+</div>
          <input type="text" name="title" maxlength="120" placeholder="Add a thread" />
          <div className="thread-window" id="detailed-thread">
            <div className="thread-options">
              <div className="upload-image">
                <img src="../img/add picture.svg" />
              </div>
              <textarea
			  	className="thread-text"
                name="content"
                placeholder="Description here"
				rows={10}
              />
              <div className="thread-tags">
			  	<input type="checkbox" id="catCheckbox" className="hidden-checkbox" name="categories" value="UX/UI" />
				<label htmlFor="uxuiCheckbox" className="thread-subject">UX/UI</label>
				<button type="button" className="thread-subject">Cybersecurity</button>
                <p className="thread-subject" id="tag-active">JS</p>
                <p className="thread-subject" id="tag-active">Wisdom</p>
              </div>
            </div>
            <div className="create-post-button">
              <button className="sign__button" type="submit">Create Post</button>
            </div>
          </div>
        </form>
		{posts
			.sort((a, b) => new Date(b.creation_date) - new Date (a.creation_date))
			.map((post) => {

			return (
		<div className="post__box">
			<div className="post__header">
				<div className="user__info">
					<div className="user__info_picture">
					<a href="../html/profile-page.html"
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
					<a href="../html/post-comment.html"
					><img src="../img/message-square.svg"
					/></a>
					<p>{post.comment_count}</p>
					<img src="../img/thumbs-up.svg" />
					<p onClick={() => handleLike("like", post.id)}>{post.likes}</p>
					<img  src="../img/thumbs-down.svg" />
					<p onClick={() => handleLike("dislike", post.id)}>{post.dislikes}</p>
				</div>
				</div>
			</div>
			)
		})}
		</div>
	)
}
