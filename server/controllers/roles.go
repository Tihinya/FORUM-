package controllers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/router"
	"forum/validation"
	"net/http"
)

func RequestPromotion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req database.GetRoleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	//Get User ID
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return

	}
	//Get Role ID
	roleId, err := database.GetRoleId(req.NewRole)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
	}

	if req.NewRole == "" {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
	}

	// Check role
	if req.NewRole != "" {
		// Check if role exist
		exist, err := validation.ValidateRole(database.DB, req.NewRole)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !exist {
			ReturnMessageJSON(w, "Role doesn't exist ", http.StatusBadRequest, "error")
			return
		}

		exist, err = validation.HasPendingRoleRequest(database.DB, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exist {
			ReturnMessageJSON(w, "You already have a request pending", http.StatusBadRequest, "error")
			return
		}

	}

	err = database.CreateRoleRequest(userID, roleId)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	ReturnMessageJSON(w, "Request has been sent", http.StatusOK, "success")

}

func ReadPromotions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqs = make([]database.SendRoleRequest, 0)
	roleRequests, err := database.SelectAllPromotion()
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	for _, roleRequest := range roleRequests {
		// Get user Name
		userName, err := validation.GetUserName(database.DB, roleRequest.UserID)
		if err != nil {
			ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
			return
		}

		// Get role Name
		roleName, err := validation.GetRoleName(database.DB, roleRequest.NewRoleID)
		if err != nil {
			ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
			return
		}

		// Create a SendRoleRequest instance with UserName and RoleName
		req := database.SendRoleRequest{
			UserID:   roleRequest.UserID,
			UserName: userName,
			RoleName: roleName,
		}

		// Append the request to the list
		reqs = append(reqs, req)
	}

	json.NewEncoder(w).Encode(reqs)

}

func PromoteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	// Get status string
	status, err := router.GetFieldString(r, "response")
	if err != nil {
		fmt.Println(err)
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	if status == "" || userID == 0 {
		ReturnMessageJSON(w, "Invalid request body, body cannot be empty and user id cannot equal 0", http.StatusBadRequest, "error")
	}
	if status != "accept" && status != "decline" {
		ReturnMessageJSON(w, "Invalid request body, body needs parameter 'accept' or 'decline'", http.StatusBadRequest, "error")
	}
	// Сhecking RoleRequest
	exist, err := validation.HasPendingRoleRequest(database.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		ReturnMessageJSON(w, "No pending request", http.StatusBadRequest, "error")
		return
	}

	if status == "accept" {
		// Get RoleID
		roleId, err := database.GetNewRoleIDByUserID(userID)
		if err != nil {
			ReturnMessageJSON(w, "Error fetching role", http.StatusInternalServerError, "error")
			return
		}
		// Update User Role
		err = database.UpdateUser("", "", roleId, userID)
		if err != nil {
			ReturnMessageJSON(w, "Error updating user", http.StatusInternalServerError, "error")
			return
		}
		// Delete RoleRequest
		err = database.DeleteRoleRequest(userID)
		if err != nil {
			ReturnMessageJSON(w, "Error declining role request", http.StatusInternalServerError, "error")
			return
		}
		ReturnMessageJSON(w, "Application approved", http.StatusOK, "success")
	}
	if status == "decline" {
		// Delete RoleRequest
		err := database.DeleteRoleRequest(userID)
		if err != nil {
			ReturnMessageJSON(w, "Error declining role request", http.StatusInternalServerError, "error")
			return
		}
		ReturnMessageJSON(w, "Application denied", http.StatusOK, "success")
	}

}

func DemoteUser(w http.ResponseWriter, r *http.Request) {
	var req database.GetRoleRequest
	//Get user ID
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	role, err := database.GetRoleId(req.NewRole)
	if err != nil {
		ReturnMessageJSON(w, "Error fetching role id", http.StatusInternalServerError, "error")
		return
	}

	// Delete RoleRequest
	err = database.DemoteUser(role, userID)
	if err != nil {
		ReturnMessageJSON(w, "Error demoting user", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "User demoted", http.StatusOK, "success")
}
