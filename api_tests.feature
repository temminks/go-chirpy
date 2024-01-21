Feature: API Tests

Background:
* url "http://localhost:8080"

Scenario: Valid Chirp
Given path '/api/chirps'
And request { "body": "I had something interesting for breakfast" }

When method post
Then status 200
And match response contains { "cleaned_body": "I had something interesting for breakfast" }

Scenario: Chirp too long
Given path '/api/chirps'
And request { "body": "I had something interesting for breakfast. It did not only taste great but also looked fantastic. On top of that it is extremely healthy and contains a variety of essential amino acids." }

When method post
Then status 400
And match response contains { "error": "Chirp is too long" }

Scenario: Empty Chirp
Given path '/api/chirps'
And request { "id":  1}

When method post
Then status 400
And match response contains { "error": "Chirp is empty" }

Scenario: Chrip Contains Bad Words but one with Punctuation
Given path '/api/chirps'
And request {"body": "Sharbert! This is a kerfuffle opinion I need to share with the world" }

When method post
Then status 200
And match response contains { "cleaned_body": "Sharbert! This is a **** opinion I need to share with the world" }