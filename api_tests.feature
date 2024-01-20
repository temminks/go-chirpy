Feature: API Tests

Background:
* url "http://localhost:8080"

Scenario: Valid Chirp
Given path '/api/validate_chirp'
And request { "body": "I had something interesting for breakfast" }

When method post
Then status 200
And match response contains { "valid": true }

Scenario: Chirp too long
Given path '/api/validate_chirp'
And request { "body": "I had something interesting for breakfast. It did not only taste great but also looked fantastic. On top of that it is extremely healthy and contains a variety of essential amino acids." }

When method post
Then status 400
And match response contains { "error": "Chirp is too long" }

Scenario: Chirp too long
Given path '/api/validate_chirp'
And request { "id":  1}

When method post
Then status 400
And match response contains { "error": "Chirp is empty" }