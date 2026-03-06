# youwont — Architecture & Technical Design

---

## 1. System Overview

```
┌─────────────────────┐         ┌─────────────────────┐
│   React Native      │         │   Supabase           │
│   (Expo Router)     │────────▶│   (Auth only)        │
│                     │  JWT    │                      │
│   TanStack Query    │◀────────│   email/password     │
└────────┬────────────┘         └──────────────────────┘
         │
         │  REST + Supabase JWT
         │  in Authorization header
         ▼
┌─────────────────────┐
│   Go API            │
│   (Echo v5)         │
│                     │
│   Verifies JWT      │
│   All business logic│
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│   MongoDB Atlas     │
│                     │
│   users             │
│   groups            │
│   bets              │
│   invites           │
│   notifications     │
└─────────────────────┘
```

**Auth flow:** Supabase handles signup/login and issues a JWT. The frontend sends that JWT on every request to the Go API. The Go API verifies it, extracts the Supabase user ID, looks up the user in MongoDB, and proceeds.

---

## 2. Backend Architecture (Go)

### Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point — wires dependencies, starts server
├── internal/
│   ├── config/
│   │   └── config.go               # Environment variables, app configuration
│   ├── middleware/
│   │   └── auth.go                 # JWT verification, attaches user to context
│   ├── model/
│   │   ├── user.go                 # User domain model + MongoDB tags
│   │   ├── group.go                # Group + embedded Member
│   │   ├── bet.go                  # Bet + embedded Wager
│   │   ├── invite.go               # Invite model
│   │   └── notification.go         # Notification model
│   ├── handler/
│   │   ├── user.go                 # POST /users, GET /users/me, GET /users/search
│   │   ├── group.go                # CRUD groups, join via code
│   │   ├── bet.go                  # CRUD bets, wagers, resolve, change decider
│   │   ├── invite.go               # Send, list, accept, decline
│   │   └── notification.go         # List, unread count, mark read
│   ├── service/
│   │   ├── user.go                 # User business logic
│   │   ├── group.go                # Group business logic
│   │   ├── bet.go                  # Bet + wager + resolution logic
│   │   ├── invite.go               # Invite logic + notification creation
│   │   └── notification.go         # Notification queries
│   └── repository/
│       ├── user.go                 # MongoDB operations for users
│       ├── group.go                # MongoDB operations for groups
│       ├── bet.go                  # MongoDB operations for bets
│       ├── invite.go               # MongoDB operations for invites
│       └── notification.go         # MongoDB operations for notifications
├── Dockerfile
├── go.mod
└── go.sum
```

### Layer Responsibilities

```
Request → Middleware (auth) → Handler → Service → Repository → MongoDB
                                ↓           ↓
                            Validates    Business logic,
                            input,       transactions,
                            binds JSON,  point calculations,
                            returns HTTP notification fan-out
```

- **Middleware**: Verifies JWT, looks up user, attaches to request context. Every route except `POST /users` requires auth.
- **Handler**: HTTP layer only. Parses request body, validates input shape, calls service, returns JSON response. No business logic.
- **Service**: All business logic lives here. Validates business rules (enough points? is user the decider? is bet still open?). Orchestrates transactions across multiple repositories.
- **Repository**: Pure data access. One method = one MongoDB operation. No business logic, no knowledge of HTTP.

### Dependency Injection

Each layer depends on the one below it via **interfaces**. This is idiomatic Go — the consumer defines the interface.

```go
// service/bet.go
type BetRepository interface {
    FindByID(ctx context.Context, id primitive.ObjectID) (*model.Bet, error)
    FindByGroupID(ctx context.Context, groupID primitive.ObjectID, status *string) ([]model.Bet, error)
    Create(ctx context.Context, bet *model.Bet) error
    PushWager(ctx context.Context, betID primitive.ObjectID, wager model.Wager) error
    Resolve(ctx context.Context, betID primitive.ObjectID, winningSide string) error
}

