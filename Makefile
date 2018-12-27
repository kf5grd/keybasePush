keybasePush : main.go repo.go web_handlers.go \
              web_logger.go web_message.go \
              web_router.go web_routes.go
	go build -o keybasePush
clean :
	rm keybasePush
