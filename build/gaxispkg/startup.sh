#!/bin/sh
#apt install -y gcc libgmp3-dev
if [ -z "`dpkg --list|grep libgmp3-dev`" ]; then
	echo "Not found libgmp3-dev, install it!"
        apt install -y libgmp3-dev
fi
if [ -z "`dpkg --list|grep gcc`" ]; then
	echo "Not found gcc, install it!"
        apt install -y gcc
fi
show_usage="args: [-d ,-k, -p, -n,-r,-h]\
                                  [--datadir=,--keystore=, --port=, --net=, --rpc=,--help]"
export DYLD_LIBRARY_PATH="./czero/lib/"
export LD_LIBRARY_PATH="./czero/lib/"
DEFAULT_DATD_DIR="./data"
LOGDIR="./log"
DEFAULT_PORT=33896
CONFIG_PATH="./gaxisConfig.toml"
DATADIR_OPTION=${DEFAULT_DATD_DIR}
NET_OPTION=""
RPC_OPTION=""
PORT_OPTION=${DEFAULT_PORT}
KEYSTORE_OPTION=""


GETOPT_ARGS=`getopt -o d:k:p:n:r:h -al datadir:,keystore:,port:,net:,rpc:,help -- "$@"`
eval set -- "$GETOPT_ARGS"
while [ -n "$1" ]
do
        case "$1" in
                -d|--datadir) DATADIR_OPTION=$2; shift 2;;
                -p|--port) PORT_OPTION=$2; shift 2;;
                -n|--net) NET_OPTION=--$2; shift 2;;
                -k|--keystore) KEYSTORE_OPTION="--keystore $2"; shift 2;;
                -r|--rpc)
                        #localhost=$(hostname -I|awk -F ' ' '{print $1}')
			localhost="0.0.0.0"
                        RPC_OPTION="$cmd --rpc --rpcport $2 --rpcaddr $localhost  --rpccorsdomain=* --rpcvhosts=*"; shift 2;;
                -h|--help) echo $show_usage;exit 0;;
                --) break ;;
        esac
done

cmd="bin/gaxis --config ${CONFIG_PATH} --lightNode --exchange --confirmedBlock 12 --mineMode --stake --recordBlockShareNumber --datadir ${DATADIR_OPTION} --port ${PORT_OPTION} ${NET_OPTION} ${RPC_OPTION} ${KEYSTORE_OPTION}  --rpcapi axis,light,stake,net,txpool,exchange"
mkdir -p $LOGDIR

echo $cmd
current=`date "+%Y-%m-%d"`
logName="gaxis_$current.log"
sh stop.sh
nohup ${cmd} >> "${LOGDIR}/${logName}" 2>&1 & echo $! > "./pid"