type BetService struct {
    bets    BetRepository
    users   UserRepository
    notifs  NotificationRepository
}
```

The concrete repository structs implement these interfaces. `main.go` wires everything:

```go
func main() {
    cfg := config.Load()
    db := connectMongo(cfg.MongoURI)

    // Repositories
    userRepo := repository.NewUserRepo(db)
    groupRepo := repository.NewGroupRepo(db)
    betRepo := repository.NewBetRepo(db)
    inviteRepo := repository.NewInviteRepo(db)
    notifRepo := repository.NewNotificationRepo(db)

    // Services
    userSvc := service.NewUserService(userRepo)
    groupSvc := service.NewGroupService(groupRepo, userRepo)
    betSvc := service.NewBetService(betRepo, userRepo, notifRepo, db)
    inviteSvc := service.NewInviteService(inviteRepo, groupRepo, notifRepo, db)
    notifSvc := service.NewNotificationService(notifRepo)

    // Handlers
    userHandler := handler.NewUserHandler(userSvc)
    groupHandler := handler.NewGroupHandler(groupSvc)
    betHandler := handler.NewBetHandler(betSvc)
    inviteHandler := handler.NewInviteHandler(inviteSvc)
    notifHandler := handler.NewNotificationHandler(notifSvc)

    // Router
    e := echo.New()
    auth := middleware.NewAuth(cfg.SupabaseJWTSecret, userSvc)

    // Public
    e.POST("/users", userHandler.Create)

    // Protected — all require auth middleware
    api := e.Group("", auth.Required)

    api.GET("/users/me", userHandler.Me)
    api.GET("/users/search", userHandler.Search)

    api.POST("/groups", groupHandler.Create)
    api.GET("/groups", groupHandler.List)
    api.GET("/groups/:id", groupHandler.Get)
    api.POST("/groups/join", groupHandler.JoinByCode)

    api.POST("/groups/:id/invites", inviteHandler.Send)
    api.GET("/invites", inviteHandler.ListMine)
    api.POST("/invites/:id/accept", inviteHandler.Accept)
    api.POST("/invites/:id/decline", inviteHandler.Decline)

    api.GET("/notifications", notifHandler.List)
    api.GET("/notifications/unread-count", notifHandler.UnreadCount)
    api.POST("/notifications/:id/read", notifHandler.MarkRead)
    api.POST("/notifications/read-all", notifHandler.MarkAllRead)

    api.GET("/groups/:id/bets", betHandler.ListByGroup)
    api.POST("/groups/:id/bets", betHandler.Create)
    api.GET("/bets/:id", betHandler.Get)
    api.POST("/bets/:id/wagers", betHandler.PlaceWager)
    api.PUT("/bets/:id/decider", betHandler.ChangeDecider)
    api.POST("/bets/:id/resolve", betHandler.Resolve)

    e.Start(":" + cfg.Port)
}
```

### Auth Middleware

Supabase signs JWTs with the **JWT secret** (HS256) found in Supabase Dashboard → Settings → API.

```go
// middleware/auth.go
//
// 1. Extract Bearer token from Authorization header
// 2. Verify signature using Supabase JWT secret (HS256)
// 3. Extract `sub` claim — this is the Supabase user ID
// 4. Look up user in MongoDB by supabase_id
// 5. Attach user to request context
// 6. If no token or invalid → 401
// 7. If valid token but no user in MongoDB → 401 (they haven't called POST /users yet)
```

The `POST /users` endpoint is the only **public** endpoint. It's called once right after Supabase signup to create the user profile in MongoDB. It still needs the Supabase JWT to extract the `supabase_id`, but it doesn't require an existing MongoDB user.

### Error Handling

All errors returned as consistent JSON:

```json
{
  "error": {
    "code": "INSUFFICIENT_POINTS",
    "message": "You need at least 10 points to place this wager"
  }
}
```

Standard HTTP status codes:
| Status | When                                                        |
| ------ | ----------------------------------------------------------- |
| 400    | Validation error — bad input, missing fields                |
| 401    | No token, invalid token, or user doesn't exist              |
| 403    | Not a group member, not the decider, not the creator        |
| 404    | Bet/group/invite not found                                  |
| 409    | Duplicate — already wagered, already invited, already member |
| 500    | Unexpected server error                                     |

Define custom error types in the service layer that the handler maps to status codes:

```go
// service/errors.go
var (
    ErrNotFound           = errors.New("not found")
    ErrForbidden          = errors.New("forbidden")
    ErrAlreadyExists      = errors.New("already exists")
    ErrInsufficientPoints = errors.New("insufficient points")
    ErrBetNotOpen         = errors.New("bet is not open")
    ErrAlreadyWagered     = errors.New("already wagered on this bet")
    ErrDeciderCannotWager = errors.New("decider cannot wager")
    ErrCannotSelfDecide   = errors.New("creator cannot be decider")
    ErrNoOpposingSide     = errors.New("cannot resolve with no opposing wagers")
)
```

### Config

```go
// config/config.go
type Config struct {
    Port              string  // default "8080"
    MongoURI          string  // required
    MongoDB           string  // default "youwont"
    SupabaseJWTSecret string  // required — from Supabase dashboard
    StartingPoints    int     // default 1000
    MinWager          int     // default 10
}
```

Loaded from environment variables via `os.Getenv` + godotenv for local development.

---

## 3. Frontend Architecture (React Native / Expo)

### Project Structure

```
app/                              # Expo Router — file-based routing
├── _layout.tsx                   # Root layout (providers, auth guard)
├── (auth)/
│   └── login.tsx                 # Sign in / sign up
├── (tabs)/
│   ├── _layout.tsx               # Tab bar layout
│   ├── index.tsx                 # Home — balance, recent activity
│   ├── groups.tsx                # My groups list
│   └── notifications.tsx         # Notification feed (new tab)
├── group/
│   └── [id].tsx                  # Group detail — members, bets
├── bet/
│   └── [id].tsx                  # Bet detail — pool, wagers
├── create-bet/
│   └── [groupId].tsx             # Create bet form
├── create-group.tsx              # Create group form
├── invite/
│   └── [id].tsx                  # View & accept/decline invite
└── modal.tsx                     # Reusable modal wrapper

