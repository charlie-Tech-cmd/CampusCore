# CampusCore

CampusCore is a modular School Management System (SMS) backend built with Go for higher education institutions. It follows a layered architecture with clear separation of concerns, making the codebase scalable, maintainable, and easy to extend.

 **Status:** Under active development. Core authentication, academic, and infrastructure modules are implemented, with additional administrative services and frontend applications planned.

---

## Features

### Authentication & Authorization

- JWT authentication
- Refresh token support
- Role-Based Access Control (RBAC)
- User registration and profile management

### Academic Management

- Faculty management
- Department management
- Course management
- Course registration
- Enrollment management
- Academic result processing
- Transcript generation
- Timetable management
- Attendance management

### Infrastructure

- PostgreSQL integration
- Database migrations with `golang-migrate`
- Docker & Docker Compose support
- Environment-based configuration
- GitHub Actions continuous integration

### In Progess

- Payment and financial management
- Clearance workflow
- Library management
- Hostel management
- Real-time notifications
- Student, Lecturer, and Admin frontend applications

---

## Technology Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.22+ |
| Database | PostgreSQL |
| API | REST |
| Authentication | JWT |
| Database Migration | golang-migrate |
| Containerization | Docker & Docker Compose |
| CI | GitHub Actions |

---

## Architecture

CampusCore follows a layered architecture where each layer has a single responsibility.

```text
                Client
                   │
                   ▼
              REST API
                   │
                   ▼
            HTTP Handlers
                   │
                   ▼
         Business Services
                   │
                   ▼
      Repository Interfaces
                   │
                   ▼
             PostgreSQL
```

This structure keeps business logic independent from data access and delivery, making the application easier to test, maintain, and extend.

---

## Project Structure

```text
CampusCore/
├── cmd/
│   └── server/
├── database/
├── docker/
├── internal/
│   ├── academic/
│   ├── api/
│   ├── auth/
│   ├── config/
│   ├── governance/
│   ├── middleware/
│   ├── models/
│   ├── notification/
│   ├── repository/
│   ├── services/
│   └── websocket/
├── .github/
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

- Go 1.22 or later
- PostgreSQL
- Docker (optional)

### Clone the Repository

```bash
git clone https://github.com/charlie-Tech-cmd/CampusCore.git
cd CampusCore
go mod download
```

### Configure Environment

Create a `.env` file in the project root.

```env
PORT=8080
DB_SOURCE=postgresql://username:password@localhost:5432/campuscore?sslmode=disable
JWT_SECRET=your_super_secret_key
```

### Run the Application

```bash
go run ./cmd/server
```

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

---

## Implemented Modules

- ✅ Authentication
- ✅ Authorization (RBAC)
- ✅ Student Profile Management
- ✅ Faculty Management
- ✅ Department Management
- ✅ Course Management
- ✅ Enrollment
- ✅ Course Registration
- ✅ Academic Results
- ✅ Transcript Generation
- ✅ Timetable Management
- ✅ Attendance Management
- ✅ Database Migrations
- ✅ Docker Support
- ✅ GitHub Actions CI

---

## Roadmap

- [ ] Payment & Financial Management
- [ ] Clearance Management
- [ ] Library Management
- [ ] Hostel Management
- [ ] Real-time Notifications
- [ ] Swagger / OpenAPI Documentation
- [ ] Continuous Deployment (CD)
- [ ] Student Portal
- [ ] Lecturer Portal
- [ ] Admin Dashboard

---

## Development Status

CampusCore is actively evolving into a comprehensive academic management platform.

Recent milestones include:

- JWT authentication with refresh token support
- Complete academic management foundation
- Transcript generation engine
- PostgreSQL repository layer
- Dockerized development environment
- Automated GitHub Actions CI pipeline

Future development will focus on administrative modules, frontend applications, and deployment.

---

## Contributing

Contributions, suggestions, and feedback are welcome. Feel free to fork the repository, open an issue, or submit a pull request.

---

## License

This project is currently developed for learning, portfolio, and educational purposes. A project license will be added before the first stable release.