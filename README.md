# Project description  
This is a example wallet service. User stories:
- I want to be able to send payments from one account to another
- I want to be able to see all payments.  
- I want to be able to see all accounts.  
# details  
- Written with go-kit as per instructions.  
- Test with database migrations for business logic
- Postgres for SQL Database
- Pagination is omitted for simplicity
- Sigma float comparison omitted for simplicity
# Usage
- USE config.json file for service config, use configTest.json for tests config  
- Postgres schema lies in `db/postrgres/001_schema.down.sql`
- Tests `go test coinsWallet/service`    
- Run `go run main.go`  

Creating transaction:

`curl -X POST \
   http://127.0.0.1:7000/wallet/v1/transaction \
   -F 'sender=mr house' \
   -F 'receiver=yes man' \
   -F amount=1000 \
   -F currency=CAPS`
   
Getting transaction list:

 `curl -X GET \
   http://127.0.0.1:7000/wallet/v1/transaction \`
   
Getting accounts list:
 `curl -X GET \
   http://127.0.0.1:7000/wallet/v1/account \`
