# Runnning the Code

Install dependencies:
```
go get
```

Run the server:
```
go run main.go
```

## The Calls
`/signup`
```cURL
curl --location --request POST 'http://localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Michael Colby",
    "email": "email@email.com",
    "password": "supersecret123!"
}'
```

`/login`
```cURL
curl --location --request POST 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Michael Colby",
    "email": "email@email.com",
    "password": "supersecret123!"
}'
```

`/update`
```cURL
curl --location --request POST 'http://localhost:8080/update' \
--header 'Authorization: Bearer ABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABC' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Michael Crolby",
    "email": "email@email.com",
    "password": "supersecret123!"
}'
```

`/delete`
```cURL
curl --location --request DELETE 'http://localhost:8080/delete' \
--header 'Authorization: Bearer ABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABC'
```

`/logout`
```cURL
curl --location --request GET 'http://localhost:8080/logout' \
--header 'Authorization: Bearer ABCABCABCABCABCABCABCABCABCABCABCABCABCABCABCABC'
```

# My Train of Thought 🚂 
Firstly, I would like to justify the use of third-party libraries. I find that it's best not to reinvent the wheel when you can, but we also have to tread the line of making sure that we do not get stuck with outdated/buggy/insecure code. This is why I always check to see how "alive" a project is first by checking the last time an update was made. Then checking to see how often issues are being handled by the crew or community. 

I felt confident in using `go-oauth2/oauth2` since it had a recent PR in April. Importing this project allowed me to worry less about some of the mechanics of OAuth2 so that I could focus on writing the necessary functions. I did not, however, use a third-party library for REST since I felt that the limited amount of operations I needed to perform didn't warrant lugging around all of the robust features of such a library. Fortunately, Go has just what I needed to make the REST work "good enough."

## User Storage
I chose to use the `clientStore` from the oauth2 library as my storage for my users. This allowed me not to have to worry about also setting up data storage interfaces. If I were to operationalize this I would have without a doubt stored users in a more conventional way such as an SQL database. The oauth2 library has several integrations with various SQL and NoSQL storage solutions for storing tokens as well. I would have also written such an interface to allow for dependency inject allowing for swapping out the storage solution if need be. This ability to swap out how things are stored allows for better mocking of the methods so that we can test easier.

## `/signup`
For this endpoint, I wanted to use a simple POST with the necessary data for a user. If the user is able to be created, the endpoint will respond with `202` that the sign up has been accepted. From here a user is now able to login with their email and password.

## `/login`
For this endpoint, the user submits their email and password. The request body is then unmarshalled, and placed into an instance of my `User` data model. This is then used to call get on the server to access the token method. Since we previously made the user a client it uses this information to get the token. If entered correctly, an access token is provided and proxied back to the user in the response.

## `/update`
For this endpoint, we get to check for the token the user just logged in to get. We make use of the oauth2 server's `ValidationBearerToken()` function that will check on the access token for us. Since it is returned to the client and it is told it is a Bearer token the request is expected to have `"Authorization": "Bearer ..."` in the header. While the validation function takes care of that it is important to note this header is required to get things done on this endpoint (as well as `/delete`). I made this one accept a POST since that felt more consistent with the other endpoints, but UPDATE is just as valid. In order to be RESTful you just need to be consistent.

It should also be noted that this endpoint will not allow the user to update their email as it is their id. This is mostly a limitation based around the storage solution and is in no way good UX.

## `/delete`
For this endpoint, it's fairly simple. If the user is validated via the access token, the clientStore sets the client as `nil` and then the token is removed from the tokenStore. The user is now deleted and can no longer log in. I chose to use the DELETE method since it's not as common it will feels like it's harder to accidentally do.

## `/logout`
For this endpoint, it is also very simple. If the token passes, then it is removed from the tokenStore. That's it. A message informing that the log out has happened is returned. Removing the token means that the client (user) no logger has a token to validate with so they will be blocked out of the calls needing authentication until `/login` is called again and a new token is issued.

# Moving Forward
My first order of business for this would be to get a better storage solution for users and tokens. Currently, this solution does not offer much in the way of scalability since a user would need to be locked to a particular instance which isn't very reliable when load balancers are also at play. I would also make sure that the error handling and REST methods are really tight and consistent. Since I am getting to over the recommended time, I wasn't able to write test, but I would also make sure those are written. 