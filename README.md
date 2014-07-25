![ScreenShot](/screenshot.png)

Photoshare is a learning app for Go and Angular.

The final aim:

- User can create posts with photos.
- User can make posts private/public.
- Other users can comment on posts.
- User can switch comments on/off.
- Simple search.
- Email user when someone comments on their posts.
- RSS feeds on all posts.
- Admins can edit/delete all posts and enable/disable user accounts.
- Users can sign up/in with at least one oauth client (Google, Facebook, Twitter)
- Runs on PostgreSQL backend.
- Twitter Boostrap-based design (yawn, but I'm not a designer...)
- Fully tested (or as much as possible) in front and backend.
- We can host the whole thing on Heroku/S3.

Of course not all of the above will happen at once, rather incremental steps.

Getting started
---------------

You need Go (1.3), node.js/npm and PostgreSQL (9.1+).

- `go get && go build`
- Set the correct environment variables. See sample_env for a template.
- `npm install bower && bower install`
- `go get goose`
- Copy and edit db/db.yml.sample to point to the correct databases.
- `goose -env=development up`
- `cd commands/server && go build`
- `./commands/server/server'

Tested on Chrome and Firefox 30+.
