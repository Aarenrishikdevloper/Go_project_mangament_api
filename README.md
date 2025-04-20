# 🗂️ Project Management API

A RESTful Project Management API built with [Go](https://golang.org/), [Fiber](https://gofiber.io/), [MongoDB](https://www.mongodb.com/), and [JWT](https://jwt.io/) for secure authentication. This API provides features for user authentication, project and task management.

## 🚀 Features

- User Registration & Login (JWT Auth)
- Create, Read, Update, Delete (CRUD) for:
  - Projects
  - Tasks
- Role-based access control (optional)
- MongoDB integration for data persistence
- Fast performance with Fiber framework

---

## 📦 Tech Stack

- **Backend**: Go, Fiber
- **Database**: MongoDB
- **Authentication**: JWT
- **Others**: Go Modules, Fiber Middleware, .env for config

---


---

## ⚙️ Getting Started

### Prerequisites

- Go 1.18+
- MongoDB
- Git

### Installation

```bash
# Clone the repo
git clone https://github.com/yourusername/project-management-api.git
cd project-management-api

# Install dependencies
go mod tidy

# Run the app
go run main.go    

