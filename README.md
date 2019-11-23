# fa19-281-marss : Air BNB Booking

## Team - MARSS

## Team Members and work done:
         
    1. Jaya Sindhu Manda - Bookings API and AWS Gateway, SNS
    2. Rohan - Frontend, Mongo sharding, Profile API
    3. Apeksha - Poperty listing API, AWS gateway, Kuberneted, S3 bucket
    4. Mrinalini - Users API, Firebase authentication
    5. Siddesh - Search API, Redis Cache

## Summary of Project:
It is similar to AIRBNB where a propert owner and customer can login. Property owner to register their property. Customer to book a property. A customer can search a property and book a property. After booking a property notification will be sent to customer from configured SNS based on the details provided by customer during booking.

## Architecture Diagram:
![alt text](281Diagram3.png)


## Summary of Key Features:
1. **Microservices** : We have used Golang for the development of all the services like Login, Profile creation, Property Publishing, Property searching and Booking.

2. **MongoDB**: For storing the users data and consistent fast retrieval, We chose mongo db. Since mongo db is a non-SQL database, it is very convenient to store the data without strict schema definition.

3. **JWT**: For maintaining sessions, we are using JWTs over the traditional cookie-based sessions. JWTs are as much secure and much more scalable since the server overhead is very less. The JWT signature and the client secret ensures that there was no tampering with the JWT.

4. **Redis**: For storing the JWT secret, We are using Redis as the in-memory caching layer. Since this will be invoked on every API call, there was a need for in-memory storage which is highly available. The Redis also provided in-built TTL which ensured that a session becomes invalid after 60 minutes.

5. **Docker**: Deployed few of the services in AWS instances. This made it even more scalable.

6. **MongoDB Sharding**: As the size of the data increase, data partitioning becomes very important. For setting this, we created two sharded cluster, two config servers, and one mongos for each service.

7. **GKE(Kubernetes)**: A managed Kubernetes cluster by Google Cloud is used to deploy docker containers. Load balancer is exposed on the same.

8. **SNS**: AWS SNS is used to notify the customers regarding their booking updastes. For every booking corresponding user is subscribed and updates are being sent to registered e-mail id's.

9. **S3**: Property listing images are stored in S3 bucket and retrieved on GET.

10. **FireBase**: For authentication purposes.

11. AWS API Gateway: All the LB's and kubernetis IP is configured in API Gateway.

## Team Contributions:

- Jaya Sindhu Manda:
	- Intial POC using GO+MONGO+REST -> Done (11-Nov)
	- Docker configuration for the Go API -> Done (13-Nov)
	- Backend API Development for Booking -> Done (16-Nov)
	- Backend API integrating with mongo cluster -> Done (19-Nov)
	- SNS Implementation as part of Booking,Updation and Deletion (To notify user with booking ID)-> Done (19-Nov)
	- API Gateway configuration -> Done
	- AWS Cloud Formation - Done
         
 - Rohan:
	- Create an API contract
         - Develop Backend API for User profile
	- Build frontend Web application using React.js
	- Integrate all the APIs with frontend 
	- Deploy frontend on Heroku platform
         
- Apeksha:
	- Create/update/get/delete Listings API
	- s3 bucket implementation for all CRUD operations -
	- Kubernetes Deployment 
	- VPC Link	

- Mrinalini:
	- Evaluate CI/CD Tools 
	- Develop backend API for login/sign up	
	- Setting up mongo db cluster and integration with user apis
	- Investigate firebase and Auth0 authentication mechanisms
	- Implement user authentication using firebase 

- Siddesh:
	- Backend API development for Search - Done
	- Backend API intergration with mongo cluster - Done
	- VPC Peering - In Progress
	- Kubernetes Deployment - In Progress


