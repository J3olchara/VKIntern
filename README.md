# Run on macOS
1. Download docker-compose and docker  
```bash
brew install docker-compose
brew install docker
```
2. Clone this repository to current directory and create virtual environment dir
```bash
git clone git@github.com:J3olchara/VKIntern.git
cd VKIntern
mkdir env
```
3. 1. If you want to develop something
        ```bash
        cp env_example/dev.db.env-example env/dev.db.env
        cp env_example/dev.env-example env/dev.env
        ```
   2. If you want to make tests
        ```bash
        cp env_example/test.db.env-example env/test.db.env
        cp env_example/test.env-example env/test.env
        ```
   3. If you want to deploy in production
        ```bash
        cp env_example/prod.db.env-example env/prod.db.env
        cp env_example/prod.env-example env/prod.env
        ```

4. Open ./env directory and fill .env files with your own data
5. 1. If you want to develop something
        ```bash
        docker-compose -f docker-compose.dev.yml -p vk_dev up 
        ```
   2. If you want to make tests
        ```bash
        docker-compose -f docker-compose.test.yml -p vk_test up --remove-orphans --exit-code-from test
        ```
   3. If you want to deploy in production
        ```bash
        docker-compose -f docker-compose.prod.yml -p vk up --remove-orphans
        ```
      
6. Application has been started on localhost:8000 or HOST:PORT from .env file

# Linux
1. Download docker-compose and docker
```bash
sudo apt-get install docker-compose
sudo apt-get install docker
```
2. Clone this repository to current directory and create virtual environment dir
```bash
git clone git@github.com:J3olchara/VKIntern.git
cd VKIntern
mkdir env
```
3. 1. If you want to develop something
        ```bash
        cp env_example/dev.db.env-example env/dev.db.env
        cp env_example/dev.env-example env/dev.env
        ```
   2. If you want to make tests
        ```bash
        cp env_example/test.db.env-example env/test.db.env
        cp env_example/test.env-example env/test.env
        ```
   3. If you want to deploy in production
        ```bash
        cp env_example/prod.db.env-example env/prod.db.env
        cp env_example/prod.env-example env/prod.env
        ```

4. Open ./env directory and fill .env files with your own data
5. 1. If you want to develop something
        ```bash
        docker-compose -f docker-compose.dev.yml -p vk_dev up 
        ```
   2. If you want to make tests
        ```bash
        docker-compose -f docker-compose.test.yml -p vk_test up --remove-orphans --exit-code-from test
        ```
   3. If you want to deploy in production
        ```bash
        docker-compose -f docker-compose.prod.yml -p vk up --remove-orphans
        ```
      
6. Application has been started on localhost:8000 or HOST:PORT from .env file
