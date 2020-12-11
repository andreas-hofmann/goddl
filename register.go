package main

type RegisterResult struct {
	Success struct {
		Username string
	}
	Error struct {
		Address     string
		Description string
		Type        int
	}
}
