package lembaas

func IsValidLogin(response *AppUserAuthResponse) bool {
	return response != nil && response.UserID > 0 && (response.ExpiresIn > 0 || !response.ExpiresAt.IsZero()) && response.SessionToken != ""
}

func IsTOTPRequired(response *AppUserAuthResponse) bool {
	return response != nil && response.LoginCode != "" && !response.LoginCodeValidUntil.IsZero()
}
