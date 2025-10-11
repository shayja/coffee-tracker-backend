# Coffee Tracker Backend Architecture

```text
           ┌──────────────────────────────┐
           │         HTTP Layer           │
           │ ┌─────────────────────────┐ │
           │ │  Handlers (controllers) │ │
           │ └─────────────────────────┘ │
           │ ┌─────────────────────────┐ │
           │ │ Middleware (auth, logging, │
           │ │ cors, user, etc.)          │
           │ └─────────────────────────┘ │
           │ ┌─────────────────────────┐ │
           │ │ DTOs / Request/Response │ │
           │ └─────────────────────────┘ │
           └─────────────▲──────────────┘
                         │
                         │ calls
                         ▼
           ┌──────────────────────────────┐
           │        Use Cases Layer       │
           │ ┌─────────────────────────┐ │
           │ │  createCoffeeEntry       │ │
           │ │  getCoffeeEntries        │ │
           │ │  generateOTP             │ │
           │ │  getUserSettings         │ │
           │ └─────────────────────────┘ │
           └─────────────▲──────────────┘
                         │
                         │ uses
                         ▼
           ┌──────────────────────────────┐
           │      Repository Layer        │
           │ ┌─────────────────────────┐ │
           │ │ Postgres / Supabase DB   │ │
           │ │ Implementations          │ │
           │ │  - CoffeeEntryRepoImpl   │ │
           │ │  - UserRepoImpl          │ │
           │ │  - KVRepoImpl            │ │
           │ └─────────────────────────┘ │
           └─────────────▲──────────────┘
                         │
                         │ accesses / persists
                         ▼
           ┌──────────────────────────────┐
           │ Infrastructure / Services    │
           │ ┌─────────────────────────┐ │
           │ │ JWTService / Auth         │ │
           │ │ SMS Service / Twilio      │ │
           │ │ Storage / Supabase        │ │
           │ └─────────────────────────┘ │
           └──────────────────────────────┘
```
