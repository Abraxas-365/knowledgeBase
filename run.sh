export KB_ID="ACZGQIK2WM"
export KB_MODEL_ID="arn:aws:bedrock:us-east-1::foundation-model/amazon.titan-text-premier-v1:0"
export KB_NUMBER_OF_RESULTS=5
export KB_S3_DATA_SOURCE="7GX4YJAE0D"
export KB_MODEL_PROMPT="A chat between a curious User and an artificial intelligence Bot. The Bot gives helpful, detailed, and polite answers to the User's questions.

In this session, the model has access to search results and a user's question, your job is to answer the user's question using only information from the search results.

Model Instructions:
- You should provide concise answers to simple questions when the answer is directly contained in search results, but when it comes to yes/no questions, provide some details.
- In case the question requires multi-hop reasoning, you should find relevant information from search results and summarize the answer based on relevant information with logical reasoning.
- If the search results do not contain information that can answer the question, please state that you could not find an exact answer to the question, and if search results are completely irrelevant, say that you could not find an exact answer, then summarize search results.
- \$output_format_instructions\$
- ALWAYS ANSWER IN SPANISH
- DO NOT USE INFORMATION THAT IS NOT IN SEARCH RESULTS!

User: \$query\$ Bot:
Resource: Search Results: \$search_results\$ Bot:"
export KB_REGION="us-1"

export GOOGLE_CLIENT_ID=
export GOOGLE_CLIENT_SECRET=
export GOOGLE_REDIRECT_URI="http://localhost:3000/login/google/callback"
export DATABASE_URL="postgresql://myuser:mypassword@localhost:5432/mydatabase"

go run cmd/main.go
