## Simple-task

    Used 
        Goland Ide 11.0.9.1
        Docker Desktop 3.3.3
        Windows 10
        golang 1.15+

## Initialization
    
    Docker Desktop must be installed for this task.
    Firstly please clone repository to your machine. 
    Below command to create and run simple-task image

        $ docker-compose up
    
## Testing golang

    After initialization finished - below command must be used for testing.
    
        $ docker exec processing go test -bench=. ./...

    Testing 
        1. Bad request
        2. Not acceptable source type
        3. No State
        4. Wrong state
        5. No transaction Id
        6. No amount
        7. Used transaction Id
        8. Not logged 
        9. Not registered
        10. Win state
        11. Lose state
        12. Negative balance check
        
        13. Benchmark 1000 random generated messages

## Testing postman

    After server and database started below url and commands can be used for testing

        1. Register new user for testing
           http://127.0.0.1/api/register

            Headers: 
                "Accept" "application/json"
                "Content-type" "application/json"

            Below json object must be posted for registration
            {
                "UserId":"NewUserID"
            }
        
        2. Send post request for testing
           http://127.0.0.1/api/processing

            Headers: 
                "Accept" "application/json"
                "Content-type" "application/json"
                "Source-type"  "server" || "game" || "payment"
                "Authorization"  "NewUserId"
            
                States can be "win" || "lose"
                Transaction id must be generated unique for every request
                
            Below json object must be posted for processing
            {
                "state": "lose", 
                "amount": "8.78", 
                "transactionId": "some generated identificator"
            }

        More: 
            random generated 1000 messages
            https://www.json-generator.com/
            
            generating script
            
            [
             '{{repeat(1000,1000)}}',
             {
               state:'{{random("win", "lose")}}',
               amount:'{{floating(0, 100, 2, "0,0.00")}}',
               transactionId:'{{guid()}}'
             }
            ]