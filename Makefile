postgres:
		docker run --name user-management-system-go -e POSTGRES_USER=cisco -e POSTGRES_PASSWORD=secret -e POSTGRESS_DB=fiber-jwt-user-management-system -p 5432:5432 -d postgres
.PHONY: postgres