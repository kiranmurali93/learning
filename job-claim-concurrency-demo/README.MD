# Job Claim Concurrency Demo

A demonstration of handling concurrency in a job claiming system using Go and PostgreSQL.

## Overview
This project illustrates how to properly handle concurrent requests in a web application where multiple users (recruiters) might try to claim a job simultaneously, with a business rule that limits claims to 4 recruiters per job.

## Project Structure
```
|- db/
   |- db.go             # Database connection setup
|- handlers/
   |- jobs.go           # API handlers for job claiming
|-.env                  # Environment variables (not included)
|- main.go              # Server initialization
|- test.go              # Concurrency test simulator
```

## Database Setup
The project uses PostgreSQL. Create the following tables:

`companies` - Organizations that post jobs
`jobs` - Job listings
`recruiter` - Users who can claim jobs
`job_claims` - Many-to-many relationship tracking job claims
`job_claims_count` - Tracks number of claims per job (for optimized implementation)

## Running the Demo
1. Set up your PostgreSQL database and update the .env file:
``` 
DB_CONNECTION_STRING=postgres://username:password@localhost:5432/dbname
```

2. Start the server:
```
go run main.go
```

3. Run the concurrency test:
```
go run test.go
```