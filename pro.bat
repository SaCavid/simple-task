@ECHO OFF
git add .
git commit -m "Finishing"
PAUSE
#"C:\Program Files\Docker\Docker\DockerCli.exe" -SwitchDaemon
#PAUSE
git push origin master
PAUSE
docker image build -t sacavid/simple-task .
PAUSE
docker push sacavid/simple-task:latest
PAUSE
docker-compose up
