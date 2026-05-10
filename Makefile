.PHONY: dev-backend dev-frontend test build quality-gate

dev-backend:
	cd backend && PORT=3000 go run main.go

dev-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...

build:
	cd backend && go build -o viddl-server main.go
	cd frontend && npm run build

quality-gate: test build
