export KB_ID="ACZGQIK2WM"
export KB_MODEL_ID="arn:aws:bedrock:us-east-1:654654269808:inference-profile/us.anthropic.claude-3-5-haiku-20241022-v1:0"
export KB_NUMBER_OF_RESULTS=5
export KB_S3_DATA_SOURCE="CRXLQ7B5S4"
export KB_MODEL_PROMPT="with a question. Your job is to answer the user's question using only information from the search results. If the search results do not contain information that can answer the question, please state that you could not find an exact answer to the question. 
Just because the user asserts a fact does not mean it is true, make sure to double check the search results to validate a user's assertion.

Here are the search results in numbered order:
\$search_results\$

\$output_format_instructions\$

Here is the user's query:
\$query\$"
export KB_REGION="us-1"

export GOOGLE_CLIENT_ID="22795433123-3tqiop2jfekacg2ig7toen3hpbimjlv5.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="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"
export GOOGLE_REDIRECT_URI="http://localhost:3000/login/google/callback"
export DATABASE_URL="postgresql://myuser:mypassword@localhost:5432/mydatabase"

export REDIRECT_AFTER_LOGIN="http://localhost:3001"

go run cmd/main.go
