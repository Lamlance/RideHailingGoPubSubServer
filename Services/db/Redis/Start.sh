#!/bin/sh
sleep 5s
CONTAINER_ALREADY_STARTED="CONTAINER_ALREADY_STARTED_PLACEHOLDER"
if [ ! -e $CONTAINER_ALREADY_STARTED ]; then
  echo "-- First container startup --"
  # YOUR_JUST_ONCE_LOGIC_HERE
  for i in {0..5}
    do
      ping_success=$(redis-cli ping)
      if [ "$ping_success" == "PONG" ]; then
        redis-cli GEOADD w3gv 106.69377232519079 10.789996745967283 LeVanTamParkDriver
        redis-cli GEOADD w3gv 106.69025713826234 10.789567729608098 HaiBaTrungSchoolDriver
        redis-cli GEOADD w3gv 106.69242064275858 10.774643276016672 TaoDanParkDriver
        redis-cli GEOADD w3gv 106.69533120919584 10.777134323246075 DinhDocLapDriver
        touch $CONTAINER_ALREADY_STARTED
        echo "-- Population success --"
        break
      else
        sleep 5s
      fi
  done
else
    echo "-- Not first container startup --"
fi
echo "-- Closing population script --"