package main

func main() {
	/*
		This is a proposed endpoint design that groups actions by the data type they operate on
		because real-time-forum and social-network projects are going to be API based
		It's important to note that this is just a design and the actual implementation may vary.

		Logging, Notification, Access Level, Status, Security, Performance will be implemented as middleware on top
		of the enpoints.
	*/

	// Post:
	// --Create
	// --Read
	// --Update
	// --Delete

	// User:
	// --Create
	// --Read
	// --Update
	// --Delete
	// --Change permissions

	// Login:
	// --Login
	// --Logout
	// ( register is User:Create )

	// Comment:
	// --Create
	// --Read
	// --Update
	// --Delete

	// Pages:
	// --Main
	// --Loggin/registration
	// --Profile
	// --Error page
	// --Performance limit page
}
