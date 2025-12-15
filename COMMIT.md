# 📝 Commit Guidelines

> **Bengi Investment System** ใช้ **Conventional Commits** specification เพื่อให้ commit history อ่านง่าย, สร้าง changelog อัตโนมัติได้, และรองรับ semantic versioning

---

## 📋 Table of Contents

- [Commit Message Format](#-commit-message-format)
- [Types](#-types)
- [Scopes](#-scopes)
- [Emojis](#-emojis)
- [Examples](#-examples)
- [Best Practices](#-best-practices)
- [Branch Naming](#-branch-naming)
- [Pull Request Guidelines](#-pull-request-guidelines)

---

## 📐 Commit Message Format

```
<type>(<scope>): <emoji> <subject>

[optional body]

[optional footer(s)]
```

### Structure Breakdown

| Part | Required | Description |
|------|----------|-------------|
| `type` | ✅ | ประเภทของการเปลี่ยนแปลง (feat, fix, docs, etc.) |
| `scope` | ❌ | ส่วนของ codebase ที่เปลี่ยนแปลง (auth, user, config) |
| `emoji` | ✅ | Visual indicator สำหรับ type |
| `subject` | ✅ | คำอธิบายสั้นๆ (50 ตัวอักษรหรือน้อยกว่า) |
| `body` | ❌ | คำอธิบายเพิ่มเติม (72 ตัวอักษรต่อบรรทัด) |
| `footer` | ❌ | Breaking changes, issue references |

### Rules

1. **Subject**
   - ใช้ **imperative mood** (Add, Fix, Update ไม่ใช่ Added, Fixed, Updated)
   - **ไม่ขึ้นต้นด้วยตัวพิมพ์ใหญ่** หลัง emoji
   - **ไม่มีจุด** ท้ายประโยค
   - จำกัด **50 ตัวอักษร**

2. **Body** (ถ้ามี)
   - แยกจาก subject ด้วย **บรรทัดว่าง**
   - อธิบาย **ทำไม** ไม่ใช่ **อย่างไร**
   - จำกัด **72 ตัวอักษร** ต่อบรรทัด

3. **Footer** (ถ้ามี)
   - `BREAKING CHANGE:` สำหรับ breaking changes
   - `Closes #123` สำหรับ link issues
   - `Refs #456` สำหรับ related issues

---

## 🏷️ Types

| Type | Description | When to Use | Bumps |
|------|-------------|-------------|-------|
| `feat` | New feature | เพิ่ม feature ใหม่ที่ user เห็น | MINOR |
| `fix` | Bug fix | แก้ไข bug | PATCH |
| `docs` | Documentation | เปลี่ยนแปลง documentation เท่านั้น | - |
| `style` | Code style | Formatting, semicolons, whitespace (ไม่กระทบ logic) | - |
| `refactor` | Refactoring | ปรับโครงสร้าง code โดยไม่เปลี่ยน behavior | - |
| `perf` | Performance | ปรับปรุง performance | PATCH |
| `test` | Tests | เพิ่ม/แก้ไข tests | - |
| `build` | Build system | เปลี่ยนแปลง build scripts, dependencies | - |
| `ci` | CI/CD | เปลี่ยนแปลง CI configuration | - |
| `chore` | Maintenance | งานบำรุงรักษาทั่วไป | - |
| `revert` | Revert | ย้อนกลับ commit ก่อนหน้า | - |

### Type Decision Tree

```
การเปลี่ยนแปลงนี้...
│
├─ เพิ่ม feature ใหม่? → feat
├─ แก้ bug? → fix
├─ เปลี่ยน docs เท่านั้น? → docs
├─ เปลี่ยน formatting/style? → style
├─ ปรับโครงสร้าง code? → refactor
├─ เพิ่มความเร็ว? → perf
├─ เพิ่ม/แก้ tests? → test
├─ เปลี่ยน build/deps? → build
├─ เปลี่ยน CI/CD? → ci
├─ ย้อนกลับ commit? → revert
└─ อื่นๆ (cleanup, typo)? → chore
```

---

## 🎯 Scopes

Scopes สำหรับ **Bengi Investment System**:

### Core Modules
| Scope | Description | Example |
|-------|-------------|---------|
| `auth` | Authentication & Authorization | `feat(auth): ✨ add JWT refresh token` |
| `user` | User management | `fix(user): 🐛 fix profile update` |
| `account` | Cash accounts & transactions | `feat(account): ✨ add deposit endpoint` |
| `portfolio` | Portfolio management | `refactor(portfolio): ♻️ extract position logic` |
| `order` | Order management | `feat(order): ✨ add limit order support` |
| `trade` | Trade execution | `fix(trade): 🐛 fix matching engine` |
| `instrument` | Market data | `perf(instrument): ⚡️ cache price data` |

### Infrastructure
| Scope | Description | Example |
|-------|-------------|---------|
| `config` | Configuration | `chore(config): 🔧 update env variables` |
| `db` | Database | `perf(db): ⚡️ add MongoDB indexes` |
| `redis` | Redis cache | `feat(redis): ✨ add Pub/Sub support` |
| `kafka` | Kafka messaging | `fix(kafka): 🐛 fix consumer offset` |
| `ws` | WebSocket | `feat(ws): ✨ add real-time prices` |

### Shared
| Scope | Description | Example |
|-------|-------------|---------|
| `middleware` | HTTP middlewares | `fix(middleware): 🐛 fix CORS headers` |
| `common` | Shared utilities | `refactor(common): ♻️ improve error types` |
| `dto` | Data transfer objects | `feat(dto): ✨ add pagination response` |
| `api` | API general | `docs(api): 📝 update OpenAPI spec` |

### No Scope
สำหรับ changes ที่กระทบหลายส่วนหรือ project-wide:
```
docs: 📝 update README
style: 🎨 run go fmt on all files
chore: 🔧 update .gitignore
```

---

## 😀 Emojis

### Required Emojis (by Type)

| Type | Emoji | Code |
|------|-------|------|
| `feat` | ✨ | `:sparkles:` |
| `fix` | 🐛 | `:bug:` |
| `docs` | 📝 | `:memo:` |
| `style` | 🎨 | `:art:` |
| `refactor` | ♻️ | `:recycle:` |
| `perf` | ⚡️ | `:zap:` |
| `test` | 🧪 | `:test_tube:` |
| `build` | 📦 | `:package:` |
| `ci` | 👷 | `:construction_worker:` |
| `chore` | 🔧 | `:wrench:` |
| `revert` | ⏪ | `:rewind:` |

### Additional Emojis (Context-specific)

| Context | Emoji | When to Use |
|---------|-------|-------------|
| 🔒 | Security fix | `fix(auth): 🔒 fix SQL injection` |
| 🔥 | Remove code/files | `refactor: 🔥 remove deprecated API` |
| 💥 | Breaking change | `feat(api): 💥 change response format` |
| 🚀 | Deploy | `chore: 🚀 prepare v1.0.0 release` |
| 🗃️ | Database | `feat(db): 🗃️ add migration script` |
| 🐳 | Docker | `chore: 🐳 update Dockerfile` |
| ➕ | Add dependency | `build: ➕ add fiber v2` |
| ➖ | Remove dependency | `build: ➖ remove unused package` |
| 🔀 | Merge | `chore: � merge develop into main` |
| 🚧 | WIP | `feat(order): 🚧 work in progress` |

---

## 📚 Examples

### Simple Commits

```bash
# Feature
feat(auth): ✨ add user registration endpoint

# Bug fix
fix(order): 🐛 fix order validation error

# Documentation
docs(api): 📝 add authentication docs

# Refactoring
refactor(trade): ♻️ extract matching logic to service

# Performance
perf(db): ⚡️ add index on users.email

# Tests
test(auth): 🧪 add login unit tests

# Style
style: 🎨 apply gofmt to all files

# Chore
chore: 🔧 update .env.example
```

### Commits with Body

```bash
feat(portfolio): ✨ add position tracking

Implement real-time position tracking with:
- FIFO lot-based cost calculation
- Unrealized P&L computation
- Automatic position updates on trade

Part of the core trading feature set.
```

### Commits with Footer

```bash
fix(auth): 🐛 fix token expiration check

The token was being validated against local time
instead of UTC, causing premature expiration for
users in different timezones.

Closes #42
```

### Breaking Changes

```bash
feat(api): 💥 change response format

BREAKING CHANGE: API responses now use a standardized format:
{
  "success": true,
  "data": {...},
  "message": "..."
}

Previous format was just the raw data object.
Migration: Update all client code to access response.data

Refs #100
```

### Revert Commit

```bash
revert: ⏪ feat(order): add market order support

This reverts commit a1b2c3d4e5f6.

Reason: Market orders causing unexpected behavior
in paper trading mode. Will re-implement with
proper safeguards.
```

---

## ✅ Best Practices

### DO ✅

```bash
# Specific and descriptive
feat(auth): ✨ add password reset with email verification

# Clear scope
fix(order): 🐛 fix quantity validation for fractional shares

# Explain why in body
perf(db): ⚡️ add compound index on trades collection

Adding index on (portfolioId, createdAt) to speed up
trade history queries from 500ms to 20ms.
```

### DON'T ❌

```bash
# Too vague
fix: 🐛 fix bug

# No emoji
feat(auth): add login

# Past tense
feat(user): ✨ added registration

# With period
fix(order): 🐛 fix validation.

# Too long subject
feat(portfolio): ✨ add comprehensive portfolio management with positions, lots, and real-time P&L calculation
```

### Commit Frequency

| Situation | Recommendation |
|-----------|----------------|
| New feature | Commit เมื่อ feature ทำงานได้ (ไม่ต้องรอ perfect) |
| Bug fix | 1 commit ต่อ 1 bug |
| Refactoring | Commit เมื่อ tests ยังผ่าน |
| WIP | ใช้ `🚧` แล้วค่อย squash ทีหลัง |

### Atomic Commits

แต่ละ commit ควร:
1. ✅ ทำเรื่องเดียว
2. ✅ ไม่ทำให้ build พัง
3. ✅ Tests ยังผ่าน
4. ✅ Revert ได้โดยไม่กระทบส่วนอื่น

---

## 🌿 Branch Naming

### Format

```
<type>/<short-description>
```

### Examples

| Type | Branch Name | Description |
|------|-------------|-------------|
| Feature | `feat/user-registration` | เพิ่ม feature ใหม่ |
| Bug fix | `fix/order-validation` | แก้ bug |
| Hotfix | `hotfix/auth-crash` | แก้ bug เร่งด่วนใน production |
| Refactor | `refactor/trade-service` | Refactoring |
| Docs | `docs/api-guide` | Documentation |
| Chore | `chore/update-deps` | Maintenance |

### Branch Flow

```
main (production)
  │
  └── develop (staging)
        │
        ├── feat/user-auth
        ├── feat/order-management
        ├── fix/login-bug
        └── refactor/database-layer
```

---

## 🔀 Pull Request Guidelines

### PR Title Format

```
<type>(<scope>): <emoji> <description>
```

เหมือน commit message แต่อธิบายภาพรวมของ PR

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] ✨ New feature
- [ ] 🐛 Bug fix
- [ ] 📝 Documentation
- [ ] ♻️ Refactoring
- [ ] ⚡️ Performance
- [ ] 🧪 Tests

## Related Issues
Closes #123

## Checklist
- [ ] Code follows project style
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented)

## Screenshots (if applicable)
```

---

## 🔧 Git Hooks (Recommended)

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run go fmt
go fmt ./...

# Run tests
go test ./... -short

# Run linter
golangci-lint run
```

### Commit-msg Hook

```bash
#!/bin/sh
# .git/hooks/commit-msg

commit_msg=$(cat "$1")

# Check commit format
if ! echo "$commit_msg" | grep -qE "^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .+ .+"; then
    echo "❌ Invalid commit message format!"
    echo "Expected: <type>(<scope>): <emoji> <subject>"
    echo "Example: feat(auth): ✨ add login endpoint"
    exit 1
fi

echo "✅ Commit message is valid"
```

---

## 📌 Quick Reference Card

```
┌─────────────────────────────────────────────────────────┐
│                    COMMIT CHEATSHEET                    │
├─────────────────────────────────────────────────────────┤
│ feat(scope): ✨ add new feature                         │
│ fix(scope): 🐛 fix bug                                  │
│ docs(scope): 📝 update docs                             │
│ style(scope): 🎨 format code                            │
│ refactor(scope): ♻️ refactor code                       │
│ perf(scope): ⚡️ improve performance                     │
│ test(scope): 🧪 add tests                               │
│ build(scope): 📦 change build                           │
│ ci(scope): 👷 update CI                                 │
│ chore(scope): 🔧 maintenance                            │
│ revert: ⏪ revert commit                                │
├─────────────────────────────────────────────────────────┤
│ 💥 Breaking | 🔒 Security | 🔥 Remove | 🚧 WIP          │
└─────────────────────────────────────────────────────────┘
```

---

<div align="center">

**Follow these guidelines to maintain a clean and readable Git history! 🚀**

</div>
