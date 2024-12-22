Calculation Service

About:
This is a simple web-calculation-service. 

Features:
    //It can do the addition, subtraktion, division and multiplikation
    //It can check your expression for mistakes
    //It is really fast

Preconditions:
    //You have to install the 1.23.1 go version.

Installation
    //First create a folder and open a terminal.
    //Next clone the repository with: git clone https://github.com/ArteShow/Calculation_Service.git
    //Then use this command: go mod tidy
    //At the end check if you have any programms that can do the http requests. If no then I will recomand Postman.(It isn't a advertising)

//Befor using:
    //Check the port. If you want to change it than open application/application.go and there scroll to the bottom. There you will sea a function and the port in it. You can change it there.
    //Prepear the request: http://localhost:(your_port)
    //Please use json as body formate. For example:
    {
    "expression": "(your expression)"
    }

How to use:
    //First open the termina lfor the second time and write this: go run cmd/main.go
    //If it want the primory to activate the service allowe it.
    //Then you can run your request and get the answer.

Anwser Possibilitys:
    //If your expression is correct than you will get this as result and code 200:
    {
    "result": "(Result of your expression)"
    }
    //If you did a mistake for example 2/0 than you will get this with code 422:
    {
        "error": "Division by zero"
    }
    //If you got an error and the code 404 than you should try again.
    //If you got the code 500 and error message:"International server error" than you should restart your programm(close the terminal or "exit")

Contact me:
    //If you have a problem, cantact me: sokartemax@gmail.com