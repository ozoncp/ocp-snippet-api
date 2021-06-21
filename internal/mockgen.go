package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/ozoncp/ocp-snippet-api/internal/repo Repo

//go:generate mockgen -destination=./mocks/api_mock.go -package=mocks github.com/ozoncp/ocp-snippet-api/internal/api Api
