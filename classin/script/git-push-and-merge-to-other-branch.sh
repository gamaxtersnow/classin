# 运行方式: sh git-push-and-merge-to-other-branch.sh | 参数1: 目标branch(1:master,2:待定) | 参数2: 提交代码的备注信息

#获取当前分支名称
now_branch=$(git branch --show-current)
echo "当前分支:" $now_branch

#提交本地代码
if [ "$1" == 1 ] ;then
      remote_branch="master"
   # elif [ "$1" == 2 ] ;then
      # remote_branch="xxx"
   else
     echo "参数1:目标branch的值错误(1:master)"
     exit
   fi
echo "目标合并分支:" $remote_branch

#备注
if [ ! -n "$2" ] ;then
   echo "备注不能为空"
   exit 2
else
   mark=$2
fi
echo "提交备注:" $mark

#开始操作
git pull
git add .
git commit -m "$mark"
git push

#切换到 {remote_branch} 分支,然后合并当前分支,push后,再次回到当前分支
git checkout $remote_branch
git pull
git merge $now_branch
git push
git checkout $now_branch