

cd gateway
rm -f gateway
cd ..

cd msg_server
rm -f msg_server
cd ..

cd router
rm -f router
cd ..

cd manager
rm -f manager
cd ..

git add -A

git commit -m "init"

git push origin master
