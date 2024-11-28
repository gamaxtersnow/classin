if [ ! -n "$1" ] ;then
    echo "表名不能为空"
    exit 2
else
    tablename=$1
fi

goctl model mysql datasource -url="meishiedu:Meishi@2021@tcp(rm-2zefj5975i499k6q05o.mysql.rds.aliyuncs.com:3306)/updeducrm" -table="$tablename" -dir="./model/oldcrm"