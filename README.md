# Emerging Trends in IT Assignment 1
**By Nur Hakimi B Mohd Yasman (S10206177C)**


This repository constitutes a ride-sharing platform named _**'DRIVE'**_ that makes use of microservice architechture. The platform has 2 primary group of users, namely the passengers and drivers. Users can create either account. Users can update most of the information in their account, but they are not able to delete their accounts for audit purposes. A passenger can request for a trip, and this process should follow similarly to other platforms such as Grab and Gojek. The passengers can retrieve their rides history in reverse chronological order.

## **Design considerations of microservices**

1. Security
    * As with any platform that handles sensitive user information, security is a crucial consideration for a ride-sharing platform. Sensitive information should be protected and not thoughtlessly shared. Certain functions should only be accessed by certain users. 
    * This is achieved through the use of a JWT microservice, that generates a new token everytime users log in. This token is then shared to the console microservice, that can then call the different functions of different microservices using that token. On each microservice, there is a JWT validator function that only allows the call to the function to go through if the token provided is valid. This way, no unauthorized access to the methods and information is explicitly allowed.

1. Service decomposition
    * Each microservice is designed to focus on a specific aspect of the overall platform. In the terms of Domain-Driven Design, under the domain of user management would be the subdomains of account creation, account sign-in, etc. The microservice is represented by the domain, with each subdomain being represented by the more specific functions. 
    * There would then be four domains - account management, user authentication, ride functionalities and ride history.
    * This ensures that the platform is modular and scalable, meaning that it is easier to add new functionalities or remove existing functionalities as needed. This can help to ensure that the platform can easily adapt to changing business needs and can scale to meet changing demand.
    * This also helps with the maintainability of the platform. By dividing the platform into smaller microservices, it is easier to make changes to specific parts of the system without impacting the overall platform. This can help to reduce the risk of introducing bugs or other issues into the system, and can make it easier to diagnose and fix problems when they do arise.
    
1. Database management
    * The platform should have a consistent approach to managing databases, including how data is stored, accessed, and shared between microservices. 
    * In this case, each database should only be accessed through a certain microservice, and the console itself should not have any access to any databases. If data from a certain database is needed, it can be retrieved via an API call to the microservice in the same domain.
    * The data should also stay consistent across the different databases and tables. This is done by using primary keys to compare records and transfer data, accomplished using MySQL's auto_increment primary key function for the User IDs.
    
 ## **Architecture diagram**
    
![Architecture Diagram](https://github.com/nhaaki/ETI_Assignment1/blob/4820e0a2970f2a2fc23c07dc48ef6544e3f4c4c3/ETI%20Assignment%201%20Microservice%20Diagram.png)

 ## **Instructions to set up and run these microservices**
 
1. Clone this repository in Visual Studio Code.
2. Run the MySQL script using a MySQL connection with the IP (127.0.0.01:3306). The MySQL script will create the user **`user`** with the password of **`password`**. Followed by creating three databases with their own respective tables.
2. Open 6 terminals, 2 for each user account (Passenger and Driver) and 1 for every other microservice.
3. Follow the code for each terminal.
    * Terminal 1
      ```
      cd ./console
      go run main.go
      ```
    * Terminal 2
      ```
      cd ./console
      go run main.go
      ```
    * Terminal 3
      ```
      cd ./jwt_client
      go run main.go
      ```
    * Terminal 4
      ```
      cd ./account
      go run main.go
      ```
    * Terminal 5
      ```
      cd ./bookride
      go run main.go
      ```
    * Terminal 6
      ```
      cd ./ridehistory
      go run main.go
      ```
      
4. From there, you can begin using the platform through the console terminals. You should not need to look at the other terminals.
