# Commit Guide - Organized by Module

Use this guide to commit changes separately by module.

---

## Backend Commits

### 1. Backend: Demo Account Module
```bash
git add backend/module/account/controller/demo.controller.go \
        backend/module/account/dto/demo.dto.go \
        backend/module/account/routes/demo.route.go \
        backend/module/account/service/demo.service.go \
        backend/module/account/model/account.go \
        backend/module/account/repository/account.repository.go

git commit -m "feat(backend): add demo account module with deposit and reset"
```

### 2. Backend: Instrument & Market Data
```bash
git add backend/module/instrument/controller/instrument.controller.go \
        backend/module/instrument/dto/instrument.dto.go \
        backend/module/instrument/routes/instrument.route.go \
        backend/module/instrument/service/instrument.service.go \
        backend/module/instrument/service/marketdata.service.go

git commit -m "feat(backend): add Yahoo Finance API for historical charts and URL decode for crypto symbols"
```

### 3. Backend: Instrument Seeder
```bash
git add backend/pkg/seeder/instrument.seeder.go \
        backend/pkg/seeder/role.seeder.go

git commit -m "feat(backend): add 130+ instruments seeder (stocks, ETFs, crypto)"
```

### 4. Backend: Auth Controller
```bash
git add backend/module/auth/controller/auth.controller.go

git commit -m "fix(backend): return tokens in response body for mobile app"
```

### 5. Backend: Trade Models
```bash
git add backend/module/trade/model/leverage_position.go \
        backend/module/trade/model/option.go

git commit -m "feat(backend): add leverage position and option models"
```

### 6. Backend: Main Entry
```bash
git add backend/main.go

git commit -m "chore(backend): register demo routes in main"
```

---

## Frontend Mobile Commits

### 7. Mobile: Theme & Colors
```bash
git add frontend/mobile/constants/theme.ts

git commit -m "style(mobile): improve gray text visibility in dark theme"
```

### 8. Mobile: Market Screen
```bash
git add frontend/mobile/app/(tabs)/market/index.tsx \
        frontend/mobile/app/(tabs)/market/[symbol].tsx

git commit -m "feat(mobile): redesign market screen with FlatList, compact cards, and TradingView chart"
```

### 9. Mobile: Portfolio Screen (Dime-style)
```bash
git add frontend/mobile/app/(tabs)/portfolio/index.tsx \
        frontend/mobile/components/portfolio/PortfolioCard.tsx \
        frontend/mobile/components/portfolio/HoldingItem.tsx \
        frontend/mobile/components/portfolio/OptionItem.tsx

git commit -m "feat(mobile): add Dime-style portfolio with gradient card, holdings, and options"
```

### 10. Mobile: Chart Component
```bash
git add frontend/mobile/components/chart/

git commit -m "feat(mobile): add TradingView Lightweight Charts component"
```

### 11. Mobile: Demo Store & API
```bash
git add frontend/mobile/stores/demo.ts \
        frontend/mobile/types/demo.ts \
        frontend/mobile/services/api.ts

git commit -m "feat(mobile): add demo account store and API service"
```

### 12. Mobile: Market Types & Store
```bash
git add frontend/mobile/types/market.ts \
        frontend/mobile/stores/market.ts

git commit -m "feat(mobile): add CandleData type and fix market store"
```

### 13. Mobile: Other Screens
```bash
git add frontend/mobile/app/(tabs)/index.tsx \
        frontend/mobile/app/(tabs)/profile/index.tsx \
        frontend/mobile/app/(tabs)/trade/index.tsx

git commit -m "feat(mobile): integrate demo balance in Home, Profile, and Trade screens"
```

---

## Quick All-in-One (If preferred)

```bash
# Backend
git add backend/
git commit -m "feat(backend): demo account, Yahoo Finance charts, 130+ instruments"

# Frontend
git add frontend/mobile/
git commit -m "feat(mobile): Dime-style portfolio, TradingView charts, demo account integration"
```

---

## Push
```bash
git push origin main
```
