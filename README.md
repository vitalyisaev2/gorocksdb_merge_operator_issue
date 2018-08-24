# gorocksdb_merge_operator_issue

Steps to reproduce:

Prerequisites:
1. ```go >= 1.10```
2. ```librocksdb.so >= 5.13```
3. ```valgrind```
4. ```massif-visualizer```


2. go get -v https://github.com/vitalyisaev2/gorocksdb_merge_operator_issue
3. Download database dump frome Google Drive: https://drive.google.com/file/d/13pn0ZW2qt4Tb9c5hPYer0HjgGt_rJtNR/view?usp=sharing
4. tar xzvf segments.tar.gz
5. gorocksdb_merge_operator_issue iteratea
