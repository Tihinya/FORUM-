import Gachi, {
	useState,
	useNavigate,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api.js"

import gitSVG from "../../img/git.svg"
import googleSVG from "../../img/google.svg"

export default function Registation() {
	const navigate = useNavigate()
	const [errorMessage, setErrorMessage] = useState("")
	const registration = "user/create"

	const [formData, setFormData] = useState({
		email: "",
		username: "",
		password: "",
		password_confirmation: "",
		gender: "",
		age: "",
	})
	const [loginGoogleUrl, setLoginGoogleUrl] = useState("")
	const [loginGihubUrl, setLoginGithubUrl] = useState("")

	const handleInputChange = (e) => {
		const { name, value } = e.target

		if (e.target.tagName === "SELECT") {
			setFormData((prevData) => ({
				...prevData,
				[name]: value,
			}))
		} else {
			setFormData((prevData) => ({
				...prevData,
				[name]: value,
			}))
		}
	}

	const handleSubmit = (e) => {
		e.preventDefault()
		fetchData(formData, registration, "POST").then((resultInJson) => {
			if (resultInJson.status === "success") {
				localStorage.setItem("id", resultInJson.id)
				navigate("/login")
			} else if (resultInJson.status === "error") {
				setErrorMessage(resultInJson.message)
				console.error("Registration error:", resultInJson.message)
			}
		})
	}
	loginWithServises("github")
	loginWithServises("google")

	const ageOptions = []
	for (let i = 8; i <= 55; i++) {
		ageOptions.push(
			<option key={i} value={i}>
				{i}
			</option>
		)
	}

	function loginWithServises(servisName) {
		fetchData(null, `login/${servisName}`, "GET").then((responseInJson) => {
			if (servisName === "github") {
				setLoginGithubUrl(responseInJson)
			} else if (servisName === "google") {
				setLoginGoogleUrl(responseInJson)
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
								<a
									href={
										loginGihubUrl != ""
											? loginGihubUrl
											: undefined
									}
								>
									<img src={gitSVG} />
								</a>

								<a
									href={
										loginGoogleUrl != ""
											? loginGoogleUrl
											: undefined
									}
								>
									<img src={googleSVG} />
								</a>
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
									<select
										name="gender"
										className="input-design"
										value={formData.gender}
										onChange={handleInputChange}
									>
										<option value="">Select gender</option>
										<option value="male">Male</option>
										<option value="female">Female</option>
										<option value="helicopter">
											Helicopter
										</option>
									</select>
									<select
										name="age"
										className="input-design"
										value={formData.age}
										onChange={handleInputChange}
									>
										{ageOptions}
									</select>
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
