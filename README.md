# Project personal_budget_app

## Routes
`router.HandleFunc("/login", s.handleLogin).Methods("POST")`

`router.HandleFunc("/register", s.handleCreateAccount).Methods("POST")`

`// Recover pass`

`router.HandleFunc("/accounts/forgetpw", s.handleForgetPassword).Methods("POST")`

`router.HandleFunc("/accounts/reset-password", s.handlePasswordReset).Methods("POST")`

`// Protected routes`
`secure := router.PathPrefix("/api").Subrouter()`

`// secure.Use(requestLoggerMiddleware) // for debug`
`secure.Use(JWTMiddleware)`

`secure.HandleFunc("/logout", s.handleLogout).Methods("GET", "POST")`

`secure.HandleFunc("/accounts", s.handleGetAccounts).Methods("GET")`

`secure.HandleFunc("/accounts/{id}", s.handleGetAccount).Methods("GET")`

`secure.HandleFunc("/accounts/{id}", s.handleDeleteAccount).Methods("DELETE")`

`secure.HandleFunc("/accounts/{id}", s.handleUpdateAccount).Methods("PUT")`

`secure.HandleFunc("/cards", s.handleAddCard).Methods("POST")`

`secure.HandleFunc("/cards", s.handleGetCards).Methods("GET")`

`secure.HandleFunc("/cards/{id}", s.handleDeleteCard).Methods("DELETE")`

`secure.HandleFunc("/cards/{id}", s.handleGetCard).Methods("GET")`

`secure.HandleFunc("/transaction/{cardId}", s.handleGetTransactions).Methods("GET")`

`secure.HandleFunc("/transaction", s.handleAddTransactionTo).Methods("POST")`

`// account settings`

`secure.HandleFunc("/accounts/settings/default-card/{cardId}", s.handleSetDefaultCard).Methods("POST")`

`secure.HandleFunc("/accounts/settings/change-password/{id}", s.handleUpdatePassword).Methods("PUT")`
## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