api/                              # Raw API functions (plain async, no hooks)
├── client.ts                     # Base fetch wrapper, auth header injection
├── users.ts                      # getMe(), searchUsers()
├── groups.ts                     # createGroup(), getGroups(), getGroup(), joinByCode()
├── bets.ts                       # createBet(), getBet(), getBetsByGroup(), placeWager(), resolve(), changeDecider()
├── invites.ts                    # sendInvite(), getMyInvites(), acceptInvite(), declineInvite()
└── notifications.ts              # getNotifications(), getUnreadCount(), markRead(), markAllRead()

hooks/                            # TanStack Query wrappers
├── use-user.ts                   # useMe(), useSearchUsers()
├── use-groups.ts                 # useGroups(), useGroup(), useCreateGroup(), useJoinByCode()
├── use-bets.ts                   # useBets(), useBet(), useCreateBet(), usePlaceWager(), useResolve()
├── use-invites.ts                # useMyInvites(), useSendInvite(), useAcceptInvite(), useDeclineInvite()
└── use-notifications.ts          # useNotifications(), useUnreadCount(), useMarkRead()

components/                       # Shared UI components
├── ui/                           # Base primitives (icon-symbol, etc.)
└── ...                           # Feature components as needed

lib/
├── supabase.ts                   # Supabase client (auth only)
└── query-client.ts               # TanStack Query client configuration

constants/
└── theme.ts                      # Colors, fonts

ctx.tsx                           # Auth context (session provider)
```

### API Client

```typescript
// api/client.ts
//
// - Base URL from environment config (points to Go API)
// - Automatically attaches Supabase JWT to every request:
//     Authorization: Bearer <token>
// - Handles JSON serialization/deserialization
// - Throws typed errors that match the backend error format
// - All functions are plain async — no React, no hooks

