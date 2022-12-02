CREATE TABLE IF NOT EXISTS internal.mailinglist
(
  email      VARCHAR PRIMARY KEY NOT NULL UNIQUE,
  subscribed timestamptz         NOT NULL DEFAULT NOW()
)