# youwont — Product Design Document

---

## 1. Jobs to Be Done

> "When someone in my friend group makes a bold claim or prediction, I want a quick way to bet points on it so we can track who was right, talk trash, and have fun — without dealing with real money."

---

## 2. Personas

### The Instigator — "Orion"
- Always starting arguments and making bold claims
- Wants to create bets fast, right in the moment, from his phone
- Cares about proving people wrong
- Will create 80% of the bets

### The Gambler — "Allison"
- Doesn't create many bets but wagers on everything
- Wants to see what's active, pick a side, and drop points on it
- Checks the app to see how her bets are doing
- Competitive, wants to climb the leaderboard

### The Judge — "Paul"
- Trusted friend who gets picked as the decider
- Needs a clear, simple way to call the outcome
- Doesn't want drama — wants the resolution to feel fair and final
- Might not open the app often, so notifications matter

---

## 3. Core User Journeys

### Journey 1: "Someone says something bold" (Create a Bet)

Orion and his friends are at lunch. Sam says "I'll hike the Y before April, easy."
Orion pulls out his phone.

1. Opens youwont, taps into "The Squad" group
2. Hits the + button to create a new bet
3. Types the title: "Sam won't hike the Y before spring"
4. Adds a short description for context
5. Picks an end date (April 1)
6. Selects Paul as the decider (someone neutral)
7. Places his own opening wager — 100 pts on FOR
8. Submits. The bet appears in the group instantly
9. Everyone at the table gets a notification and starts wagering

**What he's feeling:** Excited, wants this to be fast (under 30 seconds), wants to share it immediately.

---

### Journey 2: "I want in on this" (Place a Wager)

Allison gets a push notification: "New bet in The Squad: Sam won't hike the Y before spring."

1. Taps the notification, lands on the bet detail screen
2. Sees the title, description, who created it, who the decider is
3. Sees the current pool — 175 pts FOR, 0 pts AGAINST
4. She thinks Sam will actually do it. Taps "Place a Wager"
5. Picks her side: AGAINST
6. Enters her amount: 200 pts
7. Confirms. Her wager appears in the list. Pool updates in real time
8. She checks her point balance — still has 1,900 pts left

**What she's feeling:** Wants to see the odds quickly, wants placing a wager to be 2-3 taps max.

---

### Journey 3: "Time to call it" (Resolve a Bet)

It's April. Sam never hiked the Y. Paul gets a notification:
"The bet 'Sam won't hike the Y before spring' is past its end date. You're the decider — time to call it."

1. Paul opens the app, taps the notification
2. Sees the bet detail — all the wagers, the pool breakdown
3. Taps "Resolve Bet"
4. Selects the winning side: FOR (Sam indeed did not hike it)
5. Confirms his decision
6. Points are distributed automatically:
   - FOR bettors split the AGAINST pool proportionally to their wagers
   - Everyone gets a notification with their result
7. The bet is marked RESOLVED with a trophy showing "FOR wins"

**What he's feeling:** Wants it to be simple and unambiguous. Pick a side, done. Doesn't want to do math.

---

### Journey 4: "Let's get the group going" (Create/Join a Group)

It's the start of the semester. Orion wants to set up a betting group for his roommates.

1. Opens the app, goes to the Groups tab
2. Taps the + button to create a new group
3. Names it "The Squad", adds a description
4. Taps "Invite Members"
5. Searches for friends by username, sends invites to Sam, Allison, Paul
6. Sam gets an in-app notification (bell icon shows a badge): "Orion invited you to The Squad"
7. Sam taps the notification, sees the group name and description, taps "Accept"
8. He's in. He can see the group and start betting
9. Allison accepts later that evening. Paul accepts the next morning. The group fills up organically.

**Fallback:** The group still has an invite code (SQUAD2024) that can be shared over text for people who aren't on the platform yet, or if someone just wants to share a code quickly.

**What he's feeling:** Wants setup to take less than a minute. Inviting friends by username feels natural — like adding people to a group chat.

---

## 4. Feature Extraction