const API_BASE = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
    const session = await supabase.auth.getSession();
    const token = session.data.session?.access_token;

    const res = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...(token && { Authorization: `Bearer ${token}` }),
            ...options?.headers,
        },
    });

    if (!res.ok) {
        const body = await res.json();
        throw new ApiError(res.status, body.error.code, body.error.message);
    }

    return res.json();
}
```

### TanStack Query Setup

```typescript
// lib/query-client.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 1000 * 30,     // 30 seconds — fresh enough for a social app
            retry: 1,
        },
    },
});
```

Provider added to the root layout:

```typescript
// app/_layout.tsx
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '@/lib/query-client';

export default function RootLayout() {
    return (
        <QueryClientProvider client={queryClient}>
            <ThemeProvider ...>
                <AuthProvider>
                    <RootLayoutNav />
                </AuthProvider>
            </ThemeProvider>
        </QueryClientProvider>
    );
}
```

### Hook Pattern

Each hook file follows the same pattern — queries for reads, mutations for writes:

```typescript
// hooks/use-bets.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as betsApi from '@/api/bets';

// Read — auto-fetches, caches, refetches on focus
export function useBets(groupId: string, status?: string) {
    return useQuery({
        queryKey: ['bets', groupId, status],
        queryFn: () => betsApi.getBetsByGroup(groupId, status),
    });
}

export function useBet(betId: string) {
    return useQuery({
        queryKey: ['bet', betId],
        queryFn: () => betsApi.getBet(betId),
    });
}

