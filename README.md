<!-- go build . && sudo ./TM-Rodocker run -ti sh -->
go build . && sudo ./TM-Rodocker run -ti -v /root/volume:/containerVolume sh 


非標記參數要在最後面
go build . && sudo ./TM-Rodocker run -ti sh -v /root/volume:/containerVolume
這樣下會錯是因為 -ti 是boolFlag 所以他不會帶參數 自然sh 就會被判定為非標記參數 其後面都會變成非標記參數