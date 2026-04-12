CREATE TABLE tokens(
  user_id UUID NOT NULL REFERENCES users,
  token_hash BYTEA NOT NULL PRIMARY KEY,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  capabilities BYTEA
);
