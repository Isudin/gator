# aggreGATOR

It's a simple CLI program to aggregate and browse RSS feeds.

Go 1.23 and PostgreSQL 16.6 required.

To install, run ***go install*** command from repo's root folder.

To configure the app, create ***.gatorconfig.json*** file in your home directory.
It should contain ***db_url*** key containing a postgres db connection string
and ***current_user_name*** key containing empty string or any text.

Commands:<br/>
-- register - register a new user<br/>
-- login - login to an existing user<br/>
-- reset - reset the database<br/>
-- users - list users<br/>
-- addfeed - add new feed<br/>
-- agg - aggregate feeds<br/>
-- feeds - list feeds<br/>
-- follow - follow a feed<br/>
-- unfollow - unfollow a feed<br/>
-- following - list feeds followed by current user<br/>
-- browse - list posts<br/>