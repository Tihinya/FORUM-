export function convertTime(creationDate) {
	const timeSinceCreation = (Date.now() - new Date(creationDate)) / 1000 / 60 / 60
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