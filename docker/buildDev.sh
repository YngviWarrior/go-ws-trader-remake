iddocker=$(sudo docker ps | grep ws-go-trader-dev | cut -d' ' -f1);
sudo docker rmi ws-go-trader-dev -f
sudo docker build -f ../wsTrader.Dockerfile -t ws-go-trader-dev ../
