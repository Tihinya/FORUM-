import Gachi, { useState, useNavigate } from "../../../core/framework"
import { fetchData } from "../../additional-funcitons/api.js"

export default function Registation() {
	const navigate = useNavigate()
	const [errorMessage, setErrorMessage] = useState("")
	const registration = "user/create"

	const [formData, setFormData] = useState({
		email: "",
		username: "",
		password: "",
		password_confirmation: "",
	})

	const handleInputChange = (e) => {
		const { name, value } = e.target
		setFormData((prevData) => ({
			...prevData,
			[name]: value,
		}))
	}

	const handleSubmit = (e) => {
		e.preventDefault()
		fetchData(formData, registration, "Post").then((resultInJson) => {
			if (resultInJson.status === "success") {
				localStorage.setItem("id", resultInJson.id)
				navigate("/")
			} else if (resultInJson.status === "error") {
				setErrorMessage(resultInJson.message)
				console.error("Registration error:", resultInJson.message)
			}
		})
	}

	return (
		<>
			{errorMessage != "" ? (
				<div className="error-window">
					<button
						onClick={() => {
							setErrorMessage("")
						}}
					>
						[X]
					</button>
					{errorMessage}
				</div>
			) : (
				""
			)}
			<div className="main__block">
				<div className="sign-up__block">
					<div className="small_part-sign-up">
						<div className="small_part_content">
							<p>Already have an account?</p>
							<h6>Sign in for better experience!</h6>
							<a
								className="sign__button"
								onClick={() => navigate("/login")}
							>
								Sign In
							</a>
							<h6>Or return to</h6>
							<a
								className="sign__button"
								onClick={() => navigate("/")}
							>
								Main Page
							</a>
						</div>
					</div>
					<div className="big_part">
						<div className="big_part_content">
							<p>Sign Up</p>
							<div className="auth">
								<img src="../img/git.svg" />
								<img src="../img/goggle.svg" />
							</div>
							<h3>Sign in for better experience!</h3>
							<form className="form" onSubmit={handleSubmit}>
								<div className="input-fields">
									<input
										name="email"
										className="input-design"
										type="text"
										placeholder="Email"
										value={formData.email}
										onChange={handleInputChange}
									/>
									<input
										name="username"
										className="input-design"
										type="text"
										placeholder="Nickname"
										value={formData.username}
										onChange={handleInputChange}
									/>
									<input
										name="password"
										className="input-design"
										placeholder="Password"
										type="password"
										value={formData.password}
										onChange={handleInputChange}
									/>
									<input
										name="password_confirmation"
										className="input-design"
										placeholder="Repeat Password"
										type="password"
										value={formData.password_confirmation}
										onChange={handleInputChange}
									/>
								</div>
								<button className="sign__button" type="submit">
									Sign Up
								</button>
							</form>
						</div>
					</div>
				</div>
			</div>
		</>
	)
}
