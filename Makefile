deploy:
	go build app.go
	scp app foo.db restart.sh root@104.248.132.28:/root
	scp views/show.html root@104.248.132.28:/root/views
	ssh root@104.248.132.28 'chmod +x /root/restart.sh && /bin/bash /root/restart.sh'

