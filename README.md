# Zecas
Zecas is a general tool for running scans that gathers data about targets, different scanners can be used to gather different information such as findings, port openings, new targets.

The goal with the project is to make a platform that enables easy management and accessibility to security information regarding multiple targets.

## Getting started
To start off running this project.

Add relevant values to `.env` files in frontend and backend directories.

### General
```bash
cd backend
docker compose up -d
```

### Frontend

Move to frontend directory
```
cd frontend
```

Install dependencies
```bash
npm i --force
```

Run astro dev
```bash
npm run dev
```

### Backend

Run backend go file.
```bash
cd backend
go run /cmd/api/main.go
```

### Worker
Run worker go file.
```bash
cd backend
go run /cmd/worker/main.go
```

## Documentation
Coming soon...

## Contribute
Coming soon...