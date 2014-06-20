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
- set the following environment variables (I recommend dotenv):

	export DB_NAME=<my database name>
	export DB_USER=<my database user>
	export DB_PASS=<my database password>

	export TEST_DB_NAME=<something different from DB_NAME>
	export TEST_DB_USER=<my database user>
	export TEST_DB_PASS=<my database password>

	# create a private/public keypair using whatever your OS uses
	# this is used for generating JWT tokens

	export PRIVATE_KEY = "$(pwd)/keys/sample_key"
	export PUBLIC_KEY = "$(pwd)/keys/sample_key.pub"

	# optional, runs on 5000 by default

	export PORT = 6000

- `npm install bower && bower install`
- Run the ./photoshare executable.