| Feature                        | Source Journey        | Notes                                          |
| ------------------------------ | --------------------- | ---------------------------------------------- |
| Create a group                 | Journey 4             | Name, description, auto-generated invite code  |
| Invite users to a group        | Journey 4             | Search by username, send invite                |
| Accept / decline invites       | Journey 4             | Pending invites with accept/decline actions    |
| In-app notification feed       | Journey 2, 3, 4       | Bell icon, badge count, tap to navigate        |
| Join a group via invite code   | Journey 4             | Fallback for sharing outside the app           |
| View my groups                 | All                   | List with member count, active bet count       |
| View group detail              | Journey 1, 2          | Members, bets list, filters                    |
| Create a bet                   | Journey 1             | Title, description, end date, decider, opening wager |
| View bet detail                | Journey 2, 3          | Pool breakdown, wager list, status             |
| Place a wager                  | Journey 2             | Pick side, enter amount, confirm               |
| Resolve a bet (decider only)   | Journey 3             | Pick winning side, confirm                     |
| Point system / balances        | Journey 2, 3          | Starting points, deduct on wager, payout on resolve |
| Point payout calculation       | Journey 3             | Proportional split of losing pool to winners   |
| Push notifications             | Journey 2, 3          | New bet, bet resolved, decider reminder        |
| User profile                   | All                   | Name, username, point balance                  |
| Leave a group                  | —                     | Basic group management                         |
| Leaderboard (per group)        | —                     | Rank members by points or win rate             |
| Bet comments / trash talk      | —                     | Social layer                                   |
| Bet history / stats            | —                     | Past performance                               |
| Deep link invites              | —                     | Tap a link to join, skip code entry             |
| Real-money / crypto payouts    | —                     | Future — regulatory complexity                  |

---

## 5. Prioritization

### MVP (what you need for the app to function)

1. **Auth** — sign up, sign in, sign out *(already done via Supabase)*
2. **User profile** — name, username, point balance (seeded with starting points)
3. **Create a group** — name, description, generates invite code
4. **Invite users to a group** — search by username, send invite
5. **Accept / decline invites** — pending invites screen
6. **Notification feed** — bell icon with badge, shows invites, bet results, wins/losses
7. **Join a group via invite code** — fallback for external sharing
8. **Search users** — find people by username to invite
9. **View my groups** — list screen
10. **View group detail** — members, bets
11. **Create a bet** — title, description, end date, select decider, opening wager
12. **View bet detail** — pool, wagers, status
13. **Place a wager** — pick side, enter amount, deduct from balance
14. **Resolve a bet** — decider picks winner, points distributed

That's the core loop: **create group → create bet → wager → resolve → points move**.

### Post-MVP (nice to have, in rough priority)

15. Push notifications (via Expo push — extends in-app notifications to phone-level alerts)
16. Leaderboard per group
17. Bet comments / trash talk
18. Leave / manage group (kick members, transfer admin)
19. User stats / bet history
20. Deep link invites
21. Cancel a bet (refund all wagers)
22. Edit bet before first outside wager
23. Dark mode
24. Crypto/real-money payouts (way later, if ever)

---

## 5b. Business Rules & Constraints

### Users
- New users start with **1,000 points**
- Username must be **unique**, lowercase, alphanumeric + underscores, 3-20 characters
- Points are a **single global balance** across all groups — not per-group
- A user's point balance can **never go below 0** — enforce this at wager time
- Username is set once during profile creation (post-MVP: allow changes)

### Groups
- Any authenticated user can **create a group**
- The creator is automatically the group **ADMIN**
- A group must have a **name** (1-50 chars) and optionally a **description** (0-200 chars)
- Invite codes are **auto-generated**, unique, 6-8 alphanumeric characters
- A user **cannot be in the same group twice** — deduplicate on join/accept
- No upper limit on group size for MVP, but realistically these are friend groups (3-20 people)
- A user **cannot leave a group if they have OPEN bets** in it (as creator, decider, or wager placer) — resolve or cancel first
- All data is **scoped to group membership** — you must be in a group to see its bets, members, and activity

