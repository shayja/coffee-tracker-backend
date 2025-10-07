# Coffee Tracker Backend Architecture

```text
                         ┌───────────────────────┐
                         │      HTTP/API         │
                         │ (mux router + DTOs)  │
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │       Middleware      │
                         │ - AuthMiddleware      │
                         │ - UserMiddleware      │
                         │ - CorsMiddleware      │
                         │ - LoggingMiddleware   │
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │       Handlers        │
                         │ - auth_handler.go     │
                         │ - coffee_entry_handler│
                         │ - user_handler.go     │
                         │ - generic_kv_handler  │
                         │   Uses DTOs           │
                         │   (requests/responses)│
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │       Usecases        │
                         │ - Create/Edit/Delete  │
                         │   coffee entries      │
                         │ - Generate/Validate OTP│
                         │ - Save/Get/Delete     │
                         │   Refresh Tokens      │
                         │ - Get/Update Profile  │
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │     Repositories      │
                         │ Interfaces:           │
                         │ - UserRepository      │
                         │ - CoffeeRepo          │
                         │ - AuthRepository      │
                         │ Implementations:     │
                         │ - Postgres/Supabase  │
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │   Database & Storage  │
                         │ - PostgreSQL/Supabase │
                         │ - Supabase Storage    │
                         └─────────┬────────────┘
                                   │
                                   ▼
                         ┌───────────────────────┐
                         │   JWT / Auth Layer    │
                         │ - jwt_service.go      │
                         │ - jwt_utils.go        │
                         │ - Token generation    │
                         │ - Token validation    │
                         └───────────────────────┘
```
