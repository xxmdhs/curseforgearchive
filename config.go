package main

type config struct {
	SelID []idConfig
}

type idConfig struct {
	ID        int
	Page      int
	Reacquire bool
}
