import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"
// import ErrorWindow from "../error-window/error-window.jsx"

import { loginRequest } from "../../additional-funcitons/authorization.js"

export default function Login() {
	const navigate = useNavigate()

	const [formData, setFormData] = useState({
		email: "",
		password: "",
	})
	const [errorMessage, setErrorMessage] = useState("")

	const handleInputChange = (e) => {
		const { name, value } = e.target
		setFormData((prevData) => ({
			...prevData,
			[name]: value,
		}))
	}

	const handleSubmit = (e) => {
		e.preventDefault()

		loginRequest(formData)
			.then((resultInJson) => {
				if (resultInJson.status === "success") {
					localStorage.setItem("id", resultInJson.id)
					navigate("/")
				} else if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
					console.error("Login error:", resultInJson.message)
				}
			})
			.catch((error) => {
				console.error("Error during login:", error)
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
					<div className="big_part">
						<div className="big_part_content">
							<p>Sign In</p>
							<div className="auth">
								<img src="../img/git.svg" />
								<img src="../img/goggle.svg" />
							</div>
							<h3>Or you can login with your email</h3>
							<form className="form" onSubmit={handleSubmit}>
								<div className="input-fields">
									<input
										className="input-design"
										placeholder="Email"
										type="text"
										name="email"
										value={formData.email}
										onChange={handleInputChange}
									/>
									<input
										className="input-design"
										placeholder="Password"
										type="password"
										name="password"
										value={formData.password}
										onChange={handleInputChange}
									/>
									<h3 className="forgot_pass">
										Forgot your password?
									</h3>
								</div>{" "}
								<button className="sign__button" type="submit">
									Sign In
								</button>
							</form>
						</div>
					</div>
					<div className="small_part-login">
						<div className="small_part_content">
							<p>Hello Friend!</p>
							<h6>
								Join our family and start your journey with us!
							</h6>
							<a
								className="sign__button"
								onClick={() => navigate("/registration")}
							>
								Sign Up
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
				</div>
			</div>
		</>
	)
}
