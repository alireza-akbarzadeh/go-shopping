package controllers

type Container struct {
	Health  *HealthController
	Auth    *AuthController
	Profile *ProfileController
}
