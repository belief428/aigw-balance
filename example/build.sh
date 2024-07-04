name=aigw-balance

docker stop $name

docker rm $name && docker rmi $name

docker build -t $name .

docker run -d -p 3000:3000 --restart=always --name=$name --privileged=true $name
