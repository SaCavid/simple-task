@ECHO OFF
git add .
git commit -m "Finishing"
PAUSE
git push origin master
PAUSE
ECHO Congratulations! Your first batch file executed successfully.
PAUSE
docker image build -t sacavid/simple-task .
PAUSE
docker push sacavid/simple-task:latest
PAUSE
docker-compose up
