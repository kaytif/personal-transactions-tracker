# Transactions API

A REST API built in Go for tracking transactions in multiple accounts and categories, allowing for transfers in between accounts, and supports multiple users. 

It uses server side session authentication using cookies, middleware and user authorization. Account balances are derived from transactions. Basic atomicity is enforced in transfers and account creation. The project incorporates soft deletion, validation, database migrations, HTTP error handling, and basic pagination.

A use case for the project is as a backend API for applications that need to manage financial accounts, transactions, categories, and transfers across multiple users. For example, banking apps that use user transaction activity to provide more detailed spending analytics. The project could be advanced to become a full accounting ledger system, which would strengthen the use case. Another use case includes personal finance tracking.

I built this project to develop my skills in backend engineering in Go such as API development and design, working with relational databases and SQL, authentication and authorization, 
database transactions, atomicity, testing, migrations, and pagination.

## Tech Stack
The API is deployed on Render:

https://transactions-accounts-tracker.onrender.com

The PostgreSQL database is hosted on Neon.

There is no frontend for this project. The API can be tested using curl, Postman or PowerShell. Sroll down to the end for sample tests in PowerShell.

- Go
- PostgreSQL
- Render — API hosting
- Neon — PostgreSQL hosting

## Features
- Supports multiple users, accounts, and transaction categories
- User registration and login
- Session based authentication using cookies
- User level authorization
- Create, read, update, and delete accounts
- Create, read, update, and delete transactions
- Create, read, update, and delete categories
- Create and view within account transfers
- Atomic database operations
- Soft deletion
- Transaction pagination
- PostgreSQL migrations
- Basic HTTP handler tests

## Designs limitations and improvements
Building the API exposed several design and system limitations that I discovered in the process of building this projects, and hope to continue to make commits to improve the API.

### Transfers as First Class Entities
I have coded each transfer to be stored as two separate transactions, i.e., a negative transaction on the source account and a positive transaction on the destination account. I chose this approach because I wanted transactions to remain the single source of truth for transaction data and derived values such as account balances.

As the project developed, I realized that transfers should also have a first class representation in the database. Without this, the system needs additional logic within the transaction model to identify and group the two transactions belonging to a transfer. This introduces unnecessary complexity and requires the relationship between the transactions to be inferred rather than explicitly stored in the database.

A separate `transfers` table could represent the transfer itself, while the two transaction entries could reference the same transfer. This would allow the two sides of a transfer to be easily grouped and tracked while still keeping transactions as the source of truth for account activity and balances.

This would also make it easier to enforce rules across both sides of a transfer. Currently, the relationship between the two transaction entries is not explicitly stored, making it more difficult to ensure that both sides remain consistent.


### Double-Entry Accounting
The current system is not based on double entry accounting. Transactions directly increase or decrease the balance of an account. A more robust system could instead use a double-entry ledger. This would provide stronger consistency and make the system more suitable for more complex use cases.

### Transaction Immutability
The current API allows regular transactions to be updated and deleted using `PUT` and `DELETE`.

While building the project, I learned that preventing deleting and update is a better design for transactions. Once a transaction has been recorded, it should remain immutable. Instead of modifying or deleting an existing transaction, corrections could be represented by a reversal transaction followed by a new corrected transaction. This prevents mistakes by user and maintains and audit history.

### Future Design
A future version should include:

- immutable financial transactions
- feature to reverse transactions
- first class transfers
- a double-entry ledger model

These changes would improve consistency, traceability, and strengthen the use case.


## API Endpoints
| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/register` | Register a user |
| POST | `/login` | Login |
| GET | `/accounts` | Get accounts |
| POST | `/accounts` | Create an account |
| PUT | `/accounts` | Update an account |
| DELETE | `/accounts` | Delete an account |
| GET | `/transactions` | Get transactions |
| POST | `/transactions` | Create a transaction |
| PUT | `/transactions` | Update a transaction |
| DELETE | `/transactions` | Delete a transaction |
| GET | `/categories` | Get categories |
| POST | `/categories` | Create a category |
| PUT | `/categories` | Update a category |
| DELETE | `/categories` | Delete a category |
| GET | `/transfers` | Get transfers |
| POST | `/transfers` | Create a transfer |

Accounts, transactions, categories, and transfers require authentication.

## Testing
The current automated tests cover unsupported request methods on the HTTP handler.

The deployed API has also been manually tested against the Render deployment and Neon PostgreSQL.

## Testing the Live API
The examples below use PowerShell.

### Register
```powershell
Invoke-WebRequest `
  -Uri "https://transactions-accounts-tracker.onrender.com/register" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"username":"testuser","password":"test123"}' `
  -UseBasicParsing
```

### Login
Login creates a session cookie. The `$session` variable stores the cookie so it can be reused for authenticated requests.

```powershell
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession

Invoke-WebRequest `
  -Uri "https://transactions-accounts-tracker.onrender.com/login" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"username":"testuser","password":"test123"}' `
  -WebSession $session `
  -UseBasicParsing
```

### Create an account

```powershell
Invoke-WebRequest `
  -Uri "https://transactions-accounts-tracker.onrender.com/accounts" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"name":"DBS","balance":1000}' `
  -WebSession $session `
  -UseBasicParsing
```

An opening balance is stored as a transaction rather than directly as a balance on the account.

### Create a Transaction

```powershell
Invoke-WebRequest `
  -Uri "https://transactions-accounts-tracker.onrender.com/transactions" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"account_id":1,"name":"Coffee","amount":-120,"transaction_date":"2026-09-27"}' `
  -WebSession $session `
  -UseBasicParsing
```

Positive amounts represent money entering an account and negative amounts represent money leaving an account.

### Get Transactions
Transactions are paginated using `page` and `limit`.

```powershell
Invoke-WebRequest `
  -Uri "https://transactions-accounts-tracker.onrender.com/transactions?page=1&limit=10" `
  -Method Get `
  -WebSession $session `
  -UseBasicParsing
```

## Note on AI usage
I wrote and designed the project myself but took help from AI tools and tutorials for guidance, learning syntax, and as a thought partner. I used judgement to use AI only where it would add value without hurting my learning.  
