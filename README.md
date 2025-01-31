# aggreGATOR

It's a simple CLI program to aggregate and browse RSS feeds.

Go 1.23 and PostgreSQL 16.6 required.

To install, run ***go install*** command from repo's root folder.

To configure the app, create ***.gatorconfig.json*** file in your home directory.
It should contain ***db_url*** key containing a postgres db connection string
and ***current_user_name*** containing empty string or any text.

Commands:
-- register - register a new user
-- login - login to an existing user
-- reset - reset the database
-- users - list users
-- addfeed - add new feed
-- agg - aggregate feeds
-- feeds - list feeds
-- follow - follow a feed
-- unfollow - unfollow a feed
-- following - list feeds followed by current user
-- browse - list posts