// Write — call mutate(), invalidate related caches on success
export function usePlaceWager(betId: string) {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (data: { side: string; amount: number }) =>
            betsApi.placeWager(betId, data),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['bet', betId] });
            qc.invalidateQueries({ queryKey: ['me'] });  // points changed
        },
    });
}
```

### Notifications Tab

Add a third tab for notifications (replacing the current 2-tab layout):

```
Tabs: [ Home ] [ Groups ] [ Notifications (🔔 badge) ]
```

The badge shows the unread count via `useUnreadCount()`, which polls or refetches on app focus.

---

## 4. API Contracts

Request and response JSON for every endpoint. This is what the frontend builds against.

### Users

#### `POST /users` — Create profile after signup
```
Auth: Bearer token (Supabase JWT) — extracts supabase_id from sub claim
```
Request:
```json
{
  "name": "Sam",
  "username": "sam"
}
```
Response `201`:
```json
{
  "id": "665a...",
  "supabase_id": "uuid-from-supabase",
  "name": "Sam",
  "username": "sam",
  "avatar_url": null,
  "points": 1000,
  "created_at": "2026-03-05T10:00:00Z"
}
```
Errors: `400` username taken or invalid format, `409` user already exists

---

#### `GET /users/me` — Current user profile
Response `200`:
```json
{
  "id": "665a...",
  "name": "Sam",
  "username": "sam",
  "avatar_url": null,
  "points": 1240,
  "created_at": "2026-03-05T10:00:00Z"
}
```

---

#### `GET /users/search?q=ori` — Search by username
Response `200`:
```json
{
  "users": [
    {
      "id": "665b...",
      "name": "Orion",
      "username": "orion",
      "avatar_url": null
    }
  ]
}
```
Note: Does not return points (that's private). Max 20 results.

---

### Groups

#### `POST /groups` — Create a group
Request:
```json
{
  "name": "The Squad",
  "description": "Our main friend group."
}
```
Response `201`:
```json
{
  "id": "665c...",
  "name": "The Squad",
  "description": "Our main friend group.",
  "invite_code": "XK9F2M",
  "created_by": "665a...",
  "members": [
    {
      "user_id": "665a...",
      "name": "Sam",
      "username": "sam",
      "avatar_url": null,
      "role": "ADMIN",
      "joined_at": "2026-03-05T10:00:00Z"
    }
  ],
  "created_at": "2026-03-05T10:00:00Z"
}
```

---

#### `GET /groups` — List my groups
Response `200`:
```json
{
  "groups": [
    {
      "id": "665c...",
      "name": "The Squad",
      "description": "Our main friend group.",
      "invite_code": "XK9F2M",
      "member_count": 5,
      "open_bet_count": 2,
      "created_at": "2026-03-05T10:00:00Z"
    }
  ]
}
```
Note: List view returns summary (counts, not full members/bets). Detail view returns everything.

---

#### `GET /groups/:id` — Group detail
Response `200`:
```json
{
  "id": "665c...",
  "name": "The Squad",
  "description": "Our main friend group.",
  "invite_code": "XK9F2M",
  "created_by": "665a...",
  "members": [
    {
      "user_id": "665a...",
      "name": "Sam",
      "username": "sam",
      "avatar_url": null,
      "role": "ADMIN",
      "joined_at": "2026-03-05T10:00:00Z"
    },
    {
      "user_id": "665b...",
      "name": "Orion",
      "username": "orion",
      "avatar_url": null,
      "role": "MEMBER",
      "joined_at": "2026-03-06T14:30:00Z"
    }
  ],
  "stats": {
    "total_bets": 5,
    "open_bets": 2,
    "resolved_bets": 2,
    "canceled_bets": 1
  },
  "created_at": "2026-03-05T10:00:00Z"
}
```
Errors: `403` not a member, `404` group not found

---

#### `POST /groups/join` — Join via invite code
Request:
```json
{
  "invite_code": "XK9F2M"
}
```
Response `200`: Same shape as `GET /groups/:id`

Errors: `404` invalid code, `409` already a member

---

### Invites

#### `POST /groups/:id/invites` — Send invite
Request:
```json
{
  "user_id": "665b..."
}
```
Response `201`:
```json
{
  "id": "665d...",
  "group_id": "665c...",
  "group_name": "The Squad",
  "invited_by": "665a...",
  "invited_by_name": "Sam",
  "invitee_id": "665b...",
  "status": "PENDING",
  "created_at": "2026-03-05T10:00:00Z"
}
```
Errors: `403` not a group member, `409` already invited or already a member

---

#### `GET /invites` — My pending invites
Response `200`:
```json
{
  "invites": [
    {
      "id": "665d...",
      "group_id": "665c...",
      "group_name": "The Squad",
      "invited_by": "665a...",
      "invited_by_name": "Sam",
      "status": "PENDING",
      "created_at": "2026-03-05T10:00:00Z"
    }
  ]
}
```

---

#### `POST /invites/:id/accept`
Response `200`:
```json
{
  "id": "665d...",
  "status": "ACCEPTED",
  "group_id": "665c..."
}
```
Frontend navigates to the group on success.

---

#### `POST /invites/:id/decline`
Response `200`:
```json
{
  "id": "665d...",
  "status": "DECLINED"
}
```

---

### Notifications

#### `GET /notifications?page=0&limit=20` — Notification feed
Response `200`:
```json
{
  "notifications": [
    {
      "id": "665e...",
      "type": "BET_WON",
      "ref_type": "bet",
      "ref_id": "665f...",
      "message": "You won 175 pts on \"Sam won't hike the Y\"!",
      "read": false,
      "created_at": "2026-03-05T10:00:00Z"
    },
    {
      "id": "665g...",
      "type": "GROUP_INVITE",
      "ref_type": "invite",
      "ref_id": "665d...",
      "message": "Orion invited you to The Squad",
      "read": false,
      "created_at": "2026-03-04T14:00:00Z"
    }
  ],
  "has_more": true
}
```

---

#### `GET /notifications/unread-count`
Response `200`:
```json
{
  "count": 3
}
```

---

#### `POST /notifications/:id/read`
Response `200`:
```json
{
  "id": "665e...",
  "read": true
}
```

---

#### `POST /notifications/read-all`
Response `200`:
```json
{
  "updated": 3
}
```

---

### Bets

#### `POST /groups/:id/bets` — Create bet with opening wager
Request:
```json
{
  "title": "Sam won't hike the Y before spring",
  "description": "Sam has lived in Provo for 2 years and still hasn't...",
  "end_date": "2026-04-01T23:59:59Z",
  "decider_id": "665b...",
  "opening_wager": {
    "side": "FOR",
    "amount": 100
  }
}
```
Response `201`:
```json
{
  "id": "665f...",
  "group_id": "665c...",
  "title": "Sam won't hike the Y before spring",
  "description": "Sam has lived in Provo for 2 years and still hasn't...",
  "creator": {
    "id": "665a...",
    "name": "Sam",
    "username": "sam",
    "avatar_url": null
  },
  "decider": {
    "id": "665b...",
    "name": "Orion",
    "username": "orion",
    "avatar_url": null
  },
  "end_date": "2026-04-01T23:59:59Z",
  "status": "OPEN",
  "winning_side": null,
  "wagers": [
    {
      "id": "665h...",
      "user": {
        "id": "665a...",
        "name": "Sam",
        "username": "sam",
        "avatar_url": null
      },
      "side": "FOR",
      "amount": 100,
      "placed_at": "2026-03-05T10:00:00Z"
    }
  ],
  "pool": {
    "total": 100,
    "for_total": 100,
    "against_total": 0,
    "for_count": 1,
    "against_count": 0
  },
  "resolved_at": null,
  "created_at": "2026-03-05T10:00:00Z"
}
```
Errors: `400` validation, `403` not a member, `400` creator is decider, `400` insufficient points

---

#### `GET /groups/:id/bets?status=OPEN` — List bets in group
Response `200`:
```json
{
  "bets": [
    {
      "id": "665f...",
      "title": "Sam won't hike the Y before spring",
      "description": "Sam has lived in Provo for 2 years...",
      "status": "OPEN",
      "winning_side": null,
      "end_date": "2026-04-01T23:59:59Z",
      "wager_count": 4,
      "pool": {
        "total": 425,
        "for_total": 225,
        "against_total": 200,
        "for_count": 3,
        "against_count": 1
      },
      "created_at": "2026-03-05T10:00:00Z"
    }
  ]
}
```
Note: List view omits full wager details and user objects — just counts and totals.

---

#### `GET /bets/:id` — Bet detail
Response `200`: Same full shape as the create response (with hydrated user objects on creator, decider, and each wager).

Errors: `403` not a member of the bet's group, `404` not found

---

#### `POST /bets/:id/wagers` — Place a wager
Request:
```json
{
  "side": "AGAINST",
  "amount": 200
}
```
Response `201`: Full bet object (same shape as detail) with the new wager included.

Errors: `400` bet not open, `400` below minimum (10 pts), `400` insufficient points, `409` already wagered, `403` decider cannot wager, `403` not a group member

---

#### `PUT /bets/:id/decider` — Change decider
Request:
```json
{
  "decider_id": "665x..."
}
```
Response `200`: Full bet object with updated decider.

Errors: `403` not the creator, `400` bet not open, `400` new decider has wagered, `400` new decider is creator

---

#### `POST /bets/:id/resolve` — Resolve a bet
Request:
```json
{
  "winning_side": "FOR"
}
```
Response `200`:
```json
{
  "id": "665f...",
  "status": "RESOLVED",
  "winning_side": "FOR",
  "resolved_at": "2026-04-02T15:30:00Z",
  "payouts": [
    { "user_id": "665a...", "name": "Sam", "amount": 175, "net": 75 },
    { "user_id": "665b...", "name": "Orion", "amount": 130, "net": 30 }
  ]
}
```
`amount` = total they receive back. `net` = profit (amount - original wager).

Errors: `403` not the decider, `400` bet not open, `400` no opposing wagers
