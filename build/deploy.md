##Docker Compose部署
檔案結構
```
│  docker-compose.yml
│  Dockerfile
│  nginx.conf
├─api
│  └─main //SimpleAccount API
│     │  config.json
│     └─ main
├─html //SimpleAccount 前端檔案
│  └──index.html
├─mysql
└─ssl
```

`docker-compose.yml`
```
version: '3'

services:
  simple_account_api:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - .\api\main:/data
    entrypoint: /data/main
    working_dir: /data
    networks:
      - simple_account_api_network

  simple_account_nginx:
    image: nginx
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - simple_account_api
      - simple_account_api_oauth
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl/
      - ./html:/usr/share/nginx/html
    networks:
      - simple_account_api_network

  simple_account_mysql:
    image: mysql
    environment:
      MYSQL_ROOT_PASSWORD: example_password
      MYSQL_DATABASE: example_db
      MYSQL_USER: example_user
      MYSQL_PASSWORD: example_password
    volumes:
      - ./mysql:/var/lib/mysql
    networks:
      - simple_account_api_network

networks:
  simple_account_api_network:
    name: simple_account_api_network
```

`nginx.conf`
```
http {
    *
    *
    *
    server {
        listen 443 ssl;
        listen [::]:443 ssl;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        location /api {
            proxy_pass http://simple_account_api:80;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        location / {
            root /usr/share/nginx/html;
            index index.html;
        }
    }
}
```