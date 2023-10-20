TODO:
- [WONT] Set up the serverless framework for local testing, etc.
- [DONE] Should deploy with Serverless deploy
- [DONE] Spin up DynamoDB
- [DONE] Enable the addition of a new user with fiels like userid, email. name, DOB, etc.
- [DONE] Fetch user information based on UserID
- [DONE] Modify exusting user details using UserID
- [DONE] Remove a user record based on UserID
- [DONE] /users (POST) to add a new user
- [DONE] /users/ GET to retrieve details of the user
- [DONE] PUT
- [DONE] DELETE
- [DONE] Error handling with descriptive error code and message
- [TODO] Run all commands in docker container
- [WONT] Local integration tests
- [TODO] Add test to ensure client is working
- [TODO] Interesting decisions section for readme
- [DONE] Update date fields to strings instead of time objects
- [DONE] move repo functions into own files
- [DONE] add proper check for existing user with email
- [TODO] Update "LastModified" and "CreatedAt" to actually work ( or drop them )

using global state / init: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
