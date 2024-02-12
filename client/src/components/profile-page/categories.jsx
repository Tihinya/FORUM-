import Gachi, {
	useContext,
	useState,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { fetchData } from "../../additional-funcitons/api"

export default function Categories() {
	const { categories, setCategories } = useContext("categories")
	const { activeSubj, setActiveSubj } = useContext("currentCategory")
	const { props } = useContext("currentProps")

	const handleCreateClick = (e) => {
		e.preventDefault()
		const form = e.target
		const formDatafied = new FormData(form)
		const formJson = Object.fromEntries(formDatafied.entries())

		fetchData(formJson, `categories`, "POST")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					fetchData(null, `categories`, "GET").then(
						(resultInJson) => {
							setCategories(resultInJson)
						}
					)
				}
			})
			.catch((error) => {
				console.error("Error: ", error)
			})

		form.reset()
	}

	function handleDeleteClick() {
		const categoryId = categories.find(
			(category) => category.category == activeSubj
		).id
		fetchData(null, `categories/${categoryId}`, "DELETE")
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					fetchData(null, `categories`, "GET").then(
						(resultInJson) => {
							setCategories(resultInJson)
						}
					)
					setActiveSubj("")
				}
			})
			.catch((error) => {
				console.error("Error: ", error)
			})
	}

	return (
		<>
			{props == "categories" ? (
				<div className="post__container">
					<form onSubmit={handleCreateClick}>
						<div className="post__box category">
							<div className="add-thread">
								<input
									type="text"
									name="category"
									maxlength="25"
									placeholder="Add a category"
								></input>
							</div>
							<div className="category__buttons">
								<button
									className="sign__button category"
									type="submit"
								>
									Create Category
								</button>
							</div>
						</div>
					</form>
					<div className="post__box category">
						<h3>
							Click a category on the left-side navigation bar to
							delete
						</h3>
						<div className="category__buttons">
							<button
								className="sign__button category delete"
								onClick={handleDeleteClick}
							>
								Delete Category
							</button>
						</div>
					</div>
				</div>
			) : null}
		</>
	)
}
