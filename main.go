package main

func main() {
	go webListener(8080)
	startServer()
}