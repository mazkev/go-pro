# ⚡ Go-Pro: Enterprise REST API & Kanban Board Backend

Production-grade RESTful API backend built with **Golang**, **Clean Architecture**, **GORM**, and **Gin Web Framework**. Designed to seamlessly power both E-Commerce services and the interactive [💎 Glass Elite Trello (react-trello-kanban-suite)](https://github.com/mazkev/react-trello-kanban-suite) frontend suite.

---

## 🌟 Key Architectural Features

- **🏛️ 4-Layer Clean Architecture**: Strict separation of concerns (Domain, Repository, Usecase, Delivery/HTTP) with Interface-based Dependency Injection.
- **📋 Relational Kanban Engine**: Hierarchical multi-board architecture (`User` ➔ `Board` ➔ `List` ➔ `Card` ➔ `ChecklistItem`) featuring GORM nested preloading and drag-and-drop column reordering.
- **🔒 Enterprise Security**: JWT (`golang-jwt/jwt/v5`) authentication middleware and cryptographic password hashing via `golang.org/x/crypto/bcrypt`.
- **⚡ Asynchronous Concurrency**: Non-blocking background worker tasks (Goroutines) for operations like automated welcome emails and notifications.
- **🛡️ Data Integrity & Soft Deletes**: Zero data loss design using `gorm.DeletedAt` timestamps and database cascading relations.
- **🧪 Automated Table-Driven Testing**: High-coverage unit tests with interface mocking, running in sub-second execution times.
- **🚢 Cloud-Native DevOps**: Multi-Stage `Dockerfile` generating an ultra-lightweight **~15MB container** with automated GitHub Actions CI/CD pipelines.
- **🌐 Cross-Origin Enabled**: Integrated CORS middleware supporting seamless client connections (React, Next.js, Vite).

---

## 🛠️ Tech Stack

| Component | Technology |
| :--- | :--- |
| **Language** | Go (Golang) v1.23+ |
| **HTTP Framework** | [Gin Gonic](https://github.com/gin-gonic/gin) |
| **ORM / Data Access** | [GORM](https://gorm.io) |
| **Database** | SQLite (Pure Go via [glebarez/sqlite](https://github.com/glebarez/sqlite)) |
| **Authentication** | JWT (`golang-jwt/jwt/v5`) + Bcrypt |
| **Containerization** | Docker Multi-Stage Build (Alpine Runner) |
| **CI/CD** | GitHub Actions Pipeline |
| **Companion Frontend** | [react-trello-kanban-suite](https://github.com/mazkev/react-trello-kanban-suite) (React 19 + Tailwind v4 + Zustand) |

---

## 📡 API Endpoints Reference

### 1. Authentication & Profile (`/`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/register` | Public | Register new user (triggers async welcome email) |
| `POST` | `/login` | Public | Authenticate user & issue signed JWT Bearer token |
| `GET` | `/me` | 🔒 Private | Verify active user session and retrieve profile info |

### 2. 📋 Kanban Board (`/boards`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/boards` | 🔒 Private | Retrieve all boards owned by the authenticated user |
| `GET` | `/boards/:id` | 🔒 Private | Retrieve single board with nested lists, cards & checklists |
| `POST` | `/boards` | 🔒 Private | Create new board (auto-generates "To Do", "In Progress", "Done") |
| `PUT` | `/boards/:id` | 🔒 Private | Update board title and color theme |
| `DELETE` | `/boards/:id` | 🔒 Private | Delete board and cascade-delete all child lists and cards |

### 3. 📑 Kanban Lists / Columns (`/lists`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/boards/:boardId/lists` | 🔒 Private | Add new column into a specific board |
| `PUT` | `/lists/:id` | 🔒 Private | Rename column title |
| `PUT` | `/lists/:id/move` | 🔒 Private | **Reorder**: Change column position index |
| `DELETE` | `/lists/:id` | 🔒 Private | Delete column and cascade-delete all task cards inside |

### 4. 🗂️ Kanban Cards & Drag-and-Drop (`/cards`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/lists/:listId/cards` | 🔒 Private | Add new task card into a specific column list |
| `PUT` | `/cards/:id` | 🔒 Private | Update card title, description, label, due date, or status |
| `PUT` | `/cards/:id/move` | 🔒 Private | **Drag & Drop**: Move card across columns or reorder position |
| `DELETE` | `/cards/:id` | 🔒 Private | Delete a task card |

### 5. ✅ Card Checklists / Sub-Tasks (`/checklist`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/cards/:cardId/checklist` | 🔒 Private | Add new sub-task item to a card |
| `PUT` | `/checklist/:id` | 🔒 Private | Toggle sub-task status (completed / pending) or edit title |
| `DELETE` | `/checklist/:id` | 🔒 Private | Delete a sub-task item |

### 6. 📦 Product Catalog (`/products`)
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/products` | Public | List active products |
| `GET` | `/products/:id` | Public | Get single product |
| `GET` | `/products/search?keyword=` | Public | Search products by name |
| `POST` | `/products` | 🔒 Private | Add new product |
| `PUT` | `/products/:id` | 🔒 Private | Update product details |
| `DELETE` | `/products/:id` | 🔒 Private | **Soft Delete**: Archive product without physical deletion |

---

## 📁 Directory Structure (Clean Architecture)

```
├── .github/
│   └── workflows/
│       └── ci.yml                      # Automated GitHub Actions test & build pipeline
├── cmd/
│   └── api/
│       └── main.go                     # Application entry point & Dependency Injection wiring
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── middleware/
│   │       │   ├── auth_middleware.go  # JWT Bearer Token validation guard
│   │       │   └── cors_middleware.go  # Cross-Origin Resource Sharing handler
│   │       ├── board_handler.go        # Kanban Board HTTP controller
│   │       ├── card_handler.go         # Kanban Card & Drag-and-Drop controller
│   │       ├── checklist_handler.go    # Task Checklist controller
│   │       ├── list_handler.go         # Kanban List/Column controller
│   │       ├── product_handler.go      # Product HTTP controller
│   │       └── user_handler.go         # Auth & Profile HTTP controller
│   ├── domain/
│   │   ├── board.go                    # Kanban models (Board, List, Card, Checklist) & contracts
│   │   ├── product.go                  # Product model with Soft Delete & contracts
│   │   └── user.go                     # User model & contracts
│   ├── repository/
│   │   ├── board_repository.go         # GORM nested preload repository
│   │   ├── card_repository.go          # Card CRUD & atomic move repository
│   │   ├── checklist_repository.go     # Checklist item repository
│   │   ├── list_repository.go          # List CRUD & column reordering repository
│   │   ├── product_repository.go       # Product SQLite repository
│   │   └── user_repository.go          # User database access
│   ├── usecase/
│   │   ├── board_usecase.go            # Board business logic (auto-default lists)
│   │   ├── card_usecase.go             # Card business logic
│   │   ├── checklist_usecase.go        # Checklist item business logic
│   │   ├── list_usecase.go             # List/Column business logic
│   │   ├── product_usecase.go          # Product business logic & validations
│   │   ├── product_usecase_test.go     # Table-Driven Unit Tests with Mock Repository
│   │   └── user_usecase.go             # Auth logic & async Goroutine email triggers
│   └── utils/
│       └── token.go                    # JWT token generation & parsing utilities
├── Dockerfile                          # Multi-Stage Build (~15MB image)
├── .dockerignore
├── go.mod
├── go.sum
├── README.md
└── rest.http                           # Interactive API test suite
```

---

## 🏃 Getting Started

### Local Development
1. Clone this repository:
   ```bash
   git clone https://github.com/mazkev/go-pro.git
   cd go-pro
   ```
2. Run the application:
   ```bash
   go run ./cmd/api
   ```
3. The server starts at `http://localhost:8080`.
4. Open `rest.http` to test all endpoints interactively.

### Running Automated Tests
```bash
go test -v ./...
```

### Docker Multi-Stage Build
```bash
# Build lightweight container image
docker build -t go-pro-api .

# Run container
docker run -p 8080:8080 go-pro-api
```