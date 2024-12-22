Hi! This is a simple calculation service. You can use it everywhere!

Features:
    //Perform addition, subtraction, division, and multiplication.
    //Validate your expression for mistakes.
    //It is really fast.

Preconditions:
    //You need to have Go version 1.23.1 installed.

Installation:
    //Create a folder and open a terminal.
    //Clone the repository with: git clone https://github.com/ArteShow///Calculation_Service.git
    //Run this command: go mod tidy
    //Check if you have a program to make HTTP requests. If not, I recommend Postman (this is not an advertisement).

Before Using:
    //Check the port. If you want to change it, open application///application.go and scroll to the bottom. There you will see a function with the port configuration, which you can modify.
    //Prepare the request URL: http://localhost:(your_port)
    //Use JSON as the body format. For example:
    {
        "expression": "(your expression)"
    }

How to Use:
    //Open the terminal again and write: go run cmd/main.go
    //If prompted for permission to activate the service, allow it.
    //Then you can run your request and get the answer.

Answer Possibilities:
    {
        "result": "(Result of your expression)"
    }
    {
        "error": "Division by zero"
    }
    //Error (404): If you get this code, try again.
    //Error (500): If you see the message "Internal server error," restart your program (close the terminal or type "exit").

Contact Me:
    //If you have a problem, contact me: soka.rtemax@gmail.com