<div align="center">

# 🚀 Bengi Investment System

**A Modern Stock Trading & Investment Platform**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=for-the-badge)](https://gofiber.io/)
[![MongoDB](https://img.shields.io/badge/MongoDB-6.0+-47A248?style=for-the-badge&logo=mongodb)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

*Inspired by [Webull](https://www.webull.com/) & [Dime](https://dime.co.th/) — Built for Performance & Simplicity*

</div>

---

## 📖 Overview

**Bengi Investment System** is a full-featured stock trading and investment platform built with a **Modular Monolith Architecture**. Designed to handle high-frequency trading operations while maintaining clean, maintainable code structure.

This project takes inspiration from leading trading platforms like **Webull** and **Dime**, implementing their best features while addressing common architectural pain points.

---

## ✨ Features

### Phase 1 — Core Trading (Current)
| Feature | Description |
|---------|-------------|
| 🔐 **Authentication** | JWT-based auth with role-based access control (RBAC) |
| 💰 **Account Management** | Multi-currency accounts with balance tracking |
| 📊 **Portfolio Management** | Multiple portfolios per user |
| 📈 **Order Management** | Market & Limit orders with GTC/GTD support |
| 💹 **Trade Execution** | Real-time trade matching and execution |
| 📦 **Position Tracking** | FIFO/LIFO lot-based cost calculation |
| 🏦 **Transaction History** | Complete audit trail for all cash movements |

### Phase 2 — Enhanced Features (Planned)
| Feature | Description |
|---------|-------------|
| 👀 **Watchlists** | Custom watchlists for tracking instruments |
| 🔔 **Price Alerts** | Real-time price notifications |
| 📱 **Notifications** | Push notifications for order fills, dividends |
| 💵 **Dividends** | Automatic dividend tracking & distribution |
| 📋 **Audit Logs** | Comprehensive logging for compliance |

---

## 🏗️ Architecture

### Why Modular Monolith?

We chose **Modular Monolith** over Microservices for these reasons:

| Aspect | Microservices | Modular Monolith ✅ |
|--------|---------------|---------------------|
| **Complexity** | High (network, deployment) | Low (single binary) |
| **Development Speed** | Slower | Faster iteration |
| **Debugging** | Distributed tracing needed | Simple stack traces |
| **Deployment** | Multiple services | Single deployment |
| **Scalability** | Scale individual services | Scale entire app |
| **Team Size** | Large teams | Small to medium teams |

> 💡 **Key Insight**: Start as Modular Monolith, evolve to Microservices when scale demands it.

### Project Structure

```
bengi-investment-system/
├── backend/
│   ├── module/                    # Feature Modules
│   │   ├── auth/                  # Authentication & Users
│   │   │   ├── controller/        # HTTP handlers
│   │   │   ├── dto/               # Request/Response DTOs
│   │   │   ├── model/             # Domain models
│   │   │   ├── service/           # Business logic
│   │   │   └── module.go          # Module registration
│   │   ├── account/               # Cash & Transactions
│   │   ├── portfolio/             # Portfolios & Positions
│   │   ├── instrument/            # Market Data
│   │   ├── order/                 # Order Management
│   │   └── trade/                 # Trade Execution
│   │
│   ├── pkg/                       # Shared Packages
│   │   ├── common/                # Errors, Response, Pagination
│   │   ├── config/                # Configuration
│   │   ├── core/                  # Infrastructure
│   │   │   ├── database/          # MongoDB connection
│   │   │   ├── kafka/             # Event streaming
│   │   │   └── redis/             # Caching
│   │   ├── middleware/            # Auth, RBAC, Rate Limit
│   │   └── utils/                 # JWT, Hash, Validator
│   │
│   ├── main.go                    # Entry point
│   ├── go.mod
│   └── go.sum
│
└── README.md
```

### Design Patterns

| Pattern | Usage |
|---------|-------|
| **Repository Pattern** | Data access abstraction |
| **Service Layer** | Business logic encapsulation |
| **DTO Pattern** | Request/Response transformation |
| **Dependency Injection** | Loose coupling between modules |
| **Event-Driven** | Async processing via Kafka |

---

## 🛠️ Tech Stack

### Backend
| Technology | Purpose |
|------------|---------|
| **Go 1.21+** | Primary language — fast, typed, concurrent |
| **Go Fiber v2** | Web framework — Express-like, high performance |
| **MongoDB** | Database — flexible schema, horizontal scaling |
| **Redis** | Caching — session, rate limiting, real-time data |
| **Apache Kafka** | Event streaming — order events, price feeds |

### External APIs
| Service | Purpose |
|---------|---------|
| **Twelve Data** | Real-time & historical market data |

### DevOps (Planned)
| Technology | Purpose |
|------------|---------|
| **Docker** | Containerization |
| **Docker Compose** | Local development |
| **GitHub Actions** | CI/CD pipeline |

---

## 🔄 Data Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Fiber     │────▶│  Service    │
│   (App)     │◀────│   Router    │◀────│   Layer     │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                    ┌──────────────────────────┼──────────────────────────┐
                    │                          │                          │
                    ▼                          ▼                          ▼
             ┌─────────────┐           ┌─────────────┐           ┌─────────────┐
             │   MongoDB   │           │    Redis    │           │    Kafka    │
             │  (Primary)  │           │   (Cache)   │           │  (Events)   │
             └─────────────┘           └─────────────┘           └─────────────┘
```

---

## 📊 Database Schema

### Core Entities

```
Users ──┬──▶ Accounts ──▶ Transactions
        │
        └──▶ Portfolios ──┬──▶ Orders ──▶ Trades
                          │
                          └──▶ Positions ──▶ PositionLots

Instruments ──▶ InstrumentPrices
```

### Key Tables
- **users** — User authentication & profile
- **accounts** — Cash balances (multi-currency)
- **portfolios** — Investment portfolios
- **instruments** — Stocks, ETFs, Crypto
- **orders** — Pending orders (Market/Limit)
- **trades** — Executed trades
- **positions** — Current holdings
- **positionLots** — Cost basis tracking (FIFO/LIFO)

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- MongoDB 6.0+
- Redis 7.0+
- Kafka 3.0+ (optional for Phase 1)

### Installation

```bash
# Clone the repository
git clone https://github.com/bricksocoolxd/bengi-investment-system.git
cd bengi-investment-system/backend

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the application
go run main.go
```

### Environment Variables

```env
# Server
PORT=3000
ENV=development

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=bengi_investment

# Redis
REDIS_URI=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRES_IN=24h

# Twelve Data API
TWELVE_DATA_API_KEY=your-api-key
```

---

## 📈 Roadmap

### ✅ Phase 1 — Foundation (Current)
- [x] Project structure setup
- [ ] Authentication module (JWT + RBAC)
- [ ] Account management
- [ ] Portfolio CRUD
- [ ] Instrument data sync
- [ ] Order placement
- [ ] Trade execution
- [ ] Position tracking

### 🔜 Phase 2 — Enhanced Features
- [ ] Watchlists
- [ ] Price alerts
- [ ] Push notifications
- [ ] Dividend tracking
- [ ] Advanced order types (Stop, Stop-Limit)

### 🔮 Phase 3 — Scale & Polish
- [ ] Real-time WebSocket feeds
- [ ] Mobile app (React Native / Flutter)
- [ ] Performance optimization
- [ ] Comprehensive testing
- [ ] Production deployment

---

## 🆚 Comparison with Inspirations

| Feature | Webull | Dime | Bengi ✅ |
|---------|--------|------|----------|
| Multi-Portfolio | ✅ | ❌ | ✅ |
| Lot-based P&L | ✅ | ✅ | ✅ (FIFO/LIFO) |
| Real-time Data | ✅ | ✅ | ✅ (Twelve Data) |
| Open Source | ❌ | ❌ | ✅ |
| Self-hosted | ❌ | ❌ | ✅ |
| Customizable | ❌ | ❌ | ✅ |

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ by [@bricksocoolxd](https://github.com/bricksocoolxd)**

⭐ Star this repo if you find it useful!

</div>
