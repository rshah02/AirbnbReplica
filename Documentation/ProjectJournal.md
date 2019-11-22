Project journal:

**Week 1**

1. Project topic selection: Airbnb clone
2. Backend API selection (to be implemented in Go):
	User Sign up/Login/ Authentication using jwt - Mrinalini
	Profile page for customer and host - Rohan
	Posting ad/Removal of accomodation by host - Apeksha
	Search - Siddesh
	Booking/Cancelling/Updation accomodation - Sindhu


**Initial Responsibilities:**

- Frontend development - Rohan
- CI/CD Tools - Mrinalini
- Kubernetes - Siddesh
- Docker - Apeksha
- API Gateway - Sindhu

**To Do:**

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

- Rohan:
	- Create an API contract
    - Develop Backend API for User profile
	- Build frontend Web application using React.js
	- Integrate all the APIs with frontend 
	- Deploy frontend on Heroku platform
- Siddesh:
	- Backend API development for Search - Done
	- Backend API intergration with mongo cluster - Done
	- VPC Peering - In Progress
	- Kubernetes Deployment - In Progress

- Sindhu:
	- Intial POC using GO+MONGO+REST -> Done (11-Nov)
	- Docker configuration for the Go API -> Done (13-Nov)
	- Backend API Development for Booking -> Done (16-Nov)
	- Backend API integrating with mongo cluster -> Done (19-Nov)
	- SNS Implementation as part of Booking,Updation and Deletion (To notify user with booking ID)-> Done (19-Nov)
	- API Gateway configuration -> In Progress
	- AWS Cloud Formation - Exploring
