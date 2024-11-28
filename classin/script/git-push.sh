# 运行方式: sh git-push.sh 提交代码的备注信息

if [ ! -n "$1" ] ;then
    # mark="修改"
    echo "备注不能为空"
    exit 2
else
    mark=$1
fi

git pull
git add .
git commit -m "$mark"
git push
