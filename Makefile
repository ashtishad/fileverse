run:
	export API_HOST=127.0.0.1 \
	export API_PORT=8000 \
	export DB_USER=postgres \
	export DB_PASSWD=postgres \
	export DB_HOST=127.0.0.1 \
	export DB_PORT=5432 \
	export DB_NAME=fileverse \
	export IPFS_ADDR=127.0.0.1:5001 \
&& go run main.go
