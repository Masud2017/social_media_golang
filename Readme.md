# Social media api using golang
This a social media api that is created using dgraph database(dgo go client).

The project is dockerized means that you only need to run the docker compose command and it will work fine.

docker compose command:
```
 docker compose up -d
```


# About the api :
### Api endpoint list : 
#### The api is accessible to the port - 4442


Since this project is a test project, I didn't made any of them post. All of those endpoints are accessible by GET request
- /signup [GET] -> params : ?name=[name]&email=[email]&password=[password] 
- /userlist [GET]
- /me/:my_id [GET]
- /acceptreq/:user_id/:req_id [GET] 
- /cancelreq/:user_id/:req_id [GET]
- /addfriend [GET] -> ?my_id=[your user id]d&req_to=[user id of the user that you try to req to]
- /addfather [GET] -> ?my_id=[your user id]d&req_to=[user id of the user that you try to req to]
- /addmother [GET] -> ?my_id=[your user id]d&req_to=[user id of the user that you try to req to]
- /addson [GET] -> ?my_id=[your user id]d&req_to=[user id of the user that you try to req to]
- /myrelationlist/:user_id [GET]
- /relationship_reqs/:user_id [GET]
- /my_relationship_reqs/:user_id [GET]