### Invites
- Only group **members** (any role) can send invites — you have to be in the group to invite someone
- You **cannot invite someone who is already in the group**
- You **cannot invite someone who already has a PENDING invite** to that group
- A user **can decline and be re-invited** later (declining doesn't permanently block)
- Invites do **not expire** for MVP (post-MVP: auto-expire after 7 days)

### Bets
- Any group member can **create a bet** (not restricted to admins)
- A bet requires: **title** (1-100 chars), **description** (0-500 chars), **end date** (must be in the future), **decider** (must be a member of the group)
- The **creator cannot be the decider** — the decider should be someone neutral
- The **decider can be changed** by the bet creator while the bet is OPEN — the new decider must be a group member who hasn't wagered on the bet
- A bet starts in **OPEN** status
- The creator **must place an opening wager** when creating the bet (minimum 10 pts) — no empty bets
- End date is **informational, not enforced** — it doesn't auto-close the bet. It's just when the event is expected to conclude. The decider resolves it manually.

### Wagers
- Only members of the bet's group can **place a wager**
- A user can **only wager once per bet** — pick your side, commit
- The **decider cannot wager** on a bet they're deciding — conflict of interest
- Minimum wager: **10 points**
- Maximum wager: **your current point balance** (no debt)
- You **can wager on a bet you created** (the creator already places the opening wager, and this is their only wager)
- Wagers can only be placed while the bet is **OPEN**
- Wagers are **non-refundable** once placed (except if the bet is canceled — post-MVP feature)
- Points are **deducted immediately** when the wager is placed, not when the bet resolves

### Resolution
- Only the **decider** can resolve a bet
- The decider picks the **winning side**: FOR or AGAINST
- Resolution can happen **at any time** while the bet is OPEN — before or after the end date
- Once resolved, it's **final** — no undo for MVP
- **Payout formula:**
  - Each winner receives their original wager back
  - Plus a share of the total losing pool, proportional to their wager size
  - `payout = original_wager + (original_wager / total_winning_pool) * total_losing_pool`
  - Example: FOR pool is 200 pts (you bet 100, friend bet 100). AGAINST pool is 150 pts. FOR wins. You get: 100 + (100/200) × 150 = **175 pts**
- Points are **credited immediately** upon resolution
- If only one side has wagers (e.g., everyone bet FOR, nobody bet AGAINST), and that side wins, everyone just gets their wager back — there's no losing pool to split. If the empty side "wins," the existing wagers go to... nobody. **Rule: a bet cannot be resolved if one side has zero wagers.** The decider should wait or the bet should be canceled (post-MVP).

### Notifications
- Notifications are created as a **side effect** of domain actions — the frontend never creates them directly
- Notifications are **never deleted**, only marked as read
- The unread count badge reflects `notifications where read = false`
- Tapping a notification **marks it as read** and navigates to the relevant screen

---

## 6. Data Design (MongoDB)

### Design Principles

- **Embed** data that is always read together and won't grow unbounded
- **Reference** (store IDs) for data that is shared or independently queried
- Optimize for the app's read patterns, not for normalization

---

### Collections

#### `users`
```json
{
  "_id": ObjectId,
  "supabase_id": "uuid-from-supabase-auth",
  "name": "Sam",
  "username": "sam",
  "avatar_url": null,
  "points": 1000,
  "created_at": ISODate
}
```
- `supabase_id` links to Supabase Auth — this is how you look up the user after they log in
- `points` is the source of truth for balance — updated on wager placement and bet resolution
- Index on `supabase_id` (unique)

---

#### `groups`
```json
{
  "_id": ObjectId,
  "name": "The Squad",
  "description": "Our main friend group.",
  "invite_code": "SQUAD2024",
  "created_by": ObjectId,          // ref → users._id
  "members": [
    {
      "user_id": ObjectId,          // ref → users._id
      "role": "ADMIN",
      "joined_at": ISODate
    },
    {
      "user_id": ObjectId,
      "role": "MEMBER",
      "joined_at": ISODate
    }
  ],
  "created_at": ISODate
}
```
- Members are **embedded** — a group won't have thousands of members (it's friend groups, probably 3-20 people), and you always want the member list when you load a group
- Index on `invite_code` (unique) for join lookups
- Index on `members.user_id` for "get my groups" queries

**Why not a separate `memberships` collection?** Because the primary access pattern is "load this group with its members." Embedding avoids a join. The tradeoff is that updating a user's name requires updating it in multiple places — but you'll reference by `user_id` and hydrate user details at read time, so this isn't an issue.

---

