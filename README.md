ngx_proxy_store_file_mgr
========================


Description:
  this project is a solution to solve nginx proxy store without file manager. Since nginx 's proxy_store directive is just responsible for saving files to disk and does not consider time to delete files. ngx_proxy_store_file_mgr is a process used to delete files on the disk using LRU .


Notes:
  1. file access time saved in redis database, using sorted set;
  2. ngx_proxy_store_file_mgr uses json format config files, as follows :
  {
  "MaxFileLimit" : 100000, 
	"CheckInterval" : 20, //check every 20s 
	"ExpireDays" : 7,
	"ErrorLog": true,
	"AccessLog": false,
	"SortedSetName": "defset",
	"HashName": "defhash",
	"RedisAddr" : "127.0.0.1:6379",
	"RoutineCount" : 32
  }

  but currently the lru condition is just disk percentage(low than 20%), and later will implement the specified config parameter.
  3. provides restful api to shutdown server: curl http://ip:10000/stop
