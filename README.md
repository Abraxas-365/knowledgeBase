# KnowledgeBase Rag

This project is a serverless application designed to create a knowledge base system
where administrators can upload and delete data, also add more admins or block emails
trying to register to the admin panel. Users can query this knowledge base, with query
restrictions applied per IP address. The project is built using Go and leverages various
AWS services for storage and authentication.

## Backend
- Go Fiber Framework: Used for building the web server and handling HTTP requests.
- PostgreSQL: Database for storing user and session data.
- AWS S3: Used for storing knowledge base data.
- AWS Bedrock: Provides runtime for executing knowledge base queries.
- Google OAuth: Integrated for user authentication. using my [toolkit](https://github.com/Abraxas-365/toolkit)

## Features
- User Management: Admins can promote users to admin role, blacklist users, and delete users.
- Knowledge Base Management: Admins can upload and delete knowledge base data.
- User Queries: Users can query the knowledge base with rate limiting applied per IP address.
- OAuth Authentication: Google OAuth is used for secure user login and session management.
- Session Management: Secure cookies are used for session management.

# Getting Started
### Prerequisites
- Go 1.20+
- Docker
- AWS Account
- PostgreSQL
- Google OAuth Credentials


### Project Setup
1. **Deploy Resources**: Follow the instructions in this [Medium article](https://medium.com/@miramnair/develop-and-deploy-a-serverless-rag-solution-with-amazon-bedrock-agents-knowledge-base-and-ef8a1818bc1e) by Meera Nair.
  - We will be just using the knowledge base, so use the `cloudfromation` files in [build](./build) instead of the provided in [Medium article](https://medium.com/@miramnair/develop-and-deploy-a-serverless-rag-solution-with-amazon-bedrock-agents-knowledge-base-and-ef8a1818bc1e)
  - The steps are the same

2. **Environment Variables**: Set the following environment variables:

```
KB_ID=KowledgeBaseId
KB_MODEL_ID=knowledgebase llm model
KB_NUMBER_OF_RESULTS=number of results for the knowledge base
KB_MODEL_PROMPT=prompt for the model, use aws guide to create the prompt
KB_REGION= knowledge base region
KB_S3_DATA_SOURCE=knowledgebase data source id

REDIRECT_AFTER_LOGIN=http://localhost:30001/home
ALLOW_ORIGINS=http://localhost:3001,http://localhost:3000
GOOGLE_CLIENT_ID=google client id
GOOGLE_CLIENT_SECRET=google secret id
GOOGLE_REDIRECT_URI=redirect url
DATABASE_URL= databsase uri

```
[Env example](run.sh)

3. **Database Migration**: Ensure your PostgreSQL database is set up and migrations are applied. [migrations](./migrations/)

4. **Run without building**:
```bash
sh run.sh
```
5.	`Initial Admin Setup`: Manually change the `user.is_admin` field to true for the first admin user in the database after oauth registration.

## API Endpoints
### Authentication
- Login with Google: `/login/google`
- Google Callback: `/login/google/callback`
- Logout: `/logout`
### Knowledge Base
- Upload Data: `/generate-presigned-url` (POST)
- List Objects: `/list-objects` (GET)
- Delete Object: `/delete-object` (DELETE)
- Query: `/chat/complete-answer` (POST)
### User Management
- List Users: `/users` (GET)
- Promote to Admin: `/users/promote-to-admin` (POST)
- Delete User: `/users/:id` (DELETE)
- Blacklist Management: `/users/blacklist` (GET, POST, DELETE)

## Security
- Rate Limiting: Implemented to restrict the number of requests per IP address to prevent abuse.
- OAuth: Secures user authentication and session management.
- Cookie: Uses secure cookies for session management.