#### `bets`
```json
{
  "_id": ObjectId,
  "group_id": ObjectId,            // ref → groups._id
  "title": "Sam won't hike the Y before spring",
  "description": "Sam has lived in Provo for 2 years and...",
  "creator_id": ObjectId,          // ref → users._id
  "decider_id": ObjectId,          // ref → users._id
  "end_date": ISODate,
  "status": "OPEN",                // OPEN | RESOLVED | CANCELED
  "winning_side": null,            // FOR | AGAINST | null
  "wagers": [
    {
      "_id": ObjectId,
      "user_id": ObjectId,          // ref → users._id
      "side": "FOR",
      "amount": 100,
      "placed_at": ISODate
    },
    {
      "_id": ObjectId,
      "user_id": ObjectId,
      "side": "AGAINST",
      "amount": 200,
      "placed_at": ISODate
    }
  ],
  "resolved_at": null,
  "created_at": ISODate
}
```
- Wagers are **embedded** inside the bet — they're always read with the bet, and a bet between friends will have maybe 2-15 wagers max
- Index on `group_id` for "get all bets in this group"
- Index on `wagers.user_id` for "get all bets I'm involved in" (if needed for a profile/history page)
- `status` + `winning_side` together describe the bet's lifecycle

**Why embed wagers?** Every time you view a bet, you need the wagers. A friend-group bet will never have hundreds of wagers. Embedding means one read gets you everything. Placing a wager is an `$push` onto the array + a point deduction on the user doc (these two should happen atomically via a transaction).

---

#### `invites`
```json
{
  "_id": ObjectId,
  "group_id": ObjectId,            // ref → groups._id
  "group_name": "The Squad",       // denormalized for display
  "invited_by": ObjectId,          // ref → users._id
  "invited_by_name": "Orion",      // denormalized for display
  "invitee_id": ObjectId,          // ref → users._id
  "status": "PENDING",             // PENDING | ACCEPTED | DECLINED
  "created_at": ISODate
}
```
- Invites are their own collection — they have a lifecycle (pending → accepted/declined) and are independently queried
- Denormalize `group_name` and `invited_by_name` so the invite list renders without extra lookups
- Index on `invitee_id` + `status` for "get my pending invites"
- Index on `group_id` + `invitee_id` (unique) to prevent duplicate invites to the same person

**When an invite is accepted:** update invite status to ACCEPTED, push user into group's members array, and create a notification for the inviter ("Sam joined The Squad"). All in a transaction.

---

#### `notifications`
```json
{
  "_id": ObjectId,
  "user_id": ObjectId,             // ref → users._id (the recipient)
  "type": "GROUP_INVITE",          // see types below
  "ref_type": "invite",            // "invite" | "bet"
  "ref_id": ObjectId,              // links to the invite, bet, etc.
  "message": "Orion invited you to The Squad",
  "read": false,
  "created_at": ISODate
}
```

**Notification types:**
| Type              | When it's created                          | ref_type | What it links to        |
| ----------------- | ------------------------------------------ | -------- | ----------------------- |
| `GROUP_INVITE`    | Someone invites you to a group             | invite   | The invite (accept/decline) |
| `INVITE_ACCEPTED` | Someone accepts your invite                | group    | The group               |
| `BET_CREATED`     | New bet in a group you're in               | bet      | The bet detail          |
| `WAGER_PLACED`    | Someone wagers on a bet you created        | bet      | The bet detail          |
| `BET_RESOLVED`    | A bet you wagered on is resolved           | bet      | The bet detail          |
| `BET_WON`         | You won points from a resolved bet         | bet      | The bet detail          |
| `BET_LOST`        | You lost points from a resolved bet        | bet      | The bet detail          |
| `DECIDER_REMINDER`| A bet you're the decider for is past due   | bet      | The bet detail          |

- Notifications are **read-only from the user's perspective** — no actions on them, just read/unread. If it's a GROUP_INVITE, the app navigates to the invite where the action happens.
- Index on `user_id` + `read` (the primary query: "get my unread notifications")
- `message` is pre-rendered at creation time so the feed loads with zero extra lookups
- Notifications are cheap to create — whenever something happens, fan out a notification to each relevant user

---

### Access Patterns → Queries

