iddocker=$(sudo docker ps | grep ws-go-trader-dev | cut -d' ' -f1);
sudo docker stop $iddocker
sudo docker run -d -p 2096:2096
ws-go-trader-dev
