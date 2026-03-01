# GoAuth Integration with DocVault

## Shared Database

Both goAuth and DocVault use the same SQLite database:
- **Path**: `../docVault/docvault.db`
- **Table**: `users`

## Database Schema

```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