| What the app does               | MongoDB query                                              |
| -------------------------------- | ---------------------------------------------------------- |
| Get current user                 | `users.findOne({ supabase_id })`                           |
| Search users by username         | `users.find({ username: { $regex } })` (with index)        |
| Get my groups                    | `groups.find({ "members.user_id": userId })`               |
| Get group detail                 | `groups.findOne({ _id: groupId })` + hydrate user details  |
| Send group invite                | `invites.insertOne(...)` + `notifications.insertOne(...)` — **transaction** |
| Get my pending invites           | `invites.find({ invitee_id, status: "PENDING" })`          |
| Accept invite                    | `invites.updateOne(ACCEPTED)` + `groups.updateOne($push member)` + `notifications.insertOne(INVITE_ACCEPTED)` — **transaction** |
| Join via invite code (fallback)  | `groups.findOneAndUpdate({ invite_code }, { $push: { members } })` |
| Get notifications                | `notifications.find({ user_id }).sort({ created_at: -1 })` |
| Get unread count                 | `notifications.countDocuments({ user_id, read: false })`   |
| Mark notification read           | `notifications.updateOne({ _id }, { read: true })`         |
| Get bets for a group             | `bets.find({ group_id: groupId })`                         |
| Get single bet                   | `bets.findOne({ _id: betId })`                             |
| Create a bet + opening wager     | `bets.insertOne(...)` + `users.updateOne(points - amount)` + fan out `notifications` to group members — **transaction** |
| Place a wager                    | `bets.updateOne({ $push: { wagers } })` + `users.updateOne(points - amount)` — **transaction** |
| Resolve a bet                    | `bets.updateOne({ status, winning_side })` + `users.updateMany(point payouts)` + fan out `notifications` (BET_WON/BET_LOST) — **transaction** |

---

### Transactions (Important)

Four operations **must** be atomic (MongoDB multi-document transactions, requires replica set — Atlas gives you this by default):

1. **Send a group invite**: create invite doc AND create GROUP_INVITE notification for the invitee. (Could skip the transaction here since a missing notification isn't catastrophic, but it's cleaner.)

2. **Accept a group invite**: update invite status to ACCEPTED, push user into group's members array, AND create INVITE_ACCEPTED notification for the inviter.

3. **Place a wager**: deduct points from user AND push wager onto bet. If either fails, neither should happen. Also check `points >= amount` to prevent going negative.

4. **Resolve a bet**: update bet status, distribute points to winners, AND fan out BET_WON/BET_LOST notifications to all participants. The payout math:
   - Total losing pool = sum of losing side's wagers
   - Each winner gets: `(their_wager / total_winning_pool) * total_losing_pool`
   - Plus their original wager back

---

### Hydration Pattern

The documents store `user_id` references, not full user objects. When the API returns data to the frontend, it needs to **hydrate** — replace IDs with user details (name, username, avatar).

Two approaches:
- **Aggregation pipeline** with `$lookup` (MongoDB's version of a join) — done server-side
- **Application-level join** — fetch the bet, collect unique user IDs, batch-fetch those users, merge — done in Go

For a small app like this, either works. Application-level is simpler to write and debug. `$lookup` is fewer round trips.

---

## 7. API Design (Go / Echo)

Based on the MVP features, here are the endpoints the backend needs:

### Auth
The Go backend doesn't handle login/signup — Supabase does that. But every authenticated request should include the Supabase JWT in the `Authorization` header. The backend verifies it and extracts the `supabase_id` to identify the user.

### Endpoints

```
# Users
POST   /users                       — Create user profile (called once after Supabase signup)
GET    /users/me                    — Get current user's profile + points
GET    /users/search?q=             — Search users by username (for inviting)

# Groups
POST   /groups                      — Create a group
GET    /groups                      — List my groups
GET    /groups/:id                  — Get group detail (members, stats)
POST   /groups/join                 — Join a group via invite code (fallback)

# Invites (actions)
POST   /groups/:id/invites          — Send invite to a user
GET    /invites                     — Get my pending invites
POST   /invites/:id/accept          — Accept an invite
POST   /invites/:id/decline         — Decline an invite

# Notifications (read-only feed)
GET    /notifications               — Get my notifications (paginated, newest first)
GET    /notifications/unread-count  — Get unread count (for bell badge)
POST   /notifications/:id/read      — Mark one as read
POST   /notifications/read-all      — Mark all as read

# Bets
GET    /groups/:id/bets             — List bets in a group (?status=OPEN filter)
POST   /groups/:id/bets             — Create a bet (with opening wager)
GET    /bets/:id                    — Get bet detail (with wagers)
POST   /bets/:id/wagers             — Place a wager
PUT    /bets/:id/decider            — Change the decider (creator only, bet must be OPEN)
POST   /bets/:id/resolve            — Resolve a bet (decider only)
```

That's 20 endpoints. Clean separation: invites have actions, notifications are a feed.
