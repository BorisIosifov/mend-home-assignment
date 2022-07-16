# Home Assignment for Mend from Boris Iosifov

### Running the application with the mysql storage
The application and the mysql server could be started by the docker-compose.
Run:
```sh
docker-compose build
```
Then:
```sh
docker-compose up -d
```
Then wait about 40 seconds (no less than 30 seconds) and the application should be ready to use. It listen 443 port on your localhost.

### Working with the application
If you will use the curl command, please set '-k' flag, because I use self-signed certificate for TLS. The API available by the `https://localhost/` url. For now API provides two object types: books and cars.

##### Examples
To get list of books:
```sh
curl -k https://localhost/books/
```
To get list of cars:
```sh
curl -k https://localhost/cars/
```
To get a car:
```sh
curl -k https://localhost/cars/2/
```
To add a new car:
```sh
curl -k -X POST -d "brand=Mitsubishi&model=ASX" https://localhost/cars/
```
To update a car:
```sh
curl -k -X PUT -d "brand=Mitsubishi&model=L200" https://localhost/cars/3/
```
To remove a car:
```sh
curl -k -X DELETE https://localhost/cars/3/
```

### Running the tests
To run tests, use
```sh
docker-compose exec mend-home-assignment bash -c "go test -v -cover ./api"
```
You should run it when the docker-compose is up.
Tests don't use any external services as you wish. So you can stop the mysql container:
```sh
docker-compose stop mysql
```
And tests will still pass.

### What I don't quite like in the application
Here are the things which I would like to do or change, but I didn't do them to not overthink (as you wrote).
- I use self-signed certificate as I already said. I didn't have time to find out how to use for example let's Encrypt for the localhost.
- I didn't make a mechanism which waiting when mysql is ready when docker-compose is upping. I only made a sleeping for 30 seconds before starting the application. If I would have more time, I would write a short sh script which waiting when the 3306 port is up.
- I did just one database connection and use mutex to get rid of using it from several goroutines in the same time. If I would have more time I would make for example several goroutines to calling sql queries and each goroutine would have it own connection.
- I didn't provide references between objects.
- I didn't take out some parameters (like a paths to certificates and database credentials) to a config file.
- I didn't integrate the API with some nosql storage. I just did possibility to store data in the local memory. It helped me to write tests independent on external services. But I thinks it is not the same as a nosql storage.
- I don't well like how I did the getFieldsList function. I did it to avoid listing of fields of each object type in the mysql integration. But the realisation is not so good.
- I couldn't avoid a listing of object types in the GetList and Get methods in mysql. I would can avoid it in the Get method if I would add a pointer to object in the attributes list of the method (like it done in the json.Unmarshal). But I don't know how to avoid it in the GetList method, because []object.Car don't correspond to []object.Object.
- I didn't add any way to stop the application except stopping the container.

